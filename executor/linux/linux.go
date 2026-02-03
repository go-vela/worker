// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
)

type (
	// client manages communication with the pipeline resources.
	client struct {
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		Logger       *logrus.Entry
		Vela         *vela.Client
		Runtime      runtime.Engine
		Secrets      map[string]*api.Secret
		NoSubSecrets map[string]*api.Secret
		Hostname     string
		Version      string
		OutputCtn    *pipeline.Container

		// clients for build actions
		secret  *secretSvc
		outputs *outputSvc

		// private fields
		init                *pipeline.Container
		maxLogSize          uint
		logStreamingTimeout time.Duration
		privilegedImages    []string
		enforceTrustedRepos bool
		build               *api.Build
		pipeline            *pipeline.Build
		secrets             sync.Map
		services            sync.Map
		serviceLogs         sync.Map
		steps               sync.Map
		stepLogs            sync.Map

		streamRequests chan message.StreamRequest

		err error
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
		reflect.DeepEqual(a.NoSubSecrets, b.NoSubSecrets) &&
		a.Hostname == b.Hostname &&
		a.Version == b.Version &&
		reflect.DeepEqual(a.init, b.init) &&
		a.maxLogSize == b.maxLogSize &&
		reflect.DeepEqual(a.privilegedImages, b.privilegedImages) &&
		a.enforceTrustedRepos == b.enforceTrustedRepos &&
		reflect.DeepEqual(a.build, b.build) &&
		reflect.DeepEqual(a.pipeline, b.pipeline) &&
		reflect.DeepEqual(&a.secrets, &b.secrets) &&
		reflect.DeepEqual(&a.services, &b.services) &&
		reflect.DeepEqual(&a.serviceLogs, &b.serviceLogs) &&
		reflect.DeepEqual(&a.steps, &b.steps) &&
		reflect.DeepEqual(&a.stepLogs, &b.stepLogs) &&
		errors.Is(a.err, b.err)
}

// New returns an Executor implementation that integrates with a Linux instance.
//
//nolint:revive // ignore unexported type as it is intentional
func New(opts ...Opt) (*client, error) {
	// create new Linux client
	c := new(client)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
	c.Logger = logrus.NewEntry(logger)

	// instantiate streamRequests channel (which may be overridden using withStreamRequests()).
	// messages get sent during ExecBuild, then ExecBuild closes this on exit.
	c.streamRequests = make(chan message.StreamRequest)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// instantiate map for non-plugin secrets
	c.Secrets = make(map[string]*api.Secret)

	// instantiate map for non-substituted secrets
	c.NoSubSecrets = make(map[string]*api.Secret)

	// instantiate all client services
	c.secret = &secretSvc{client: c}
	c.outputs = &outputSvc{client: c}

	return c, nil
}
