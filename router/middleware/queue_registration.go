// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/types/library"
)

// QueueRegistration RegisterToken is a middleware function that attaches the
// queue-registration channel to the context of every http.Request.
func QueueRegistration(r chan library.QueueRegistration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue-registration", r)
		c.Next()
	}
}