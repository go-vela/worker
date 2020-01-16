// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health check the status of the application
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
