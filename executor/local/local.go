// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"sync"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
)

type (
	// client manages communication with the pipeline resources.
	client struct {
		Vela     *vela.Client
		Runtime  runtime.Engine
		Hostname string
		Version  string

		// private fields
		init           *pipeline.Container
		build          *library.Build
		pipeline       *pipeline.Build
		repo           *library.Repo
		services       sync.Map
		steps          sync.Map
		user           *library.User
		err            error
		streamRequests chan message.StreamRequest
	}
)

// New returns an Executor implementation that integrates with the local system.
//
// nolint: revive // ignore unexported type as it is intentional
func New(opts ...Opt) (*client, error) {
	// create new local client
	c := new(client)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// instantiate streamRequests channel
	c.streamRequests = make(chan message.StreamRequest)

	return c, nil
}
