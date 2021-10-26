// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/step"
)

// create a stage logging pattern.
const stagePattern = "[stage: %s][step: %s]"

// CreateStage prepares the stage for execution.
func (c *client) CreateStage(ctx context.Context, s *pipeline.Stage) error {
	// create a stage pattern for log output
	_pattern := fmt.Sprintf(stagePattern, c.init.Name, c.init.Name)

	// output init progress to stdout
	fmt.Fprintln(os.Stdout, _pattern, "> Pulling step images for stage", s.Name, "...")

	// create the steps for the stage
	for _, _step := range s.Steps {
		// update the container environment with stage name
		_step.Environment["VELA_STEP_STAGE"] = s.Name

		// create the step
		err := c.CreateStep(ctx, _step)
		if err != nil {
			return err
		}

		// inspect the step image
		image, err := c.Runtime.InspectImage(ctx, _step)
		if err != nil {
			return err
		}

		// output the image information to stdout
		fmt.Fprintln(os.Stdout, _pattern, string(image))
	}

	return nil
}

// PlanStage prepares the stage for execution.
func (c *client) PlanStage(ctx context.Context, s *pipeline.Stage, m *sync.Map) error {
	// ensure dependent stages have completed
	for _, needs := range s.Needs {
		// check if a dependency stage has completed
		stageErr, ok := m.Load(needs)
		if !ok { // stage not found so we continue
			continue
		}

		// wait for the stage channel to close
		select {
		case <-ctx.Done():
			return fmt.Errorf("errgroup context is done")
		case err := <-stageErr.(chan error):
			if err != nil {
				return err
			}

			continue
		}
	}

	return nil
}

// ExecStage runs a stage.
func (c *client) ExecStage(ctx context.Context, s *pipeline.Stage, m *sync.Map) error {
	// close the stage channel at the end
	defer func() {
		errChan, ok := m.Load(s.Name)
		if !ok {
			return
		}

		close(errChan.(chan error))
	}()

	// execute the steps for the stage
	for _, _step := range s.Steps {
		// check if the step should be skipped
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Skip
		if step.Skip(_step, c.build, c.repo) {
			continue
		}

		// plan the step
		err := c.PlanStep(ctx, _step)
		if err != nil {
			return fmt.Errorf("unable to plan step %s: %w", _step.Name, err)
		}

		// execute the step
		err = c.ExecStep(ctx, _step)
		if err != nil {
			return fmt.Errorf("unable to exec step %s: %w", _step.Name, err)
		}
	}

	return nil
}

// DestroyStage cleans up the stage after execution.
func (c *client) DestroyStage(ctx context.Context, s *pipeline.Stage) error {
	var err error

	// destroy the steps for the stage
	for _, _step := range s.Steps {
		// destroy the step
		err = c.DestroyStep(ctx, _step)
		if err != nil {
			fmt.Fprintln(os.Stdout, "unable to destroy step: ", err)
		}
	}

	return err
}
