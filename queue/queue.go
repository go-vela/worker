// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"github.com/go-vela/types"
)

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {
	// Pull defines a function that pops an item off the queue.
	Pull(string) (*types.Item, error)
}
