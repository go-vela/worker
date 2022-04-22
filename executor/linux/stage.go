// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
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
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := c.Logger.WithField("stage", s.Name)

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte(fmt.Sprintf("> Preparing step images for stage %s...\n", s.Name)))

	// create the steps for the stage
	for _, _step := range s.Steps {
		// update the container environment with stage name
		_step.Environment["VELA_STEP_STAGE"] = s.Name

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
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
		_log.AppendData(image)
	}

	return nil
}

// PlanStage prepares the stage for execution.
func (c *client) PlanStage(ctx context.Context, s *pipeline.Stage, m *sync.Map) error {
	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
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
func (c *client) ExecStage(ctx context.Context, s *pipeline.Stage, m *sync.Map, streamRequests chan<- message.StreamRequest) error {
	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
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
	// execute the steps for the stage
	for _, _step := range s.Steps {
		// check if the step should be skipped
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Skip
		if step.Skip(_step, c.build, c.repo) {
			continue
		}

		logger.Debugf("planning %s step", _step.Name)
		// plan the step
		err := c.PlanStep(ctx, _step)
		if err != nil {
			return fmt.Errorf("unable to plan step %s: %w", _step.Name, err)
		}

		logger.Infof("executing %s step", _step.Name)
		// execute the step
		err = c.ExecStep(ctx, _step, streamRequests)
		if err != nil {
			return fmt.Errorf("unable to exec step %s: %w", _step.Name, err)
		}
	}

	return nil
}

// DestroyStage cleans up the stage after execution.
func (c *client) DestroyStage(ctx context.Context, s *pipeline.Stage) error {
	// update engine logger with stage metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
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
