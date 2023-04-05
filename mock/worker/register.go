// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// postRegister returns mock JSON for a http POST.
//
// Do not pass an auth token to fail the request.
func postRegister(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if len(token) == 0 {
		c.JSON(http.StatusInternalServerError, "no token provided in Authorization header")
	}

	c.JSON(http.StatusOK, "successfully passed token to worker")
}
