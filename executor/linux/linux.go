// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"sync"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/sirupsen/logrus"
)

type (
	// client manages communication with the pipeline resources.
	client struct {
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		Logger   *logrus.Entry
		Vela     *vela.Client
		Runtime  runtime.Engine
		Secrets  map[string]*library.Secret
		Hostname string
		Version  string

		// clients for build actions
		secret *secretSvc

		// private fields
		init       *pipeline.Container
		logMethod  string
		maxLogSize uint
		build      *library.Build
		pipeline   *pipeline.Build
		repo       *library.Repo
		// nolint: structcheck,unused // ignore false positives
		secrets     sync.Map
		services    sync.Map
		serviceLogs sync.Map
		steps       sync.Map
		stepLogs    sync.Map

		streamRequests chan message.StreamRequest

		user *library.User
		err  error
	}

	// nolint: structcheck // ignore false positive
	svc struct {
		client *client
	}
)

// New returns an Executor implementation that integrates with a Linux instance.
//
// nolint: revive // ignore unexported type as it is intentional
func New(opts ...Opt) (*client, error) {
	// create new Linux client
	c := new(client)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// instantiate map for non-plugin secrets
	c.Secrets = make(map[string]*library.Secret)

	// instantiate all client services
	c.secret = &secretSvc{client: c}

	// instantiate streamRequests channel
	c.streamRequests = make(chan message.StreamRequest)

	return c, nil
}
