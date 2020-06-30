// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/pkg-executor/executor"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	exec "github.com/go-vela/worker/router/middleware/executor"
)

// swagger:operation GET /api/v1/executors/{executor} executor GetExecutor
//
// Get a currently running executor
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela server token
//   required: true
//   type: string
// - in: path
//   name: executor
//   description: The executor to retrieve
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the executor
//     type: json
//     schema:
//       "$ref": "#/definitions/Executor"
//   '500':
//     description: Unable to retrieve the executor
//     type: json

// GetExecutor represents the API handler to capture the
// executor currently running on a worker.
func GetExecutor(c *gin.Context) {
	var err error

	e := exec.Retrieve(c)
	executor := &library.Executor{}

	// TODO: Add this information from the context or helpers on executor
	// tmp.SetHost(executor.GetHost())
	executor.SetRuntime("docker")
	executor.SetDistribution("linux")

	// get build on executor
	executor.Build, err = e.GetBuild()
	if err != nil {
		msg := fmt.Errorf("unable to retrieve build: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	// get pipeline on executor
	executor.Pipeline, err = e.GetPipeline()
	if err != nil {
		msg := fmt.Errorf("unable to retrieve pipeline: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	// get repo on executor
	executor.Repo, err = e.GetRepo()
	if err != nil {
		msg := fmt.Errorf("unable to retrieve repo: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, executor)
}

// swagger:operation GET /api/v1/executors executor GetExecutors
//
// Get all currently running executors
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela server token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved all running executors
//     type: json
//     schema:
//       "$ref": "#/definitions/Executor"
//   '500':
//     description: Unable to retrieve all running executors
//     type: json

// GetExecutors represents the API handler to capture the
// executors currently running on a worker.
func GetExecutors(c *gin.Context) {
	var err error

	// capture executors value from context
	value := c.Value("executors")
	if value == nil {
		msg := fmt.Sprintf("no running executors found")

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	// cast executors value to expected type
	e, ok := value.(map[int]executor.Engine)
	if !ok {
		msg := fmt.Sprintf("unable to get executors")

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	executors := []*library.Executor{}

	for _, executor := range e {
		// create a temporary executor to append results to response
		tmp := &library.Executor{}

		// TODO: Add this information from the context or helpers on executor
		// tmp.SetHost(executor.GetHost())
		tmp.SetRuntime("docker")
		tmp.SetDistribution("linux")

		// get build on executor
		tmp.Build, err = executor.GetBuild()
		if err != nil {
			msg := fmt.Errorf("unable to retrieve build: %w", err).Error()

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		// get pipeline on executor
		tmp.Pipeline, err = executor.GetPipeline()
		if err != nil {
			msg := fmt.Errorf("unable to retrieve pipeline: %w", err).Error()

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		// get repo on executor
		tmp.Repo, err = executor.GetRepo()
		if err != nil {
			msg := fmt.Errorf("unable to retrieve repo: %w", err).Error()

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		executors = append(executors, tmp)
	}

	c.JSON(http.StatusOK, executors)
}
