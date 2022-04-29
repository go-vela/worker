// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/internal/step"
)

// create a step logging pattern.
const stepPattern = "[step: %s]"

// CreateStep configures the step for execution.
func (c *client) CreateStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// create a library step object to facilitate injecting environment as early as possible
	// (PlanStep is too late to inject environment vars for the kubernetes runtime).
	_step := library.StepFromBuildContainer(c.build, ctn)

	// update the step container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Environment
	err = step.Environment(ctn, c.build, c.repo, _step, c.Version)
	if err != nil {
		return err
	}

	// substitute container configuration
	//
	// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.Substitute
	err = ctn.Substitute()
	if err != nil {
		return err
	}

	return nil
}

// PlanStep prepares the step for execution.
func (c *client) PlanStep(ctx context.Context, ctn *pipeline.Container) error {
	// create the library step object
	_step := library.StepFromBuildContainer(c.build, ctn)
	_step.SetStatus(constants.StatusRunning)
	_step.SetStarted(time.Now().UTC().Unix())

	// update the step container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Environment
	err := step.Environment(ctn, c.build, c.repo, _step, c.Version)
	if err != nil {
		return err
	}

	// add the step to the client map
	c.steps.Store(ctn.ID, _step)

	return nil
}

// ExecStep runs a step.
func (c *client) ExecStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// load the step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_step, err := step.Load(ctn, &c.steps)
	if err != nil {
		return err
	}

	// defer taking a snapshot of the step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Snapshot
	defer func() { step.Snapshot(ctn, c.build, nil, nil, nil, _step) }()

	// run the runtime container
	err = c.Runtime.RunContainer(ctx, ctn, c.pipeline)
	if err != nil {
		return err
	}

	// trigger StreamStep goroutine with logging context
	c.streamRequests <- message.StreamRequest{
		Key:       "step",
		Stream:    c.StreamStep,
		Container: ctn,
	}

	// do not wait for detached containers
	if ctn.Detach {
		return nil
	}

	// wait for the runtime container
	err = c.Runtime.WaitContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// inspect the runtime container
	err = c.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// StreamStep tails the output for a step.
func (c *client) StreamStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// tail the runtime container
	rc, err := c.Runtime.TailContainer(ctx, ctn)
	if err != nil {
		return err
	}
	defer rc.Close()

	// create a step pattern for log output
	_pattern := fmt.Sprintf(stepPattern, ctn.Name)

	// check if the container provided is for stages
	_stage, ok := ctn.Environment["VELA_STEP_STAGE"]
	if ok {
		// check if the stage name is set
		if len(_stage) > 0 {
			// create a stage pattern for log output
			_pattern = fmt.Sprintf(stagePattern, _stage, ctn.Name)
		}
	}

	// create new scanner from the container output
	scanner := bufio.NewScanner(rc)

	// scan entire container output
	for scanner.Scan() {
		// ensure we output to stdout
		fmt.Fprintln(os.Stdout, _pattern, scanner.Text())
	}

	return scanner.Err()
}

// DestroyStep cleans up steps after execution.
func (c *client) DestroyStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// load the step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_step, err := step.Load(ctn, &c.steps)
	if err != nil {
		// create the step from the container
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#StepFromContainerEnvironment
		_step = library.StepFromContainerEnvironment(ctn)
	}

	// defer an upload of the step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Upload
	defer func() { step.Upload(ctn, c.build, nil, nil, nil, _step) }()

	// inspect the runtime container
	err = c.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// remove the runtime container
	err = c.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}
