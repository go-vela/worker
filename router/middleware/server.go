// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// ServerAddress is a middleware function that attaches the
// server address to the context of every http.Request.
func ServerAddress(addr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("server-address", addr)
		c.Next()
	}
}
