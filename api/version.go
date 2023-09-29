// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/version"
)

// swagger:operation GET /version router Version
//
// Get the version of the Vela API
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully retrieved the Vela API version
//     schema:
//       type: string

// Version represents the API handler to
// report the version information for Vela.
func Version(c *gin.Context) {
	c.JSON(http.StatusOK, version.New())
}
