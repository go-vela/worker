// SPDX-License-Identifier: Apache-2.0

package local

import (
	"fmt"
	"os"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
)

// Opt represents a configuration option to initialize the executor client for Local.
type Opt func(*client) error

// WithBuild sets the library build in the executor client for Local.
func WithBuild(b *library.Build) Opt {
	return func(c *client) error {
		// set the build in the client
		c.build = b

		return nil
	}
}

// WithHostname sets the hostname in the executor client for Local.
func WithHostname(hostname string) Opt {
	return func(c *client) error {
		// check if a hostname is provided
		if len(hostname) == 0 {
			// default the hostname to localhost
			hostname = "localhost"
		}

		// set the hostname in the client
		c.Hostname = hostname

		return nil
	}
}

// WithPipeline sets the pipeline build in the executor client for Local.
func WithPipeline(p *pipeline.Build) Opt {
	return func(c *client) error {
		// check if the pipeline provided is empty
		if p == nil {
			return fmt.Errorf("empty pipeline provided")
		}

		// set the pipeline in the client
		c.pipeline = p

		return nil
	}
}

// WithRepo sets the library repo in the executor client for Local.
func WithRepo(r *api.Repo) Opt {
	return func(c *client) error {
		// set the repo in the client
		c.repo = r

		return nil
	}
}

// WithRuntime sets the runtime engine in the executor client for Local.
func WithRuntime(r runtime.Engine) Opt {
	return func(c *client) error {
		// check if the runtime provided is empty
		if r == nil {
			return fmt.Errorf("empty runtime provided")
		}

		// set the runtime in the client
		c.Runtime = r

		return nil
	}
}

// WithVelaClient sets the Vela client in the executor client for Local.
func WithVelaClient(cli *vela.Client) Opt {
	return func(c *client) error {
		// set the Vela client in the client
		c.Vela = cli

		return nil
	}
}

// WithVersion sets the version in the executor client for Local.
func WithVersion(version string) Opt {
	return func(c *client) error {
		// check if a version is provided
		if len(version) == 0 {
			// default the version
			version = "v0.0.0"
		}

		// set the version in the client
		c.Version = version

		return nil
	}
}

// WithMockStdout adds a mock stdout writer to the client if mock is true.
// If mock is true, then you must use a goroutine to read from
// MockStdout as quickly as possible, or writing to stdout will hang.
func WithMockStdout(mock bool) Opt {
	return func(c *client) error {
		if !mock {
			return nil
		}

		// New() sets c.stdout = os.stdout, replace it if a mock is required.
		reader, writer, err := os.Pipe()
		if err != nil {
			return err
		}

		c.mockStdoutReader = reader
		c.stdout = writer

		return nil
	}
}

// withStreamRequests sets the streamRequests channel in the executor client for Linux
// (primarily used for tests).
func withStreamRequests(s chan message.StreamRequest) Opt {
	return func(c *client) error {
		// set the streamRequests channel in the client
		c.streamRequests = s

		return nil
	}
}
