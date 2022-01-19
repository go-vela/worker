// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
// x-success_http_code: '200'
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
