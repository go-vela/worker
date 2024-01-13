// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// RegisterToken is a middleware function that attaches the
// auth-token channel to the context of every http.Request.
func RegisterToken(r chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("register-token", r)
		c.Next()
	}
}
