// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// WorkerHostname is a middleware function that attaches the
// worker hostname to the context of every http.Request.
func WorkerHostname(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("worker-hostname", name)
		c.Next()
	}
}
