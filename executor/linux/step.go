// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/internal/step"
)

// CreateStep configures the step for execution.
func (c *client) CreateStep(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("step", ctn.Name)

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
	_step := library.StepFromBuildContainer(c.build, ctn)

	// update the step container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Environment
	err = step.Environment(ctn, c.build, c.repo, _step, c.Version)
	if err != nil {
		return err
	}

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

// PlanStep prepares the step for execution.
func (c *client) PlanStep(ctx context.Context, ctn *pipeline.Container) error {
	var err error

	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("step", ctn.Name)

	// create the library step object
	_step := library.StepFromBuildContainer(c.build, ctn)
	_step.SetStatus(constants.StatusRunning)
	_step.SetStarted(time.Now().UTC().Unix())

	logger.Debug("uploading step state")
	// send API call to update the step
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#StepService.Update
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
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.GetStep
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
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("step", ctn.Name)

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
	defer func() { step.Snapshot(ctn, c.build, c.Vela, c.Logger, c.repo, _step) }()

	logger.Debug("running container")
	// run the runtime container
	err = c.Runtime.RunContainer(ctx, ctn, c.pipeline)
	if err != nil {
		// set step status to error and step error
		_step.SetStatus(constants.StatusError)
		_step.SetError(err.Error())

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
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("step", ctn.Name)

	// load the logs for the step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(ctn, &c.stepLogs)
	if err != nil {
		return err
	}

	existingLog := *_log

	secretValues := getSecretValues(ctn)

	defer func() {
		// tail the runtime container
		rc, err := c.Runtime.TailContainer(ctx, ctn)
		if err != nil {
			logger.Errorf("unable to tail container output for upload: %v", err)

			return
		}
		defer rc.Close()

		// read all output from the runtime container
		data, err := io.ReadAll(rc)
		if err != nil {
			logger.Errorf("unable to read container output for upload: %v", err)

			return
		}

		// don't attempt last upload if log size exceeded
		if c.maxLogSize > 0 && uint(len(data)) >= c.maxLogSize {
			logger.Trace("maximum log size reached")

			return
		}

		// overwrite the existing log with all bytes
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#Log.SetData
		_log.SetData(append(existingLog.GetData(), data...))

		// mask secrets in the log data
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#Log.MaskData
		_log.MaskData(secretValues)

		logger.Debug("uploading logs")
		// send API call to update the logs for the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.UpdateStep
		_, err = c.Vela.Log.UpdateStep(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), ctn.Number, _log)
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

	// create new buffer for uploading logs
	logs := new(bytes.Buffer)

	// create new channel for processing logs
	done := make(chan bool)

	go func() {
		logger.Debug("polling logs for container")

		// spawn "infinite" loop that will upload logs
		// from the buffer until the channel is closed
		for {
			// sleep for "1s" before attempting to upload logs
			time.Sleep(1 * time.Second)

			// create a non-blocking select to check if the channel is closed
			select {
			// after repo timeout of idle (no response) end the stream
			//
			// this is a safety mechanism
			case <-time.After(time.Duration(c.repo.GetTimeout()) * time.Minute):
				logger.Tracef("repo timeout of %d exceeded", c.repo.GetTimeout())

				return
			// channel is closed
			case <-done:
				logger.Trace("channel closed for polling container logs")

				// return out of the go routine
				return
			// channel is not closed
			default:
				// update the existing log with the new bytes if there is new data to add
				if len(logs.Bytes()) > 0 {
					logger.Trace(logs.String())

					// update the existing log with the new bytes
					//
					// https://pkg.go.dev/github.com/go-vela/types/library#Log.AppendData
					_log.AppendData(logs.Bytes())
					// mask secrets within the logs before updating database
					//
					// https://pkg.go.dev/github.com/go-vela/types/library#Log.MaskData
					_log.MaskData(secretValues)

					logger.Debug("appending logs")
					// send API call to append the logs for the step
					//
					// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogStep.UpdateStep
					_, err := c.Vela.Log.UpdateStep(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), ctn.Number, _log)
					if err != nil {
						logger.Error(err)
					}

					// flush the buffer of logs
					logs.Reset()
				}

				// check whether we've reached the maximum log size
				if c.maxLogSize > 0 && uint(len(_log.GetData())) >= c.maxLogSize {
					logger.Trace("maximum log size reached")

					return
				}
			}
		}
	}()

	// create new scanner from the container output
	scanner := bufio.NewScanner(rc)

	// scan entire container output
	for scanner.Scan() {
		// write all the logs from the scanner
		logs.Write(append(scanner.Bytes(), []byte("\n")...))
	}

	logger.Info("finished streaming logs")

	// close channel to stop processing logs
	close(done)

	return scanner.Err()
}

// DestroyStep cleans up steps after execution.
func (c *client) DestroyStep(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	// update engine logger with step metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("step", ctn.Name)

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

// getSecretValues is a helper function that creates a slice of
// secret values that will be used to mask secrets in logs before
// updating the database.
func getSecretValues(ctn *pipeline.Container) []string {
	secretValues := []string{}
	// gather secrets' values from the environment map for masking
	for _, secret := range ctn.Secrets {
		// capture secret from environment
		s, ok := ctn.Environment[strings.ToUpper(secret.Target)]
		if !ok {
			continue
		}
		// handle multi line secrets from files
		s = strings.ReplaceAll(s, "\n", " ")

		secretValues = append(secretValues, strings.TrimSuffix(s, " "))
	}

	return secretValues
}
