// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key defines the key type for storing
// the worker Worker in the context.
const key = "worker"

// FromContext retrieves the worker Worker from the context.Context.
func FromContext(c context.Context) Worker {
	// get worker value from context.Context
	v := c.Value(key)
	if v == nil {
		return Worker{}
	}

	// cast executor value to expected Worker type
	w, ok := v.(Worker)
	if !ok {
		return Worker{}
	}

	return w
}

// FromGinContext retrieves the executor Engine from the gin.Context.
func FromGinContext(c *gin.Context) Worker {
	// get executor value from gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Get
	v, ok := c.Get(key)
	if !ok {
		return Worker{}
	}

	// cast executor value to expected Engine type
	e, ok := v.(Worker)
	if !ok {
		return Worker{}
	}

	return e
}

// WithContext inserts the executor Engine into the context.Context.
func WithContext(c context.Context, w Worker) context.Context {
	// set the executor Engine in the context.Context
	//
	//nolint:revive,staticcheck // ignore using string with context value
	return context.WithValue(c, key, w)
}

// WithGinContext inserts the executor Engine into the gin.Context.
func WithGinContext(c *gin.Context, w Worker) {
	// set the executor Engine in the gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Set
	c.Set(key, w)
}
