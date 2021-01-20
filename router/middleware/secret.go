// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// Secret is a middleware function that attaches the secret used for
// server <-> agent communication to the context of every http.Request.
func Secret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("secret", secret)
		c.Next()
	}
}
