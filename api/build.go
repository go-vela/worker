// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/worker/router/middleware/executor"
)

// swagger:operation GET /api/v1/executors/{executor}/build build GetBuild
//
// Get the currently running build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: executor
//   description: The executor running the build
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '500':
//     description: Unable to retrieve the build
//     schema:
//       type: string

// GetBuild represents the API handler to capture the
// build currently running on an executor.
func GetBuild(c *gin.Context) {
	e := executor.Retrieve(c)

	build, err := e.GetBuild()
	if err != nil {
		msg := fmt.Errorf("unable to read build: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, build)
}

// swagger:operation DELETE /api/v1/executors/{executor}/build/cancel build CancelBuild
//
// Cancel the currently running build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: executor
//   description: The executor running the build
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully canceled the build
//     type: json
//   '500':
//     description: Unable to cancel the build
//     type: json

// CancelBuild represents the API handler to cancel a
// build currently running on an executor.
//
// This function performs a hard cancellation of a build on worker.
// Any build running during this time will immediately be stopped.
func CancelBuild(c *gin.Context) {
	e := executor.Retrieve(c)

	build, err := e.CancelBuild()
	if err != nil {
		msg := fmt.Errorf("unable to cancel build: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, build)
}
