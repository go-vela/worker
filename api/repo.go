// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

// swagger:operation GET /api/v1/executors/{executor}/repo repo GetRepo
//
// Get a currently running repo
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
//   description: The executor running the build
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '500':
//     description: Unable to retrieve the repo
//     type: json

// GetRepo represents the API handler to capture the
// repo currently running on an executor.
func GetRepo(c *gin.Context) {
	e := executor.Retrieve(c)

	repo, err := e.GetRepo()
	if err != nil {
		msg := fmt.Errorf("unable to read repo: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, repo)
}
