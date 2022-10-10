// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/worker/internal/build"
	context2 "github.com/go-vela/worker/internal/context"
	"github.com/go-vela/worker/internal/step"
)

// CreateBuild configures the build for execution.
func (c *client) CreateBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, c.Vela, c.err, c.Logger, c.repo) }()

	// update the build fields
	c.build.SetStatus(constants.StatusRunning)
	c.build.SetStarted(time.Now().UTC().Unix())
	c.build.SetHost(c.Hostname)
	c.build.SetDistribution(c.Driver())
	c.build.SetRuntime(c.Runtime.Driver())

	c.Logger.Info("uploading build state")
	// send API call to update the build
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#BuildService.Update
	//nolint:contextcheck // ignore passing context
	c.build, _, c.err = c.Vela.Build.Update(c.repo.GetOrg(), c.repo.GetName(), c.build)
	if c.err != nil {
		return fmt.Errorf("unable to upload build state: %w", c.err)
	}

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

	c.Logger.Infof("creating %s step", c.init.Name)
	// create the step
	c.err = c.CreateStep(ctx, c.init)
	if c.err != nil {
		return fmt.Errorf("unable to create %s step: %w", c.init.Name, c.err)
	}

	c.Logger.Infof("planning %s step", c.init.Name)
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
	defer func() { build.Snapshot(c.build, c.Vela, c.err, c.Logger, c.repo) }()

	// load the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_init, err := step.Load(c.init, &c.steps)
	if err != nil {
		return err
	}

	// load the logs for the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(c.init, &c.stepLogs)
	if err != nil {
		return err
	}

	// defer taking a snapshot of the init step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#SnapshotInit
	defer func() { step.SnapshotInit(c.init, c.build, c.Vela, c.Logger, c.repo, _init, _log) }()

	c.Logger.Info("creating network")
	// create the runtime network for the pipeline
	c.err = c.Runtime.CreateNetwork(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to create network: %w", c.err)
	}

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Inspecting runtime network...\n"))

	// inspect the runtime network for the pipeline
	network, err := c.Runtime.InspectNetwork(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect network: %w", err)
	}

	// update the init log with network information
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData(network)

	c.Logger.Info("creating volume")
	// create the runtime volume for the pipeline
	c.err = c.Runtime.CreateVolume(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to create volume: %w", c.err)
	}

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Inspecting runtime volume...\n"))

	// inspect the runtime volume for the pipeline
	volume, err := c.Runtime.InspectVolume(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect volume: %w", err)
	}

	// update the init log with volume information
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData(volume)

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Preparing secrets...\n"))

	// iterate through each secret provided in the pipeline
	for _, secret := range c.pipeline.Secrets {
		// ignore pulling secrets coming from plugins
		if !secret.Origin.Empty() {
			continue
		}

		c.Logger.Infof("pulling %s %s secret %s", secret.Engine, secret.Type, secret.Name)

		//nolint:contextcheck // ignore passing context
		s, err := c.secret.pull(secret)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to pull secrets: %w", err)
		}

		_log.AppendData([]byte(
			fmt.Sprintf("$ vela view secret --secret.engine %s --secret.type %s --org %s --repo %s --name %s \n",
				secret.Engine, secret.Type, s.GetOrg(), s.GetRepo(), s.GetName())))

		sRaw, err := json.MarshalIndent(s.Sanitize(), "", " ")
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to decode secret: %w", err)
		}

		_log.AppendData(append(sRaw, "\n"...))

		// add secret to the map
		c.Secrets[secret.Name] = s
	}

	return nil
}

// AssembleBuild prepares the containers within a build for execution.
//
//nolint:funlen // ignore function length due to comments and logging messages
func (c *client) AssembleBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, c.Vela, c.err, c.Logger, c.repo) }()

	// load the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_init, err := step.Load(c.init, &c.steps)
	if err != nil {
		return err
	}

	// load the logs for the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(c.init, &c.stepLogs)
	if err != nil {
		return err
	}

	// defer an upload of the init step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Upload
	defer func() { step.Upload(c.init, c.build, c.Vela, c.Logger, c.repo, _init) }()

	defer func() {
		c.Logger.Infof("uploading %s step logs", c.init.Name)
		// send API call to update the logs for the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#LogService.UpdateStep
		_log, _, err = c.Vela.Log.UpdateStep(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), c.init.Number, _log)
		if err != nil {
			c.Logger.Errorf("unable to upload %s logs: %v", c.init.Name, err)
		}
	}()

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Preparing service images...\n"))

	// create the services for the pipeline
	for _, s := range c.pipeline.Services {
		// TODO: remove this; but we need it for tests
		s.Detach = true

		c.Logger.Infof("creating %s service", s.Name)
		// create the service
		c.err = c.CreateService(ctx, s)
		if c.err != nil {
			return fmt.Errorf("unable to create %s service: %w", s.Name, c.err)
		}

		c.Logger.Infof("inspecting %s service", s.Name)
		// inspect the service image
		image, err := c.Runtime.InspectImage(ctx, s)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s service: %w", s.Name, err)
		}

		// update the init log with service image info
		//
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
		_log.AppendData(image)
	}

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Preparing stage images...\n"))

	// create the stages for the pipeline
	for _, s := range c.pipeline.Stages {
		// TODO: remove hardcoded reference
		//
		//nolint:goconst // ignore making a constant for now
		if s.Name == "init" {
			continue
		}

		c.Logger.Infof("creating %s stage", s.Name)
		// create the stage
		c.err = c.CreateStage(ctx, s)
		if c.err != nil {
			return fmt.Errorf("unable to create %s stage: %w", s.Name, c.err)
		}
	}

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Preparing step images...\n"))

	// create the steps for the pipeline
	for _, s := range c.pipeline.Steps {
		// TODO: remove hardcoded reference
		if s.Name == "init" {
			continue
		}

		c.Logger.Infof("creating %s step", s.Name)
		// create the step
		c.err = c.CreateStep(ctx, s)
		if c.err != nil {
			return fmt.Errorf("unable to create %s step: %w", s.Name, c.err)
		}

		c.Logger.Infof("inspecting %s step", s.Name)
		// inspect the step image
		image, err := c.Runtime.InspectImage(ctx, s)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s step: %w", s.Name, c.err)
		}

		// update the init log with step image info
		//
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
		_log.AppendData(image)
	}

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Preparing secret images...\n"))

	// create the secrets for the pipeline
	for _, s := range c.pipeline.Secrets {
		// skip over non-plugin secrets
		if s.Origin.Empty() {
			continue
		}

		c.Logger.Infof("creating %s secret", s.Origin.Name)
		// create the service
		c.err = c.secret.create(ctx, s.Origin)
		if c.err != nil {
			return fmt.Errorf("unable to create %s secret: %w", s.Origin.Name, c.err)
		}

		c.Logger.Infof("inspecting %s secret", s.Origin.Name)
		// inspect the service image
		image, err := c.Runtime.InspectImage(ctx, s.Origin)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s secret: %w", s.Origin.Name, err)
		}

		// update the init log with secret image info
		//
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
		_log.AppendData(image)
	}

	// inspect the runtime build (eg a kubernetes pod) for the pipeline
	buildOutput, err := c.Runtime.InspectBuild(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect build: %w", err)
	}

	if len(buildOutput) > 0 {
		// update the init log with progress
		// (an empty value allows the runtime to opt out of providing this)
		//
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
		_log.AppendData(buildOutput)
	}

	// assemble runtime build just before any containers execute
	c.err = c.Runtime.AssembleBuild(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to assemble runtime build %s: %w", c.pipeline.ID, c.err)
	}

	// update the init log with progress
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
	_log.AppendData([]byte("> Executing secret images...\n"))

	c.Logger.Info("executing secret images")
	// execute the secret
	c.err = c.secret.exec(ctx, &c.pipeline.Secrets)
	if c.err != nil {
		return fmt.Errorf("unable to execute secret: %w", c.err)
	}

	return c.err
}

// ExecBuild runs a pipeline for a build.
func (c *client) ExecBuild(ctx context.Context) error {
	// defer an upload of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Upload
	defer func() { build.Upload(c.build, c.Vela, c.err, c.Logger, c.repo) }()

	// execute the services for the pipeline
	for _, _service := range c.pipeline.Services {
		c.Logger.Infof("planning %s service", _service.Name)
		// plan the service
		c.err = c.PlanService(ctx, _service)
		if c.err != nil {
			return fmt.Errorf("unable to plan service: %w", c.err)
		}

		c.Logger.Infof("executing %s service", _service.Name)
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
		if step.Skip(_step, c.build, c.repo) {
			continue
		}

		c.Logger.Infof("planning %s step", _step.Name)
		// plan the step
		c.err = c.PlanStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		c.Logger.Infof("executing %s step", _step.Name)
		// execute the step
		c.err = c.ExecStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to execute step: %w", c.err)
		}
	}

	// create an error group with the context for each stage
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#WithContext
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
		// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Go
		stages.Go(func() error {
			c.Logger.Infof("planning %s stage", stage.Name)
			// plan the stage
			c.err = c.PlanStage(stageCtx, stage, stageMap)
			if c.err != nil {
				return fmt.Errorf("unable to plan stage: %w", c.err)
			}

			c.Logger.Infof("executing %s stage", stage.Name)
			// execute the stage
			c.err = c.ExecStage(stageCtx, stage, stageMap)
			if c.err != nil {
				return fmt.Errorf("unable to execute stage: %w", c.err)
			}

			return nil
		})
	}

	c.Logger.Debug("waiting for stages completion")
	// wait for the stages to complete or return an error
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Wait
	c.err = stages.Wait()
	if c.err != nil {
		return fmt.Errorf("unable to wait for stages: %w", c.err)
	}

	return c.err
}

// StreamBuild receives a StreamRequest and then
// runs StreamService or StreamStep in a goroutine.
func (c *client) StreamBuild(ctx context.Context) error {
	// cancel streaming after a timeout once the build has finished
	delayedCtx, cancelStreaming := context2.WithDelayedCancelPropagation(ctx, c.logStreamingTimeout)
	defer cancelStreaming()

	// create an error group with the parent context
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#WithContext
	streams, streamCtx := errgroup.WithContext(delayedCtx)

	defer func() {
		c.Logger.Trace("waiting for stream functions to return")

		err := streams.Wait()
		if err != nil {
			c.Logger.Errorf("error in a stream request, %v", err)
		}

		cancelStreaming()

		c.Logger.Info("all stream functions have returned")
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
				// update engine logger with step metadata
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
				logger := c.Logger.WithField(req.Key, req.Container.Name)

				logger.Debugf("streaming %s container %s", req.Key, req.Container.ID)

				err := req.Stream(streamCtx, req.Container)
				if err != nil {
					logger.Error(err)
				}

				return nil
			})
		case <-delayedCtx.Done():
			c.Logger.Debug("streaming context canceled")
			// build done or canceled
			return nil
		}
	}
}

// DestroyBuild cleans up the build after execution.
func (c *client) DestroyBuild(ctx context.Context) error {
	var err error

	defer func() {
		c.Logger.Info("deleting runtime build")
		// remove the runtime build for the pipeline
		err = c.Runtime.RemoveBuild(ctx, c.pipeline)
		if err != nil {
			c.Logger.Errorf("unable to remove runtime build: %v", err)
		}
	}()

	// destroy the steps for the pipeline
	for _, _step := range c.pipeline.Steps {
		// TODO: remove hardcoded reference
		if _step.Name == "init" {
			continue
		}

		c.Logger.Infof("destroying %s step", _step.Name)
		// destroy the step
		err = c.DestroyStep(ctx, _step)
		if err != nil {
			c.Logger.Errorf("unable to destroy step: %v", err)
		}
	}

	// destroy the stages for the pipeline
	for _, _stage := range c.pipeline.Stages {
		// TODO: remove hardcoded reference
		if _stage.Name == "init" {
			continue
		}

		c.Logger.Infof("destroying %s stage", _stage.Name)
		// destroy the stage
		err = c.DestroyStage(ctx, _stage)
		if err != nil {
			c.Logger.Errorf("unable to destroy stage: %v", err)
		}
	}

	// destroy the services for the pipeline
	for _, _service := range c.pipeline.Services {
		c.Logger.Infof("destroying %s service", _service.Name)
		// destroy the service
		err = c.DestroyService(ctx, _service)
		if err != nil {
			c.Logger.Errorf("unable to destroy service: %v", err)
		}
	}

	// destroy the secrets for the pipeline
	for _, _secret := range c.pipeline.Secrets {
		// skip over non-plugin secrets
		if _secret.Origin.Empty() {
			continue
		}

		c.Logger.Infof("destroying %s secret", _secret.Name)
		// destroy the secret
		err = c.secret.destroy(ctx, _secret.Origin)
		if err != nil {
			c.Logger.Errorf("unable to destroy secret: %v", err)
		}
	}

	c.Logger.Info("deleting volume")
	// remove the runtime volume for the pipeline
	err = c.Runtime.RemoveVolume(ctx, c.pipeline)
	if err != nil {
		c.Logger.Errorf("unable to remove volume: %v", err)
	}

	c.Logger.Info("deleting network")
	// remove the runtime network for the pipeline
	err = c.Runtime.RemoveNetwork(ctx, c.pipeline)
	if err != nil {
		c.Logger.Errorf("unable to remove network: %v", err)
	}

	return err
}
