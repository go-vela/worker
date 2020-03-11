// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// CreateStage prepares the stage for execution.
func (c *client) CreateStage(ctx context.Context, s *pipeline.Stage) error {
	init := c.pipeline.Stages[0].Steps[0]
	// TODO: make this cleaner
	result, ok := c.stepLogs.Load(init.ID)
	if !ok {
		return fmt.Errorf("unable to get init step log from client")
	}

	l := result.(*library.Log)

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"stage": s.Name,
	})

	// update the init log with progress
	l.SetData(
		append(
			l.GetData(),
			[]byte(fmt.Sprintf("  $ Pulling step images for stage %s...\n", s.Name))...,
		),
	)

	// create the steps for the stage
	for _, step := range s.Steps {
		// TODO: make this not hardcoded
		// update the init log with progress
		l.SetData(
			append(
				l.GetData(),
				[]byte(fmt.Sprintf("    $ docker image inspect %s\n", step.Image))...,
			),
		)

		logger.Debugf("creating %s step", step.Name)
		// create the step
		err := c.CreateStep(ctx, step)
		if err != nil {
			return err
		}

		c.logger.Infof("inspecting %s step", step.Name)
		// inspect the step image
		image, err := c.Runtime.InspectImage(ctx, step)
		if err != nil {
			return err
		}

		// update the init log with step image info
		l.SetData(append(l.GetData(), image...))
	}

	return nil
}

// PlanStage prepares the stage for execution.
func (c *client) PlanStage(ctx context.Context, s *pipeline.Stage, m map[string]chan error) error {
	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"stage": s.Name,
	})

	logger.Debug("gathering stage dependency tree")
	// ensure dependent stages have completed
	for _, needs := range s.Needs {
		logger.Debugf("looking up dependency %s", needs)
		// check if a dependency stage has completed
		stageErr, ok := m[needs]
		if !ok { // stage not found so we continue
			continue
		}

		logger.Debugf("waiting for dependency %s", needs)
		// wait for the stage channel to close
		select {
		case <-ctx.Done():
			return fmt.Errorf("errgroup context is done")
		case err := <-stageErr:
			if err != nil {
				logger.WithError(err).Errorf("%s stage produced error", needs)
				return err
			}

			continue
		}
	}

	return nil
}

// ExecStage runs a stage.
func (c *client) ExecStage(ctx context.Context, s *pipeline.Stage, m map[string]chan error) error {
	b := c.build
	r := c.repo

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"stage": s.Name,
	})

	// close the stage channel at the end
	defer close(m[s.Name])

	logger.Debug("starting execution of stage")
	// execute the steps for the stage
	for _, step := range s.Steps {
		c.logger.Infof("planning %s step", step.Name)
		// plan the step
		err := c.PlanStep(ctx, step)
		if err != nil {
			return fmt.Errorf("unable to plan step %s: %w", step.Name, err)
		}

		logger.Debugf("executing %s step", step.Name)
		// execute the step
		err = c.ExecStep(ctx, step)
		if err != nil {
			return err
		}

		result, ok := c.steps.Load(step.ID)
		if !ok {
			return fmt.Errorf("unable to get step %s from client", step.Name)
		}

		cStep := result.(*library.Step)

		// check the step exit code
		if step.ExitCode != 0 {
			// check if we ignore step failures
			if !step.Ruleset.Continue {
				// set build status to failure
				b.SetStatus(constants.StatusFailure)
			}

			// update the step fields
			cStep.SetExitCode(step.ExitCode)
			cStep.SetStatus(constants.StatusFailure)
		}

		cStep.SetFinished(time.Now().UTC().Unix())
		c.logger.Infof("uploading %s step state", step.Name)
		// send API call to update the build
		_, _, err = c.Vela.Step.Update(r.GetOrg(), r.GetName(), b.GetNumber(), cStep)
		if err != nil {
			return err
		}
	}

	return nil
}

// DestroyStage cleans up the stage after execution.
func (c *client) DestroyStage(ctx context.Context, s *pipeline.Stage) error {
	// update logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"stage": s.Name,
	})

	// destroy the steps for the stage
	for _, step := range s.Steps {
		logger.Debugf("destroying %s step", step.Name)
		// destroy the step
		err := c.DestroyStep(ctx, step)
		if err != nil {
			return err
		}
	}

	return nil
}
