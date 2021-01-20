// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/go-vela/types"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/router/middleware/executor"
)

// swagger:operation GET /api/v1/executors/{executor}/pipeline pipeline GetPipeline
//
// Get a currently running pipeline
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
//   description: The executor running the pipeline
//   required: true
//   type: string
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
