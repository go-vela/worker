// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"fmt"

	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Opt represents a configuration option to initialize the client.
type Opt func(*client) error

// WithBuild sets the library build in the client.
func WithBuild(b *library.Build) Opt {
	return func(c *client) error {
		// set the build in the client
		c.build = b

		return nil
	}
}

// WithHostname sets the hostname in the client.
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

// WithPipeline sets the pipeline build in the client.
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

// WithRepo sets the library repo in the client.
func WithRepo(r *library.Repo) Opt {
	return func(c *client) error {
		// set the repo in the client
		c.repo = r

		return nil
	}
}

// WithRuntime sets the runtime engine in the client.
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

// WithUser sets the library user in the client.
func WithUser(u *library.User) Opt {
	return func(c *client) error {
		// set the user in the client
		c.user = u

		return nil
	}
}

// WithVelaClient sets the Vela client in the client.
func WithVelaClient(cli *vela.Client) Opt {
	return func(c *client) error {
		// set the Vela client in the client
		c.Vela = cli

		return nil
	}
}

// WithVersion sets the version in the client.
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
