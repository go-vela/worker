// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"context"

	"github.com/go-vela/worker/executor"
)

const key = "executor"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Executor associated with this context.
func FromContext(c context.Context) executor.Engine {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	r, ok := value.(executor.Engine)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the Executor to this context if it supports
// the Setter interface.
func ToContext(c Setter, e executor.Engine) {
	c.Set(key, e)
}
