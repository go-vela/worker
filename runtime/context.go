// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"context"
)

// key defines the key type for storing
// the runtime Engine in the context.
const key = "runtime"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the runtime Engine
// associated with this context.
func FromContext(c context.Context) Engine {
	// get runtime value from context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast runtime value to expected Engine type
	r, ok := v.(Engine)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the runtime Engine to this
// context if it supports the Setter interface.
func ToContext(c Setter, r Engine) {
	c.Set(key, r)
}
