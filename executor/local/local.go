// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"reflect"
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

// equal returns true if the other client is the equivalent.
func Equal(a, b *client) bool {
	if reflect.DeepEqual(a.Vela, b.Vela) &&
		reflect.DeepEqual(a.Runtime, b.Runtime) &&
		a.Hostname == b.Hostname &&
		a.Version == b.Version &&
		reflect.DeepEqual(a.init, b.init) &&
		reflect.DeepEqual(a.build, b.build) &&
		reflect.DeepEqual(a.pipeline, b.pipeline) &&
		reflect.DeepEqual(a.repo, b.repo) &&
		reflect.DeepEqual(&a.services, &b.services) &&
		reflect.DeepEqual(&a.steps, &b.steps) &&
		reflect.DeepEqual(a.user, b.user) &&
		reflect.DeepEqual(a.err, b.err) {
		return true
	}

	return false
}

// New returns an Executor implementation that integrates with the local system.
//
// nolint: revive // ignore unexported type as it is intentional
func New(opts ...Opt) (*client, error) {
	// create new local client
	c := new(client)

	// instantiate streamRequests channel (which may be overridden using withStreamRequests()).
	c.streamRequests = make(chan message.StreamRequest)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
