// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/executor"
)

// Executors is a middleware function that attaches the
// executors to the context of every http.Request.
func Executors(e map[int]executor.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("executors", e)
		c.Next()
	}
}
