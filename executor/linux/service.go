// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-vela/worker/version"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/sirupsen/logrus"
)

// CreateService configures the service for execution.
func (c *client) CreateService(ctx context.Context, ctn *pipeline.Container) error {
	var err error

	b := c.build
	r := c.repo

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	// update the engine service object
	s := new(library.Service)
	s.SetName(ctn.Name)
	s.SetNumber(ctn.Number)
	s.SetStatus(constants.StatusRunning)
	s.SetStarted(time.Now().UTC().Unix())

	logger.Debug("uploading service state")
	// send API call to update the service
	s, _, err = c.Vela.Svc.Update(r.GetOrg(), r.GetName(), b.GetNumber(), s)
	if err != nil {
		return err
	}

	s.SetStatus(constants.StatusSuccess)

	// add a service to a map
	c.services.Store(ctn.ID, s)

	// get the service log here
	logger.Debug("retrieve service log")
	// send API call to capture the service log
	l, _, err := c.Vela.Log.GetService(r.GetOrg(), r.GetName(), b.GetNumber(), s.GetNumber())
	if err != nil {
		return err
	}

	// add a step log to a map
	c.serviceLogs.Store(ctn.ID, l)

	return nil
}

// PlanService prepares the service for execution.
func (c *client) PlanService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	ctn.Environment["BUILD_HOST"] = c.Hostname
	ctn.Environment["VELA_HOST"] = c.Hostname
	ctn.Environment["VELA_VERSION"] = version.Version.String()
	// TODO: remove hardcoded reference
	ctn.Environment["VELA_RUNTIME"] = "docker"
	ctn.Environment["VELA_DISTRIBUTION"] = "linux"

	logger.Debug("setting up container")
	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// ExecService runs a service.
func (c *client) ExecService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	logger.Debug("running container")
	// run the runtime container
	err := c.Runtime.RunContainer(ctx, c.pipeline, ctn)
	if err != nil {
		return err
	}

	go func() {
		logger.Debug("stream logs for container")
		// stream logs from container
		err := c.StreamService(ctx, ctn)
		if err != nil {
			logrus.Error(err)
		}
	}()

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

// StreamService tails the output for a service.
func (c *client) StreamService(ctx context.Context, ctn *pipeline.Container) error {
	// TODO: remove hardcoded reference
	if ctn.Name == "init" {
		return nil
	}

	b := c.build
	r := c.repo

	result, ok := c.stepLogs.Load(ctn.ID)
	if !ok {
		return fmt.Errorf("unable to get step log from client")
	}

	l := result.(*library.Log)

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"step": ctn.Name,
	})

	// create new buffer for uploading logs
	logs := new(bytes.Buffer)

	logger.Debug("tailing container")
	// tail the runtime container
	rc, err := c.Runtime.TailContainer(ctx, ctn)
	if err != nil {
		return err
	}
	defer rc.Close()

	// create new scanner from the container output
	scanner := bufio.NewScanner(rc)

	// scan entire container output
	for scanner.Scan() {
		// write all the logs from the scanner
		logs.Write(append(scanner.Bytes(), []byte("\n")...))

		// if we have at least 1000 bytes in our buffer
		if logs.Len() > 1000 {
			logger.Trace(logs.String())

			// update the existing log with the new bytes
			l.SetData(append(l.GetData(), logs.Bytes()...))

			logger.Debug("appending logs")
			// send API call to append the logs for the service
			l, _, err = c.Vela.Log.UpdateService(r.GetOrg(), r.GetName(), b.GetNumber(), ctn.Number, l)
			if err != nil {
				return err
			}

			// flush the buffer of logs
			logs.Reset()
		}
	}
	logger.Trace(logs.String())

	// update the existing log with the last bytes
	l.SetData(append(l.GetData(), logs.Bytes()...))

	logger.Debug("uploading logs")
	// send API call to update the logs for the service
	_, _, err = c.Vela.Log.UpdateService(r.GetOrg(), r.GetName(), b.GetNumber(), ctn.Number, l)
	if err != nil {
		return err
	}

	return nil
}

// DestroyService cleans up services after execution.
func (c *client) DestroyService(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	logger.Debug("inspecting container")
	// inspect the runtime container
	err := c.Runtime.InspectContainer(ctx, ctn)
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
