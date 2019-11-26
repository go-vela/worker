// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// CreateBuild prepares the build for execution.
func (c *client) CreateBuild(ctx context.Context) error {
	var err error
	b := c.build
	p := c.pipeline
	r := c.repo

	// update engine logger with extra metadata
	c.logger = c.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"repo":  r.GetFullName(),
	})

	defer func() {
		// NOTE: When an error occurs during a build that does not have to do
		// with a pipeline we should set build status to "error" not "failed"
		// because it is worker related and not build.
		if err != nil {
			// update the build fields
			b.SetFinished(time.Now().UTC().Unix())
			b.SetStatus(constants.StatusError)
			b.SetError(err.Error())

			c.logger.Info("uploading errored stated")
			// send API call to update the build
			_, _, err = c.Vela.Build.Update(r.GetOrg(), r.GetName(), b)
			if err != nil {
				c.logger.Errorf("unable to upload errored state: %w", err)
			}
		}
	}()

	// update the build fields
	b.SetStatus(constants.StatusRunning)
	b.SetStarted(time.Now().UTC().Unix())
	b.SetHost(c.Hostname)
	// TODO: This should not be hardcoded
	b.SetDistribution("linux")
	b.SetRuntime("docker")

	c.logger.Info("uploading build state")
	// send API call to update the build
	b, _, err = c.Vela.Build.Update(r.GetOrg(), r.GetName(), b)
	if err != nil {
		return fmt.Errorf("unable to upload start state: %w", err)
	}

	c.logger.Info("creating network")
	// create the runtime network for the pipeline
	err = c.Runtime.CreateNetwork(ctx, p)
	if err != nil {
		return fmt.Errorf("unable to create network: %w", err)
	}

	c.logger.Info("creating volume")
	// create the runtime volume for the pipeline
	err = c.Runtime.CreateVolume(ctx, p)
	if err != nil {
		return fmt.Errorf("unable to create volume: %w", err)
	}

	c.logger.Info("pulling secrets")
	// pull secrets for the pipeline
	err = c.PullSecret(ctx)
	if err != nil {
		return fmt.Errorf("unable to pull secrets: %w", err)
	}

	// create the services for the pipeline
	for _, s := range p.Services {
		// TODO: remove this; but we need it for tests
		s.Detach = true
		s.Pull = true

		c.logger.Infof("creating %s service", s.Name)
		// create the service
		err = c.CreateService(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to create service: %w", err)
		}
	}

	// create the stages for the pipeline
	for _, s := range p.Stages {
		c.logger.Infof("creating %s stage", s.Name)
		// create the stage
		err = c.CreateStage(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to create stage: %w", err)
		}
	}

	// create the steps for the pipeline
	for _, s := range p.Steps {
		c.logger.Infof("creating %s step", s.Name)
		// create the step
		err = c.CreateStep(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to create step: %w", err)
		}
	}

	b.SetStatus(constants.StatusSuccess)
	c.build = b

	return nil
}

// ExecBuild runs a pipeline for a build.
func (c *client) ExecBuild(ctx context.Context) error {
	var err error
	b := c.build
	p := c.pipeline
	r := c.repo

	defer func() {
		// NOTE: When an error occurs during a build that does not have to do
		// with a pipeline we should set build status to "error" not "failed"
		// because it is worker related and not build.
		if err != nil {
			// update the build fields
			b.SetFinished(time.Now().UTC().Unix())
			b.SetStatus(constants.StatusError)
			b.SetError(err.Error())

			c.logger.Info("uploading errored stated")
			// send API call to update the build
			_, _, err = c.Vela.Build.Update(r.GetOrg(), r.GetName(), b)
			if err != nil {
				c.logger.Errorf("unable to upload errored state: %w", err)
			}
		}
	}()

	// execute the services for the pipeline
	for _, s := range p.Services {
		c.logger.Infof("planning %s service", s.Name)
		// plan the service
		err = c.PlanService(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to plan service: %w", err)
		}

		c.logger.Infof("executing %s service", s.Name)
		// execute the service
		err = c.ExecService(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to execute service: %w", err)
		}
	}

	// execute the steps for the pipeline
	for _, s := range p.Steps {
		// check if the build status is successful
		if !strings.EqualFold(b.GetStatus(), constants.StatusSuccess) {
			// break out of loop to stop running steps
			break
		}

		c.logger.Infof("planning %s step", s.Name)
		// plan the step
		err = c.PlanStep(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to plan step: %w", err)
		}

		c.logger.Infof("executing %s step", s.Name)
		// execute the step
		err = c.ExecStep(ctx, s)
		if err != nil {
			return fmt.Errorf("unable to execute step: %w", err)
		}

		// check the step exit code
		if s.ExitCode != 0 {
			// check if we ignore step failures
			if !s.Ruleset.Continue {
				// set build status to failure
				b.SetStatus(constants.StatusFailure)
			}

			// update the step fields
			c.steps[s.ID].SetExitCode(s.ExitCode)
			c.steps[s.ID].SetStatus(constants.StatusFailure)
		}

		c.steps[s.ID].SetFinished(time.Now().UTC().Unix())
		c.logger.Infof("uploading %s step state", s.Name)
		// send API call to update the build
		_, _, err = c.Vela.Step.Update(r.GetOrg(), r.GetName(), b.GetNumber(), c.steps[s.ID])
		if err != nil {
			return err
		}
	}

	// create an error group with the context for each stage
	stages, stageCtx := errgroup.WithContext(ctx)
	// create a map to track the progress of each stage
	stageMap := make(map[string]chan error)

	// iterate through each stage in the pipeline
	for _, s := range p.Stages {
		// https://golang.org/doc/faq#closures_and_goroutines
		stage := s

		// create a new channel for each stage in the map
		stageMap[stage.Name] = make(chan error)

		stages.Go(func() error {
			c.logger.Infof("executing %s stage", stage.Name)
			// execute the stage
			err := c.ExecStage(stageCtx, stage, stageMap)
			if err != nil {
				return fmt.Errorf("unable to execute stage: %w", err)
			}

			return nil
		})
	}

	c.logger.Debug("waiting for stages completion")
	// wait for the stages to complete or return an error
	err = stages.Wait()
	if err != nil {
		return err
	}

	b.SetFinished(time.Now().UTC().Unix())
	c.logger.Info("uploading build state")
	// send API call to update the build
	_, _, err = c.Vela.Build.Update(r.GetOrg(), r.GetName(), b)
	if err != nil {
		return fmt.Errorf("unable to upload final state: %v", err)
	}

	return nil
}

// DestroyBuild cleans up the build after execution.
func (c *client) DestroyBuild(ctx context.Context) error {
	var err error
	b := c.build
	p := c.pipeline
	r := c.repo

	defer func() {
		// NOTE: When an error occurs during a build that does not have to do
		// with a pipeline we should set build status to "error" not "failed"
		// because it is worker related and not build.
		if err != nil {
			// update the build fields
			b.SetFinished(time.Now().UTC().Unix())
			b.SetStatus(constants.StatusError)
			b.SetError(err.Error())

			c.logger.Info("uploading errored stated")
			// send API call to update the build
			_, _, err = c.Vela.Build.Update(r.GetOrg(), r.GetName(), b)
			if err != nil {
				c.logger.Errorf("unable to upload errored state: %w", err)
			}
		}
	}()

	// destroy the steps for the pipeline
	for _, s := range p.Steps {
		c.logger.Infof("destroying %s step", s.Name)
		// destroy the step
		err = c.DestroyStep(ctx, s)
		if err != nil {
			c.logger.Errorf("unable to destroy step: %w", err)
		}
	}

	// destroy the stages for the pipeline
	for _, s := range p.Stages {
		c.logger.Infof("destroying %s stage", s.Name)
		// destroy the stage
		err = c.DestroyStage(ctx, s)
		if err != nil {
			c.logger.Errorf("unable to destroy stage: %w", err)
		}
	}

	// destroy the services for the pipeline
	for _, s := range p.Services {
		c.logger.Infof("destroying %s service", s.Name)
		// destroy the service
		err = c.DestroyService(ctx, s)
		if err != nil {
			c.logger.Errorf("unable to destroy service: %w", err)
		}

		c.logger.Infof("uploading %s service state", s.Name)
		// send API call to update the build
		c.services[s.ID].SetExitCode(s.ExitCode)
		c.services[s.ID].SetFinished(time.Now().UTC().Unix())
		_, _, err = c.Vela.Svc.Update(r.GetOrg(), r.GetName(), b.GetNumber(), c.services[s.ID])
		if err != nil {
			c.logger.Errorf("unable to upload service status: %w", err)
		}
	}

	c.logger.Info("deleting volume")
	// remove the runtime volume for the pipeline
	err = c.Runtime.RemoveVolume(ctx, p)
	if err != nil {
		c.logger.Errorf("unable to remove volume: %w", err)
	}

	c.logger.Info("deleting network")
	// remove the runtime network for the pipeline
	err = c.Runtime.RemoveNetwork(ctx, p)
	if err != nil {
		c.logger.Errorf("unable to remove network: %w", err)
	}

	return err
}
