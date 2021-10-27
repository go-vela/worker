// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key defines the key type for storing
// the runtime Engine in the context.
const key = "runtime"

// FromContext retrieves the runtime Engine from the context.Context.
func FromContext(c context.Context) Engine {
	// get runtime value from context.Context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast runtime value to expected Engine type
	e, ok := v.(Engine)
	if !ok {
		return nil
	}

	return e
}

// FromGinContext retrieves the runtime Engine from the gin.Context.
func FromGinContext(c *gin.Context) Engine {
	// get runtime value from gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Get
	v, ok := c.Get(key)
	if !ok {
		return nil
	}

	// cast runtime value to expected Engine type
	e, ok := v.(Engine)
	if !ok {
		return nil
	}

	return e
}

// WithContext inserts the runtime Engine into the context.Context.
func WithContext(c context.Context, e Engine) context.Context {
	// set the runtime Engine in the context.Context
	//
	// nolint: golint,staticcheck // ignore using string with context value
	return context.WithValue(c, key, e)
}

// WithGinContext inserts the runtime Engine into the gin.Context.
func WithGinContext(c *gin.Context, e Engine) {
	// set the runtime Engine in the gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Set
	c.Set(key, e)
}
