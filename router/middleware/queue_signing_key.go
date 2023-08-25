// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// QueueSigningKey is a middleware function that attaches the
// auth-token channel to the context of every http.Request.
func QueueSigningKey(r chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue-signing-key", r)
		c.Next()
	}
}
