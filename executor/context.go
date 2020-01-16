// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"context"
)

// key defines the key type for storing
// the executor Engine in the context.
const key = "executor"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the executor Engine
// associated with this context.
func FromContext(c context.Context) map[int]Engine {
	// get executor value from context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast executor value to expected Engine type
	e, ok := v.(map[int]Engine)
	if !ok {
		return nil
	}

	return e
}

// ToContext adds the executor Engine to this
// context if it supports the Setter interface.
func ToContext(c Setter, e map[int]Engine) {
	c.Set(key, e)
}
