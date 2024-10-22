// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/worker/router/middleware/executor"
)

// swagger:operation GET /api/v1/executors/{executor}/repo repo GetRepo
//
// Get a currently running repo
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

	build, err := e.GetBuild()
	if err != nil {
		msg := fmt.Errorf("unable to read build: %w", err).Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, build.GetRepo())
}
