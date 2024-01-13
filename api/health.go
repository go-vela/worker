// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /health system Health
//
// Check if the worker API is available
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successful 'ping' of Vela worker API
//     schema:
//       type: string

// Health check the status of the application.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
