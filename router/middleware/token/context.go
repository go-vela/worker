// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"context"
)

const key = "token"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the User associated with this context.
func FromContext(c context.Context) *string {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	u, ok := value.(*string)
	if !ok {
		return nil
	}

	return u
}

// ToContext adds the User to this context if it supports
// the Setter interface.
func ToContext(s Setter, vc *string) {
	s.Set(key, vc)
}
