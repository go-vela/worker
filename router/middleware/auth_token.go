// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthToken is a middleware function that attaches the
// auth-token channel to the context of every http.Request.
func AuthToken(r chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("auth-token", r)
		c.Next()
	}
}
