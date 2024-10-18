// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/worker/internal/step"
)

// CreateStage prepares the stage for execution.
func (c *client) CreateStage(ctx context.Context, s *pipeline.Stage) error {
	// load the logs for the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(c.init, &c.stepLogs)
	if err != nil {
		return err
	}

	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("stage", s.Name)

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library#Log.AppendData
	_log.AppendData([]byte(fmt.Sprintf("> Preparing step images for stage %s...\n", s.Name)))

	// create the steps for the stage
	for _, _step := range s.Steps {
		// update the container environment with stage name
		_step.Environment["VELA_STEP_STAGE"] = s.Name

		_log.AppendData([]byte(fmt.Sprintf("> Preparing step image %s...\n", _step.Image)))

		logger.Debugf("creating %s step", _step.Name)
		// create the step
		err := c.CreateStep(ctx, _step)
		if err != nil {
			return err
		}

		logger.Infof("inspecting image for %s step", _step.Name)
		// inspect the step image
		image, err := c.Runtime.InspectImage(ctx, _step)
		if err != nil {
			return err
		}

		// update the init log with step image info
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#Log.AppendData
		_log.AppendData(image)
	}

	return nil
}

// PlanStage prepares the stage for execution.
func (c *client) PlanStage(ctx context.Context, s *pipeline.Stage, m *sync.Map) error {
	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("stage", s.Name)

	logger.Debug("gathering stage dependency tree")
	// ensure dependent stages have completed
	for _, needs := range s.Needs {
		logger.Debugf("looking up dependency %s", needs)
		// check if a dependency stage has completed
		stageErr, ok := m.Load(needs)
		if !ok { // stage not found so we continue
			continue
		}

		logger.Debugf("waiting for dependency %s", needs)
		// wait for the stage channel to close
		select {
		case <-ctx.Done():
			return fmt.Errorf("errgroup context is done")
		case err := <-stageErr.(chan error):
			if err != nil {
				logger.Errorf("%s stage returned error: %v", needs, err)
				return err
			}

			continue
		}
	}

	return nil
}

// ExecStage runs a stage.
func (c *client) ExecStage(ctx context.Context, s *pipeline.Stage, m *sync.Map) error {
	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("stage", s.Name)

	// close the stage channel at the end
	defer func() {
		// retrieve the error channel for the current stage
		errChan, ok := m.Load(s.Name)
		if !ok {
			logger.Debugf("error channel for stage %s not found", s.Name)

			return
		}

		// close the error channel
		close(errChan.(chan error))
	}()

	logger.Debug("starting execution of stage")

	// stop value determines when a stage's step series should stop executing
	stop := false

	// execute the steps for the stage
	for _, _step := range s.Steps {
		// first check to see if we need to stop the stage before it even starts.
		// this happens when using 'needs' block
		if !s.Independent && c.build.GetStatus() == constants.StatusFailure {
			stop = true
		}

		// set parallel rule that determines whether the step should adhere to entire build
		// status check, which is determined by independent rule and stop value.
		if !stop && s.Independent {
			_step.Ruleset.If.Parallel = true
		} else {
			_step.Ruleset.If.Parallel = false
		}

		// check if the step should be skipped
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Skip
		skip, err := step.Skip(_step, c.build)
		if err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		if skip {
			continue
		}

		// add netrc to secrets for masking in logs
		sec := &pipeline.StepSecret{
			Target: "VELA_NETRC_PASSWORD",
		}
		_step.Secrets = append(_step.Secrets, sec)

		// load any lazy secrets and inject them into container environment
		err = loadLazySecrets(c, _step)
		if err != nil {
			return fmt.Errorf("unable to plan step %s: %w", _step.Name, err)
		}

		logger.Debugf("planning %s step", _step.Name)
		// plan the step
		err = c.PlanStep(ctx, _step)
		if err != nil {
			return fmt.Errorf("unable to plan step %s: %w", _step.Name, err)
		}

		// poll outputs
		opEnv, maskEnv, err := c.outputs.poll(ctx, c.OutputCtn)
		if c.err != nil {
			return fmt.Errorf("unable to exec outputs container: %w", err)
		}

		// merge env from outputs
		//
		//nolint:errcheck // only errors with empty environment input, which does not matter here
		_step.MergeEnv(opEnv)

		// merge env from masked outputs
		//
		//nolint:errcheck // only errors with empty environment input, which does not matter here
		_step.MergeEnv(maskEnv)

		// add masked outputs to secret map so they can be masked in logs
		for key := range maskEnv {
			sec := &pipeline.StepSecret{
				Target: key,
			}
			_step.Secrets = append(_step.Secrets, sec)
		}

		// perform any substitution on dynamic variables
		err = _step.Substitute()
		if err != nil {
			return err
		}

		c.Logger.Debug("injecting non-substituted secrets")
		// inject no-substitution secrets for container
		err = injectSecrets(_step, c.NoSubSecrets)
		if err != nil {
			return err
		}

		logger.Infof("executing %s step", _step.Name)
		// execute the step
		err = c.ExecStep(ctx, _step)
		if err != nil {
			return fmt.Errorf("unable to exec step %s: %w", _step.Name, err)
		}

		// failed steps within the stage should set the stop value to true unless
		// the continue rule is set to true.
		if _step.ExitCode != 0 && !_step.Ruleset.Continue {
			stop = true
		}
	}

	return nil
}

// DestroyStage cleans up the stage after execution.
func (c *client) DestroyStage(ctx context.Context, s *pipeline.Stage) error {
	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("stage", s.Name)

	var err error

	// destroy the steps for the stage
	for _, _step := range s.Steps {
		logger.Debugf("destroying %s step", _step.Name)
		// destroy the step
		err = c.DestroyStep(ctx, _step)
		if err != nil {
			logger.Errorf("unable to destroy step: %v", err)
		}
	}

	return err
}
