// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// Server is a middleware function that attaches the vela server address used for
// server <-> agent communication to the context of every http.Request.
func Server(server string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("server-address", server)
		c.Next()
	}
}
