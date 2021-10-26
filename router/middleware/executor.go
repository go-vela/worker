// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/executor"
)

// Executors is a middleware function that attaches the
// executors to the context of every http.Request.
func Executors(e map[int]executor.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("executors", e)
		c.Next()
	}
}
