// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/worker/runtime"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Opt represents a configuration option to initialize the executor client for Linux.
type Opt func(*client) error

// WithBuild sets the library build in the executor client for Linux.
func WithBuild(b *library.Build) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring build in linux executor client")

		// check if the build provided is empty
		if b == nil {
			return fmt.Errorf("empty build provided")
		}

		// set the build in the client
		c.build = b

		return nil
	}
}

// WithLogMethod sets the method used to publish logs in the executor client for Linux.
func WithLogMethod(method string) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring log streaming method in linux executor client")

		// check if a method is provided
		if len(method) == 0 {
			return fmt.Errorf("empty log method provided")
		}

		// set the log method in the client
		c.logMethod = method

		return nil
	}
}

// WithMaxLogSize sets the maximum log size (in bytes) in the executor client for Linux.
func WithMaxLogSize(size uint) Opt {
	c.Logger.Trace("configuring maximum log size in linux executor client")

	return func(c *client) error {
		// set the maximum log size in the client
		c.maxLogSize = size

		return nil
	}
}

// WithHostname sets the hostname in the executor client for Linux.
func WithHostname(hostname string) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring hostname in linux executor client")

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

// WithLogger sets the logger in the executor client for Linux.
func WithLogger(logger *logrus.Entry) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring logger in linux executor client")

		// check if the logger provided is empty
		if logger != nil {
			// set the executor logger in the linux client
			c.Logger = logger
		}

		return nil
	}
}

// WithPipeline sets the pipeline build in the executor client for Linux.
func WithPipeline(p *pipeline.Build) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring pipeline in linux executor client")

		// check if the pipeline provided is empty
		if p == nil {
			return fmt.Errorf("empty pipeline provided")
		}

		// set the pipeline in the client
		c.pipeline = p

		return nil
	}
}

// WithRepo sets the library repo in the executor client for Linux.
func WithRepo(r *library.Repo) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring repository in linux executor client")

		// check if the repo provided is empty
		if r == nil {
			return fmt.Errorf("empty repo provided")
		}

		// set the repo in the client
		c.repo = r

		return nil
	}
}

// WithRuntime sets the runtime engine in the executor client for Linux.
func WithRuntime(r runtime.Engine) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring runtime in linux executor client")

		// check if the runtime provided is empty
		if r == nil {
			return fmt.Errorf("empty runtime provided")
		}

		// set the runtime in the client
		c.Runtime = r

		return nil
	}
}

// WithUser sets the library user in the executor client for Linux.
func WithUser(u *library.User) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring user in linux executor client")

		// check if the user provided is empty
		if u == nil {
			return fmt.Errorf("empty user provided")
		}

		// set the user in the client
		c.user = u

		return nil
	}
}

// WithVelaClient sets the Vela client in the executor client for Linux.
func WithVelaClient(cli *vela.Client) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring Vela client in linux executor client")

		// check if the Vela client provided is empty
		if cli == nil {
			return fmt.Errorf("empty Vela client provided")
		}

		// set the Vela client in the client
		c.Vela = cli

		return nil
	}
}

// WithVersion sets the version in the executor client for Linux.
func WithVersion(version string) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring version in linux executor client")

		// check if a version is provided
		if len(version) == 0 {
			// default the version to localhost
			version = "v0.0.0"
		}

		// set the version in the client
		c.Version = version

		return nil
	}
}
