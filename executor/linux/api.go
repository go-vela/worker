// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"fmt"
	"time"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/service"
	"github.com/go-vela/worker/internal/step"
)

// GetBuild gets the current build in execution.
func (c *client) GetBuild() (*library.Build, error) {
	// check if the build resource is available
	if c.build == nil {
		return nil, fmt.Errorf("build resource not found")
	}

	return c.build, nil
}

// GetPipeline gets the current pipeline in execution.
func (c *client) GetPipeline() (*pipeline.Build, error) {
	// check if the pipeline resource is available
	if c.pipeline == nil {
		return nil, fmt.Errorf("pipeline resource not found")
	}

	return c.pipeline, nil
}

// GetRepo gets the current repo in execution.
func (c *client) GetRepo() (*api.Repo, error) {
	// check if the repo resource is available
	if c.repo == nil {
		return nil, fmt.Errorf("repo resource not found")
	}

	return c.repo, nil
}

// CancelBuild cancels the current build in execution.
//
//nolint:funlen // process of going through steps/services/stages is verbose and could be funcitonalized
func (c *client) CancelBuild() (*library.Build, error) {
	// get the current build from the client
	b, err := c.GetBuild()
	if err != nil {
		return nil, err
	}

	// set the build status to canceled
	b.SetStatus(constants.StatusCanceled)

	// get the current pipeline from the client
	pipeline, err := c.GetPipeline()
	if err != nil {
		return nil, err
	}

	// cancel non successful services
	//nolint:dupl // false positive, steps/services are different
	for _, _service := range pipeline.Services {
		// load the service from the client
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Load
		s, err := service.Load(_service, &c.services)
		if err != nil {
			// create the library service object
			s = new(library.Service)
			s.SetName(_service.Name)
			s.SetNumber(_service.Number)
			s.SetImage(_service.Image)
			s.SetStarted(time.Now().UTC().Unix())
			s.SetHost(c.build.GetHost())
			s.SetRuntime(c.build.GetRuntime())
			s.SetDistribution(c.build.GetDistribution())
		}

		// if service state was not terminal, set it as canceled
		switch s.GetStatus() {
		// service is in a error state
		case constants.StatusError:
			break
		// service is in a failure state
		case constants.StatusFailure:
			break
		// service is in a killed state
		case constants.StatusKilled:
			break
		// service is in a success state
		case constants.StatusSuccess:
			break
		default:
			// update the service with a canceled state
			s.SetStatus(constants.StatusCanceled)
			// add a service to a map
			c.services.Store(_service.ID, s)
		}
	}

	// cancel non successful steps
	//nolint:dupl // false positive, steps/services are different
	for _, _step := range pipeline.Steps {
		// load the step from the client
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
		s, err := step.Load(_step, &c.steps)
		if err != nil {
			// create the library step object
			s = new(library.Step)
			s.SetName(_step.Name)
			s.SetNumber(_step.Number)
			s.SetImage(_step.Image)
			s.SetStarted(time.Now().UTC().Unix())
			s.SetHost(c.build.GetHost())
			s.SetRuntime(c.build.GetRuntime())
			s.SetDistribution(c.build.GetDistribution())
		}

		// if step state was not terminal, set it as canceled
		switch s.GetStatus() {
		// step is in a error state
		case constants.StatusError:
			break
		// step is in a failure state
		case constants.StatusFailure:
			break
		// step is in a killed state
		case constants.StatusKilled:
			break
		// step is in a success state
		case constants.StatusSuccess:
			break
		default:
			// update the step with a canceled state
			s.SetStatus(constants.StatusCanceled)
			// add a step to a map
			c.steps.Store(_step.ID, s)
		}
	}

	// cancel non successful stages
	for _, _stage := range pipeline.Stages {
		// cancel non successful steps for that stage
		for _, _step := range _stage.Steps {
			// load the step from the client
			//
			// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
			s, err := step.Load(_step, &c.steps)
			if err != nil {
				// create the library step object
				s = new(library.Step)
				s.SetName(_step.Name)
				s.SetNumber(_step.Number)
				s.SetImage(_step.Image)
				s.SetStage(_stage.Name)
				s.SetStarted(time.Now().UTC().Unix())
				s.SetHost(c.build.GetHost())
				s.SetRuntime(c.build.GetRuntime())
				s.SetDistribution(c.build.GetDistribution())
			}

			// if stage state was not terminal, set it as canceled
			switch s.GetStatus() {
			// stage is in a error state
			case constants.StatusError:
				break
			// stage is in a failure state
			case constants.StatusFailure:
				break
			// stage is in a killed state
			case constants.StatusKilled:
				break
			// stage is in a success state
			case constants.StatusSuccess:
				break
			default:
				// update the step with a canceled state
				s.SetStatus(constants.StatusCanceled)
				// add a step to a map
				c.steps.Store(_step.ID, s)
			}
		}
	}

	err = c.DestroyBuild(context.Background())
	if err != nil {
		c.Logger.Errorf("unable to destroy build: %v", err)
	}

	return b, nil
}
