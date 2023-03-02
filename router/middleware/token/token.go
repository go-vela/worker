// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Retrieve gets the token from the provided request http.Request
// to be parsed and validated. This is called on every request
// to enable capturing the user making the request and validating
// they have the proper access. The following methods of providing
// authentication to Vela are supported:
//
// * Bearer token in `Authorization` header
// .
func Retrieve(c context.Context) *string {
	return FromContext(c)
}

// Establish sets the client in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.Request.Header.Get("Authorization")

		ToContext(c, &t)
		c.Next()
	}
}
