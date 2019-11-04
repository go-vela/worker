// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/worker/version"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/drone/envsubst"
	"github.com/sirupsen/logrus"
)

// CreateStep prepares the step for execution.
func (c *client) CreateStep(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"step": ctn.Name,
	})

	logger.Debug("setting up container")
	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// inject default workspace
	if len(ctn.Directory) == 0 {
		ctn.Directory = fmt.Sprintf("/home/%s", c.pipeline.ID)
	}

	ctn.Environment["BUILD_HOST"] = c.Hostname
	ctn.Environment["VELA_HOST"] = c.Hostname
	ctn.Environment["VELA_VERSION"] = version.Version.String()
	// TODO: This should not be hardcoded
	ctn.Environment["VELA_RUNTIME"] = "docker"
	ctn.Environment["VELA_DISTRIBUTION"] = "linux"

	logger.Debug("injecting secrets")
	// inject secrets for step
	err = injectSecrets(ctn, c.Secrets)
	if err != nil {
		return err
	}

	logger.Debug("marshaling configuration")
	// marshal container configuration
	body, err := json.Marshal(ctn)
	if err != nil {
		return fmt.Errorf("unable to marshal configuration: %v", err)
	}

	// create substitute function
	subFunc := func(name string) string {
		env := ctn.Environment[name]
		if strings.Contains(env, "\n") {
			env = fmt.Sprintf("%q", env)
		}
		return env
	}

	logger.Debug("substituting environment")
	// substitute the environment variables
	subStep, err := envsubst.Eval(string(body), subFunc)
	if err != nil {
		return fmt.Errorf("unable to substitute environment variables: %v", err)
	}

	logger.Debug("unmarshaling configuration")
	// unmarshal container configuration
	err = json.Unmarshal([]byte(subStep), ctn)
	if err != nil {
		return fmt.Errorf("unable to unmarshal configuration: %v", err)
	}

	return nil
}

// PlanStep defines a function that prepares the step for execution.
func (c *client) PlanStep(ctx context.Context, ctn *pipeline.Container) error {
	var err error
	b := c.build
	r := c.repo

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"step": ctn.Name,
	})

	// update the engine step object
	c.step = &library.Step{
		Name:         vela.String(ctn.Name),
		Number:       vela.Int(ctn.Number),
		Status:       vela.String(constants.StatusRunning),
		Started:      vela.Int64(time.Now().UTC().Unix()),
		Host:         vela.String(ctn.Environment["VELA_HOST"]),
		Runtime:      vela.String(ctn.Environment["VELA_RUNTIME"]),
		Distribution: vela.String(ctn.Environment["VELA_DISTRIBUTION"]),
	}

	logger.Debug("uploading step state")
	// send API call to update the step
	c.step, _, err = c.Vela.Step.Update(r.GetOrg(), r.GetName(), b.GetNumber(), c.step)
	if err != nil {
		return err
	}
	c.step.Status = vela.String(constants.StatusSuccess)

	// get the step log here
	logger.Debug("retrieve step log")
	// send API call to capture the step log
	c.stepLog, _, err = c.Vela.Log.GetStep(r.GetOrg(), r.GetName(), b.GetNumber(), c.step.GetNumber())
	if err != nil {
		return err
	}

	return nil
}

// ExecStep runs a step.
func (c *client) ExecStep(ctx context.Context, ctn *pipeline.Container) error {
	b := c.build
	r := c.repo

	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"step": ctn.Name,
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
				c.stepLog.Data = vela.Bytes(append(c.stepLog.GetData(), logs.Bytes()...))

				logger.Debug("appending logs")
				// send API call to update the logs for the step
				c.stepLog, _, err = c.Vela.Log.UpdateStep(r.GetOrg(), r.GetName(), b.GetNumber(), ctn.Number, c.stepLog)
				if err != nil {
					return err
				}

				// flush the buffer of logs
				logs.Reset()
			}
		}
		logger.Trace(logs.String())

		// update the existing log with the last bytes
		c.stepLog.Data = vela.Bytes(append(c.stepLog.GetData(), logs.Bytes()...))

		logger.Debug("uploading logs")
		// send API call to update the logs for the step
		c.stepLog, _, err = c.Vela.Log.UpdateStep(r.GetOrg(), r.GetName(), b.GetNumber(), ctn.Number, c.stepLog)
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

// DestroyStep cleans up steps after execution.
func (c *client) DestroyStep(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with extra metadata
	logger := c.logger.WithFields(logrus.Fields{
		"step": ctn.Name,
	})

	logger.Debug("removing container")
	// remove the runtime container
	err := c.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}
