// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation POST /api/v1/shutdown router Shutdown
//
// Perform a soft shutdown of the worker
//
// ---
// x-success_http_code: '501'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela server token
//   required: true
//   type: string
// responses:
//   '501':
//     description: Endpoint is not yet implemented
//     schema:
//       type: string

// Shutdown represents the API handler to shutdown a
// executors currently running on an worker.
//
// This function performs a soft shut down of a worker.
// Any build running during this time will safely complete, then
// the worker will safely shut itself down.
func Shutdown(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "This endpoint is not yet implemented")
}
