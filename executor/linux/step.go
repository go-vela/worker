// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/step"
	"golang.org/x/sync/errgroup"
)

// CreateStep configures the step for execution.
func (c *client) CreateStep(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := c.logger.WithField("step", ctn.Name)

	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	logger.Debug("setting up container")
	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// create a library step object to facilitate injecting environment as early as possible
	// (PlanStep is too late to inject environment vars for the kubernetes runtime).
	_step := c.newLibraryStep(ctn)
	_step.SetStatus(constants.StatusPending)

	// update the step container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Environment
	err = step.Environment(ctn, c.build, c.repo, _step, c.Version)
	if err != nil {
		return err
	}

	logger.Debug("escaping newlines in secrets")
	escapeNewlineSecrets(c.Secrets)

	logger.Debug("injecting secrets")
	// inject secrets for container
	err = injectSecrets(ctn, c.Secrets)
	if err != nil {
		return err
	}

	logger.Debug("substituting container configuration")
	// substitute container configuration
	//
	// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.Substitute
	err = ctn.Substitute()
	if err != nil {
		return fmt.Errorf("unable to substitute container configuration")
	}

	return nil
}

// newLibraryStep creates a library step object.
func (c *client) newLibraryStep(ctn *pipeline.Container) *library.Step {
	_step := new(library.Step)
	_step.SetName(ctn.Name)
	_step.SetNumber(ctn.Number)
	_step.SetImage(ctn.Image)
	_step.SetStage(ctn.Environment["VELA_STEP_STAGE"])
	_step.SetHost(c.build.GetHost())
	_step.SetRuntime(c.build.GetRuntime())
	_step.SetDistribution(c.build.GetDistribution())
	return _step
}

// PlanStep prepares the step for execution.
func (c *client) PlanStep(ctx context.Context, ctn *pipeline.Container) error {
	var err error

	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := c.logger.WithField("step", ctn.Name)

	// create the library step object
	_step := c.newLibraryStep(ctn)
	_step.SetStatus(constants.StatusRunning)
	_step.SetStarted(time.Now().UTC().Unix())

	logger.Debug("uploading step state")
	// send API call to update the step
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.Update
	_step, _, err = c.Vela.Step.Update(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), _step)
	if err != nil {
		return err
	}

	// update the step container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Environment
	err = step.Environment(ctn, c.build, c.repo, _step, c.Version)
	if err != nil {
		return err
	}

	// add a step to a map
	c.steps.Store(ctn.ID, _step)

	// get the step log here
	logger.Debug("retrieve step log")
	// send API call to capture the step log
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#LogService.GetStep
	_log, _, err := c.Vela.Log.GetStep(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), _step.GetNumber())
	if err != nil {
		return err
	}

	// add a step log to a map
	c.stepLogs.Store(ctn.ID, _log)

	return nil
}

// ExecStep runs a step.
func (c *client) ExecStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := c.logger.WithField("step", ctn.Name)

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
	defer func() { step.Snapshot(ctn, c.build, c.Vela, c.logger, c.repo, _step) }()

	logger.Debug("running container")
	// run the runtime container
	err = c.Runtime.RunContainer(ctx, ctn, c.pipeline)
	if err != nil {
		return err
	}

	// create an error group with the parent context
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#WithContext
	logs, logCtx := errgroup.WithContext(ctx)

	logs.Go(func() error {
		logger.Debug("streaming logs for container")
		// stream logs from container
		err := c.StreamStep(logCtx, ctn)
		if err != nil {
			logger.Error(err)
		}

		return nil
	})

	// do not wait for detached containers
	if ctn.Detach {
		return nil
	}

	logger.Debug("waiting for container")
	// wait for the runtime container
	err = c.Runtime.WaitContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("inspecting container")
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

	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := c.logger.WithField("step", ctn.Name)

	// load the logs for the step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(ctn, &c.stepLogs)
	if err != nil {
		return err
	}

	// nolint: dupl // ignore similar code
	defer func() {
		// tail the runtime container
		rc, err := c.Runtime.TailContainer(ctx, ctn)
		if err != nil {
			logger.Errorf("unable to tail container output for upload: %v", err)

			return
		}
		defer rc.Close()

		// read all output from the runtime container
		data, err := ioutil.ReadAll(rc)
		if err != nil {
			logger.Errorf("unable to read container output for upload: %v", err)

			return
		}

		// overwrite the existing log with all bytes
		//
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.SetData
		_log.SetData(data)

		logger.Debug("uploading logs")
		// send API call to update the logs for the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#LogService.UpdateStep
		_, _, err = c.Vela.Log.UpdateStep(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), ctn.Number, _log)
		if err != nil {
			logger.Errorf("unable to upload container logs: %v", err)
		}
	}()

	logger.Debug("tailing container")
	// tail the runtime container
	rc, err := c.Runtime.TailContainer(ctx, ctn)
	if err != nil {
		return err
	}
	defer rc.Close()

	// set the timeout to the repo timeout
	// to ensure the stream is not cut off
	c.Vela.SetTimeout(time.Minute * time.Duration(c.repo.GetTimeout()))

	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.Stream
	_, err = c.Vela.Step.Stream(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), ctn.Number, rc)
	if err != nil {
		logger.Errorf("unable to stream logs: %v", err)
	}

	logger.Info("finished streaming logs")

	return nil
}

// DestroyStep cleans up steps after execution.
func (c *client) DestroyStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := c.logger.WithField("step", ctn.Name)

	// load the step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_step, err := step.Load(ctn, &c.steps)
	if err != nil {
		// create the step from the container
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#StepFromContainer
		_step = library.StepFromContainer(ctn)
	}

	// defer an upload of the step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Upload
	defer func() { step.Upload(ctn, c.build, c.Vela, logger, c.repo, _step) }()

	logger.Debug("inspecting container")
	// inspect the runtime container
	err = c.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("removing container")
	// remove the runtime container
	err = c.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}
