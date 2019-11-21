// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Shutdown represents the API handler to shutdown a
// worker currently running on an executor.
//
// This function performs a soft shut down of a worker.
// Any build running during this time will safely complete, then
// the worker will safely shut itself down.
func Shutdown(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "This endpoint is not yet implemented")
}
