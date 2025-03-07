// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/worker/executor"
	exec "github.com/go-vela/worker/router/middleware/executor"
)

// swagger:operation GET /api/v1/executors/{executor} executor GetExecutor
//
// Get a currently running executor
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: executor
//   description: The executor to retrieve
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
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
	executor := &api.Executor{}

	// TODO: Add this information from the context or helpers on executor
	// tmp.SetHost(executor.GetHost())
	executor.SetRuntime("docker")
	executor.SetDistribution("linux")

	// get build on executor
	executor.Build, err = e.GetBuild()
	if err != nil {
		msg := fmt.Errorf("unable to retrieve build: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	// get pipeline on executor
	executor.Pipeline, err = e.GetPipeline()
	if err != nil {
		msg := fmt.Errorf("unable to retrieve pipeline: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, executor)
}

// swagger:operation GET /api/v1/executors executor GetExecutors
//
// Get all currently running executors
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
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
		msg := "no running executors found"

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	// cast executors value to expected type
	e, ok := value.(map[int]executor.Engine)
	if !ok {
		msg := "unable to get executors"

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	executors := []*api.Executor{}

	for id, executor := range e {
		// create a temporary executor to append results to response
		tmp := &api.Executor{}

		// TODO: Add this information from the context or helpers on executor
		// tmp.SetHost(executor.GetHost())
		tmp.SetRuntime("docker")
		tmp.SetDistribution("linux")
		tmp.SetID(int64(id))

		// get build on executor
		tmp.Build, err = executor.GetBuild()
		if err != nil {
			msg := fmt.Errorf("unable to retrieve build: %w", err).Error()

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

			return
		}

		// get pipeline on executor
		tmp.Pipeline, err = executor.GetPipeline()
		if err != nil {
			msg := fmt.Errorf("unable to retrieve pipeline: %w", err).Error()

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

			return
		}

		executors = append(executors, tmp)
	}

	c.JSON(http.StatusOK, executors)
}
