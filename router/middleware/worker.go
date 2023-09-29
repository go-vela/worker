// SPDX-License-Identifier: Apache-2.0

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
