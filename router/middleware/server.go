// SPDX-License-Identifier: Apache-2.0

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
