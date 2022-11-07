// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"reflect"
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
		init        *pipeline.Container
		logMethod   string
		maxLogSize  uint
		build       *library.Build
		pipeline    *pipeline.Build
		repo        *library.Repo
		secrets     sync.Map
		services    sync.Map
		serviceLogs sync.Map
		steps       sync.Map
		stepLogs    sync.Map

		streamRequests chan message.StreamRequest
		buildActivity  chan message.BuildActivity

		user *library.User
		err  error
	}

	svc struct {
		client *client
	}
)

// Equal returns true if the other client is the equivalent.
func Equal(a, b *client) bool {
	// handle any nil comparisons
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	return reflect.DeepEqual(a.Logger, b.Logger) &&
		reflect.DeepEqual(a.Vela, b.Vela) &&
		reflect.DeepEqual(a.Runtime, b.Runtime) &&
		reflect.DeepEqual(a.Secrets, b.Secrets) &&
		a.Hostname == b.Hostname &&
		a.Version == b.Version &&
		reflect.DeepEqual(a.init, b.init) &&
		a.logMethod == b.logMethod &&
		a.maxLogSize == b.maxLogSize &&
		reflect.DeepEqual(a.build, b.build) &&
		reflect.DeepEqual(a.pipeline, b.pipeline) &&
		reflect.DeepEqual(a.repo, b.repo) &&
		reflect.DeepEqual(&a.secrets, &b.secrets) &&
		reflect.DeepEqual(&a.services, &b.services) &&
		reflect.DeepEqual(&a.serviceLogs, &b.serviceLogs) &&
		reflect.DeepEqual(&a.steps, &b.steps) &&
		reflect.DeepEqual(&a.stepLogs, &b.stepLogs) &&
		// do not compare streamRequests channel
		reflect.DeepEqual(a.user, b.user) &&
		reflect.DeepEqual(a.err, b.err)
}

// New returns an Executor implementation that integrates with a Linux instance.
//
//nolint:revive // ignore unexported type as it is intentional
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

	// instantiate streamRequests channel (which may be overridden using withStreamRequests()).
	c.streamRequests = make(chan message.StreamRequest)

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

	return c, nil
}
