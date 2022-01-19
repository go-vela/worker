// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"net/http"
	"strings"

	"github.com/go-vela/worker/router/middleware/token"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the user in the given context.
func Retrieve(c *gin.Context) *library.User {
	return FromContext(c)
}

// Establish sets the user in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := new(library.User)

		t, err := token.Retrieve(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		secret := c.MustGet("secret").(string)
		if strings.EqualFold(t, secret) {
			u.SetName("vela-server")
			u.SetActive(true)
			u.SetAdmin(true)
		}

		ToContext(c, u)
		c.Next()
	}
}
