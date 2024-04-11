// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types"
	"github.com/go-vela/worker/router/middleware/executor"
)

// swagger:operation GET /api/v1/executors/{executor}/pipeline pipeline GetPipeline
//
// Get a currently running pipeline
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: executor
//   description: The executor running the pipeline
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
//   '500':
//     description: Unable to retrieve the pipeline
//     type: json

// GetPipeline represents the API handler to capture the
// pipeline currently running on an executor.
func GetPipeline(c *gin.Context) {
	e := executor.Retrieve(c)

	pipeline, err := e.GetPipeline()
	if err != nil {
		msg := fmt.Errorf("unable to read pipeline: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, pipeline)
}
