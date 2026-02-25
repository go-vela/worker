// SPDX-License-Identifier: Apache-2.0

package local

import (
	"errors"
	"os"
	"reflect"
	"sync"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
)

type (
	// client manages communication with the pipeline resources.
	client struct {
		Vela      *vela.Client
		Runtime   runtime.Engine
		Hostname  string
		Version   string
		OutputCtn *pipeline.Container

		// private fields
		init           *pipeline.Container
		build          *api.Build
		pipeline       *pipeline.Build
		services       sync.Map
		steps          sync.Map
		err            error
		streamRequests chan message.StreamRequest

		outputs *outputSvc

		// internal field partially exported for tests
		stdout           *os.File
		mockStdoutReader *os.File
	}

	svc struct {
		client *client
	}

	// MockedClient is for internal use to facilitate testing the local executor.
	MockedClient interface {
		MockStdout() *os.File
	}
)

// MockStdout is for internal use to facilitate testing the local executor.
// MockStdout returns a reader over a mocked Stdout.
func (c *client) MockStdout() *os.File {
	return c.mockStdoutReader
}

// equal returns true if the other client is the equivalent.
func Equal(a, b *client) bool {
	// handle any nil comparisons
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	return reflect.DeepEqual(a.Vela, b.Vela) &&
		reflect.DeepEqual(a.Runtime, b.Runtime) &&
		a.Hostname == b.Hostname &&
		a.Version == b.Version &&
		reflect.DeepEqual(a.init, b.init) &&
		reflect.DeepEqual(a.build, b.build) &&
		reflect.DeepEqual(a.pipeline, b.pipeline) &&
		reflect.DeepEqual(&a.services, &b.services) &&
		reflect.DeepEqual(&a.steps, &b.steps) &&
		errors.Is(a.err, b.err)
}

// New returns an Executor implementation that integrates with the local system.
//

func New(opts ...Opt) (*client, error) {
	// create new local client
	c := new(client)

	// Add stdout by default
	c.stdout = os.Stdout

	c.outputs = &outputSvc{client: c}

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
