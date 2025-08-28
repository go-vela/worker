// SPDX-License-Identifier: Apache-2.0

package local

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/worker/internal/build"
	"github.com/go-vela/worker/internal/outputs"
	"github.com/go-vela/worker/internal/step"
)

// CreateBuild configures the build for execution.
func (c *client) CreateBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, nil, c.err, nil) }()
	// Check if storage client is initialized
	// and if storage is enabled
	if c.Storage == nil {
		return fmt.Errorf("storage client is not initialized")
	}

	// update the build fields
	c.build.SetStatus(constants.StatusRunning)
	c.build.SetStarted(time.Now().UTC().Unix())
	c.build.SetHost(c.Hostname)
	c.build.SetDistribution(c.Driver())
	c.build.SetRuntime(c.Runtime.Driver())

	// setup the runtime build
	c.err = c.Runtime.SetupBuild(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to setup build %s: %w", c.pipeline.ID, c.err)
	}

	// load the init step from the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadInit
	c.init, c.err = step.LoadInit(c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to load init step from pipeline: %w", c.err)
	}

	// create the step
	c.err = c.CreateStep(ctx, c.init)
	if c.err != nil {
		return fmt.Errorf("unable to create %s step: %w", c.init.Name, c.err)
	}

	// plan the step
	c.err = c.PlanStep(ctx, c.init)
	if c.err != nil {
		return fmt.Errorf("unable to plan %s step: %w", c.init.Name, c.err)
	}

	return c.err
}

// PlanBuild prepares the build for execution.
func (c *client) PlanBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, nil, c.err, nil) }()

	// load the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_init, err := step.Load(c.init, &c.steps)
	if err != nil {
		return err
	}

	// defer taking a snapshot of the init step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#SnapshotInit
	defer func() { step.SnapshotInit(c.init, c.build, nil, nil, _init, nil) }()

	// create a step pattern for log output
	_pattern := fmt.Sprintf(stepPattern, c.init.Name)

	// check if the pipeline provided has stages
	if len(c.pipeline.Stages) > 0 {
		// create a stage pattern for log output
		_pattern = fmt.Sprintf(stagePattern, c.init.Name, c.init.Name)
	}

	// create the runtime network for the pipeline
	c.err = c.Runtime.CreateNetwork(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to create network: %w", c.err)
	}

	// output init progress to stdout
	fmt.Fprintln(c.stdout, _pattern, "> Inspecting runtime network...")

	// inspect the runtime network for the pipeline
	network, err := c.Runtime.InspectNetwork(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect network: %w", err)
	}

	// output the network information to stdout
	fmt.Fprintln(c.stdout, _pattern, string(network))

	// create the runtime volume for the pipeline
	err = c.Runtime.CreateVolume(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to create volume: %w", err)
	}

	// output init progress to stdout
	fmt.Fprintln(c.stdout, _pattern, "> Inspecting runtime volume...")

	// inspect the runtime volume for the pipeline
	volume, err := c.Runtime.InspectVolume(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect volume: %w", err)
	}

	// output the volume information to stdout
	fmt.Fprintln(c.stdout, _pattern, string(volume))

	return c.err
}

// AssembleBuild prepares the containers within a build for execution.
func (c *client) AssembleBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, nil, c.err, nil) }()

	// load the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_init, err := step.Load(c.init, &c.steps)
	if err != nil {
		return err
	}

	// defer an upload of the init step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Upload
	defer func() { step.Upload(c.init, c.build, nil, nil, _init) }()

	// create a step pattern for log output
	_pattern := fmt.Sprintf(stepPattern, c.init.Name)

	// check if the pipeline provided has stages
	if len(c.pipeline.Stages) > 0 {
		// create a stage pattern for log output
		_pattern = fmt.Sprintf(stagePattern, c.init.Name, c.init.Name)
	}

	// output init progress to stdout
	fmt.Fprintln(c.stdout, _pattern, "> Preparing service images...")

	// create the services for the pipeline
	for _, _service := range c.pipeline.Services {
		// TODO: remove this; but we need it for tests
		_service.Detach = true

		fmt.Fprintln(c.stdout, _pattern, fmt.Sprintf("> Preparing service image %s...", _service.Image))

		// create the service
		c.err = c.CreateService(ctx, _service)
		if c.err != nil {
			return fmt.Errorf("unable to create %s service: %w", _service.Name, c.err)
		}

		// inspect the service image
		image, err := c.Runtime.InspectImage(ctx, _service)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s service: %w", _service.Name, err)
		}

		// output the image information to stdout
		fmt.Fprintln(c.stdout, _pattern, string(image))
	}

	// output init progress to stdout
	fmt.Fprintln(c.stdout, _pattern, "> Preparing stage images...")

	// create the stages for the pipeline
	for _, _stage := range c.pipeline.Stages {
		// TODO: remove hardcoded reference
		//

		if _stage.Name == "init" {
			continue
		}

		// create the stage
		c.err = c.CreateStage(ctx, _stage)
		if c.err != nil {
			return fmt.Errorf("unable to create %s stage: %w", _stage.Name, c.err)
		}
	}

	// output init progress to stdout
	fmt.Fprintln(c.stdout, _pattern, "> Preparing step images...")

	// create the steps for the pipeline
	for _, _step := range c.pipeline.Steps {
		// TODO: remove hardcoded reference
		if _step.Name == "init" {
			continue
		}

		fmt.Fprintln(c.stdout, _pattern, fmt.Sprintf("> Preparing step image %s...", _step.Image))

		// create the step
		c.err = c.CreateStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to create %s step: %w", _step.Name, c.err)
		}

		// inspect the step image
		image, err := c.Runtime.InspectImage(ctx, _step)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s step: %w", _step.Name, err)
		}

		// output the image information to stdout
		fmt.Fprintln(c.stdout, _pattern, string(image))
	}

	// output a new line for readability to stdout
	fmt.Fprintln(c.stdout, "")

	c.err = c.outputs.create(ctx, c.OutputCtn, (int64(60) * 30))
	if c.err != nil {
		return fmt.Errorf("unable to create outputs container: %w", c.err)
	}

	// assemble runtime build just before any containers execute
	c.err = c.Runtime.AssembleBuild(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to assemble runtime build %s: %w", c.pipeline.ID, c.err)
	}

	return c.err
}

// ExecBuild runs a pipeline for a build.
func (c *client) ExecBuild(ctx context.Context) error {
	// defer an upload of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Upload
	defer func() { build.Upload(c.build, nil, c.err, nil) }()

	// output maps for dynamic environment variables captured from volume
	var opEnv, maskEnv map[string]string

	// execute outputs container
	c.err = c.outputs.exec(ctx, c.OutputCtn)
	if c.err != nil {
		return fmt.Errorf("unable to exec outputs container: %w", c.err)
	}

	// execute the services for the pipeline
	for _, _service := range c.pipeline.Services {
		// plan the service
		c.err = c.PlanService(ctx, _service)
		if c.err != nil {
			return fmt.Errorf("unable to plan service: %w", c.err)
		}

		// execute the service
		c.err = c.ExecService(ctx, _service)
		if c.err != nil {
			return fmt.Errorf("unable to execute service: %w", c.err)
		}
	}

	// execute the steps for the pipeline
	for _, _step := range c.pipeline.Steps {
		// TODO: remove hardcoded reference
		if _step.Name == "init" {
			continue
		}

		// check if the step should be skipped
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Skip
		skip, err := step.Skip(_step, c.build, c.build.GetStatus())
		if err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		if skip {
			logrus.Infof("Skipping step %s due to ruleset", _step.Name)
			continue
		}

		// plan the step
		c.err = c.PlanStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		// poll outputs
		opEnv, maskEnv, c.err = c.outputs.poll(ctx, c.OutputCtn)
		if c.err != nil {
			return fmt.Errorf("unable to exec outputs container: %w", c.err)
		}

		opEnv = outputs.Sanitize(_step, opEnv)
		maskEnv = outputs.Sanitize(_step, maskEnv)

		// merge env from outputs
		//
		//nolint:errcheck // only errors with empty environment input, which does not matter here
		_step.MergeEnv(opEnv)

		// merge env from masked outputs
		//
		//nolint:errcheck // only errors with empty environment input, which does not matter here
		_step.MergeEnv(maskEnv)

		// execute the step
		c.err = c.ExecStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to execute step: %w", c.err)
		}
	}

	// create an error group with the context for each stage
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#WithContext
	stages, stageCtx := errgroup.WithContext(ctx)
	// create a map to track the progress of each stage
	stageMap := new(sync.Map)

	// iterate through each stage in the pipeline
	for _, _stage := range c.pipeline.Stages {
		// TODO: remove hardcoded reference
		if _stage.Name == "init" {
			continue
		}

		// https://golang.org/doc/faq#closures_and_goroutines
		stage := _stage

		// create a new channel for each stage in the map
		stageMap.Store(stage.Name, make(chan error))

		// spawn errgroup routine for the stage
		//
		// https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Go
		stages.Go(func() error {
			// plan the stage
			c.err = c.PlanStage(stageCtx, stage, stageMap)
			if c.err != nil {
				return fmt.Errorf("unable to plan stage: %w", c.err)
			}

			// execute the stage
			c.err = c.ExecStage(stageCtx, stage, stageMap)
			if c.err != nil {
				return fmt.Errorf("unable to execute stage: %w", c.err)
			}

			return nil
		})
	}

	// wait for the stages to complete or return an error
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Wait
	c.err = stages.Wait()
	if c.err != nil {
		return fmt.Errorf("unable to wait for stages: %w", c.err)
	}

	return c.err
}

// StreamBuild receives a StreamRequest and then
// runs StreamService or StreamStep in a goroutine.
func (c *client) StreamBuild(ctx context.Context) error {
	// create an error group with the parent context
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#WithContext
	streams, streamCtx := errgroup.WithContext(ctx)

	defer func() {
		fmt.Fprintln(c.stdout, "waiting for stream functions to return")

		err := streams.Wait()
		if err != nil {
			fmt.Fprintln(c.stdout, "error in a stream request:", err)
		}

		fmt.Fprintln(c.stdout, "all stream functions have returned")
	}()

	// allow the runtime to do log/event streaming setup at build-level
	streams.Go(func() error {
		// If needed, the runtime should handle synchronizing with
		// AssembleBuild which runs concurrently with StreamBuild.
		return c.Runtime.StreamBuild(streamCtx, c.pipeline)
	})

	for {
		select {
		case req := <-c.streamRequests:
			streams.Go(func() error {
				fmt.Fprintf(c.stdout, "[%s: %s] > Streaming container '%s'...\n", req.Key, req.Container.Name, req.Container.ID)

				err := req.Stream(streamCtx, req.Container)
				if err != nil {
					fmt.Fprintln(c.stdout, "error streaming:", err)
				}

				return nil
			})
		case <-ctx.Done():
			// build done or canceled
			return nil
		}
	}
}

// DestroyBuild cleans up the build after execution.
func (c *client) DestroyBuild(ctx context.Context) error {
	var err error

	defer func() {
		// remove the runtime build for the pipeline
		err = c.Runtime.RemoveBuild(ctx, c.pipeline)
		if err != nil {
			// output the error information to stdout
			fmt.Fprintln(c.stdout, "unable to destroy runtime build:", err)
		}
	}()

	// destroy the steps for the pipeline
	for _, _step := range c.pipeline.Steps {
		// TODO: remove hardcoded reference
		if _step.Name == "init" {
			continue
		}

		// destroy the step
		err = c.DestroyStep(ctx, _step)
		if err != nil {
			// output the error information to stdout
			fmt.Fprintln(c.stdout, "unable to destroy step:", err)
		}
	}

	// destroy the stages for the pipeline
	for _, _stage := range c.pipeline.Stages {
		// TODO: remove hardcoded reference
		if _stage.Name == "init" {
			continue
		}

		// destroy the stage
		err = c.DestroyStage(ctx, _stage)
		if err != nil {
			// output the error information to stdout
			fmt.Fprintln(c.stdout, "unable to destroy stage:", err)
		}
	}

	// destroy the services for the pipeline
	for _, _service := range c.pipeline.Services {
		// destroy the service
		err = c.DestroyService(ctx, _service)
		if err != nil {
			// output the error information to stdout
			fmt.Fprintln(c.stdout, "unable to destroy service:", err)
		}
	}

	// destroy output container
	err = c.outputs.destroy(ctx, c.OutputCtn)
	if err != nil {
		fmt.Fprintln(c.stdout, "unable to destroy output container:", err)
	}

	// remove the runtime volume for the pipeline
	err = c.Runtime.RemoveVolume(ctx, c.pipeline)
	if err != nil {
		// output the error information to stdout
		fmt.Fprintln(c.stdout, "unable to destroy runtime volume:", err)
	}

	// remove the runtime network for the pipeline
	err = c.Runtime.RemoveNetwork(ctx, c.pipeline)
	if err != nil {
		// output the error information to stdout
		fmt.Fprintln(c.stdout, "unable to destroy runtime network:", err)
	}

	return err
}
