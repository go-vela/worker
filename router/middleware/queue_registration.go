// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/types/library"
)

// WorkerRegistration is a middleware function that attaches the
// queue-address channel to the context of every http.Request.
func WorkerRegistration(r chan library.WorkerRegistration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("worker-registration", r)
		c.Next()
	}
}
