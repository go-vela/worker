// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/internal/service"
)

// CreateService configures the service for execution.
func (c *client) CreateService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with service metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("service", ctn.Name)

	logger.Debug("setting up container")
	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// update the service container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Environment
	err = service.Environment(ctn, c.build, c.repo, nil, c.Version)
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

	logger.Debug("injecting non-substituted secrets")
	// inject no-substitution secrets for container
	err = injectSecrets(ctn, c.NoSubSecrets)
	if err != nil {
		return err
	}

	return nil
}

// PlanService prepares the service for execution.
func (c *client) PlanService(ctx context.Context, ctn *pipeline.Container) error {
	var err error

	// update engine logger with service metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("service", ctn.Name)

	// create the library service object
	_service := new(library.Service)
	_service.SetName(ctn.Name)
	_service.SetNumber(ctn.Number)
	_service.SetImage(ctn.Image)
	_service.SetStatus(constants.StatusRunning)
	_service.SetStarted(time.Now().UTC().Unix())
	_service.SetHost(c.build.GetHost())
	_service.SetRuntime(c.build.GetRuntime())
	_service.SetDistribution(c.build.GetDistribution())

	logger.Debug("uploading service state")
	// send API call to update the service
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#SvcService.Update
	_service, _, err = c.Vela.Svc.Update(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), _service)
	if err != nil {
		return err
	}

	// update the service container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Environment
	err = service.Environment(ctn, c.build, c.repo, _service, c.Version)
	if err != nil {
		return err
	}

	// add a service to a map
	c.services.Store(ctn.ID, _service)

	// get the service log here
	logger.Debug("retrieve service log")
	// send API call to capture the service log
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.GetService
	_log, _, err := c.Vela.Log.GetService(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), _service.GetNumber())
	if err != nil {
		return err
	}

	// add a service log to a map
	c.serviceLogs.Store(ctn.ID, _log)

	return nil
}

// ExecService runs a service.
func (c *client) ExecService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with service metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("service", ctn.Name)

	// load the service from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Load
	_service, err := service.Load(ctn, &c.services)
	if err != nil {
		return err
	}

	// defer taking a snapshot of the service
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Snapshot
	defer func() { service.Snapshot(ctn, c.build, c.Vela, c.Logger, c.repo, _service) }()

	logger.Debug("running container")
	// run the runtime container
	err = c.Runtime.RunContainer(ctx, ctn, c.pipeline)
	if err != nil {
		return err
	}

	// trigger StreamService goroutine with logging context
	c.streamRequests <- message.StreamRequest{
		Key:       "service",
		Stream:    c.StreamService,
		Container: ctn,
	}

	return nil
}

// StreamService tails the output for a service.
func (c *client) StreamService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with service metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("service", ctn.Name)

	// load the logs for the service from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#LoadLogs
	_log, err := service.LoadLogs(ctn, &c.serviceLogs)
	if err != nil {
		return err
	}

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
		_log.SetData(data)

		logger.Debug("uploading logs")
		// send API call to update the logs for the service
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.UpdateService
		_, err = c.Vela.Log.UpdateService(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), ctn.Number, _log)
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

					logger.Debug("appending logs")
					// send API call to append the logs for the service
					//
					// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.UpdateService
					_, err = c.Vela.Log.UpdateService(c.repo.GetOrg(), c.repo.GetName(), c.build.GetNumber(), ctn.Number, _log)
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

// DestroyService cleans up services after execution.
func (c *client) DestroyService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with service metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := c.Logger.WithField("service", ctn.Name)

	// load the service from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Load
	_service, err := service.Load(ctn, &c.services)
	if err != nil {
		// create the service from the container
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#ServiceFromContainerEnvironment
		_service = library.ServiceFromContainerEnvironment(ctn)
	}

	// defer an upload of the service
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#LoaUploadd
	defer func() { service.Upload(ctn, c.build, c.Vela, logger, c.repo, _service) }()

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
