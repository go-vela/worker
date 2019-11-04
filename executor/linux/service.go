// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"bufio"
	"bytes"
	"context"
	"time"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// CreateService prepares the service for execution.
func (c *client) CreateService(ctx context.Context, ctn *pipeline.Container) error {

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	logger.Debug("setting up container")
	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// PlanService defines a function that prepares the service for execution.
func (c *client) PlanService(ctx context.Context, ctn *pipeline.Container) error {
	var err error
	b := c.build
	r := c.repo

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	// update the engine service object
	c.service = &library.Service{
		Name:    vela.String(ctn.Name),
		Number:  vela.Int(ctn.Number),
		Status:  vela.String(constants.StatusRunning),
		Started: vela.Int64(time.Now().UTC().Unix()),
	}

	logger.Debug("uploading service state")
	// send API call to update the service
	_, _, err = c.Vela.Svc.Update(r.GetOrg(), r.GetName(), b.GetNumber(), c.service)
	if err != nil {
		return err
	}
	c.service.Status = vela.String(constants.StatusSuccess)

	// get the service log here
	logger.Debug("retrieve service log")
	// send API call to capture the service log
	c.serviceLog, _, err = c.Vela.Log.GetService(r.GetOrg(), r.GetName(), b.GetNumber(), c.service.GetNumber())
	if err != nil {
		return err
	}

	return nil
}

// ExecService runs a service.
func (c *client) ExecService(ctx context.Context, ctn *pipeline.Container) error {
	b := c.build
	r := c.repo

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"service": ctn.Name,
	})

	// run the container in a detached state
	if ctn.Detach {
		logger.Debug("running container in detach mode")
		// run the runtime container
		err := c.Runtime.RunContainer(ctx, c.pipeline, ctn)
		if err != nil {
			return err
		}

		return nil
	}

	logger.Debug("running container")
	// run the runtime container
	err := c.Runtime.RunContainer(ctx, c.pipeline, ctn)
	if err != nil {
		return err
	}

	// create new buffer for uploading logs
	logs := new(bytes.Buffer)
	go func() error {
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
				c.serviceLog.Data = vela.Bytes(append(c.serviceLog.GetData(), logs.Bytes()...))

				logger.Debug("appending logs")
				c.serviceLog, _, err = c.Vela.Log.UpdateService(r.GetOrg(), r.GetName(), b.GetNumber(), ctn.Number, c.serviceLog)
				if err != nil {
					return err
				}

				// flush the buffer of logs
				logs.Reset()
			}
		}
		logger.Trace(logs.String())

		// update the existing log with the last bytes
		c.serviceLog.Data = vela.Bytes(append(c.serviceLog.GetData(), logs.Bytes()...))

		logger.Debug("uploading logs")
		// send API call to update the logs for the service
		c.serviceLog, _, err = c.Vela.Log.UpdateService(r.GetOrg(), r.GetName(), b.GetNumber(), ctn.Number, c.serviceLog)
		if err != nil {
			return err
		}

		return nil
	}()

	logger.Debug("waiting for container")
	// wait for the runtime container
	err = c.Runtime.WaitContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("inspecting container")
	// inspect the runtime container
	err = c.Runtime.InfoContainer(ctx, ctn)
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
	err := c.Runtime.InfoContainer(ctx, ctn)
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
