// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/executor"
)

// Executor is a middleware function that initializes the executor and
// attaches to the context of every http.Request.
func Executor(e map[int]executor.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		executor.ToContext(c, e)
		c.Next()
	}
}
