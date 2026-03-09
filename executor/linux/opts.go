// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
)

// Opt represents a configuration option to initialize the executor client for Linux.
type Opt func(*client) error

// WithBuild sets the library build in the executor client for Linux.
func WithBuild(b *api.Build) Opt {
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

// WithMaxLogSize sets the maximum log size (in bytes) in the executor client for Linux.
func WithMaxLogSize(size uint) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring maximum log size in linux executor client")

		// set the maximum log size in the client
		c.maxLogSize = size

		return nil
	}
}

// WithFileSizeLimit sets the maximum file size (in MB) for a single file upload in the executor client for Linux.
func WithFileSizeLimit(limit int) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring file size limit in linux executor client")

		// set the file size limit in the client
		c.fileSizeLimit = int64(limit) * 1024 * 1024

		return nil
	}
}

// WithBuildFileSizeLimit sets the maximum total size (in MB) for all file uploads in a single build in the executor client for Linux.
func WithBuildFileSizeLimit(limit int) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring build file size limit in linux executor client")

		// set the build file size limit in the client
		c.buildFileSizeLimit = int64(limit) * 1024 * 1024

		return nil
	}
}

// WithLogStreamingTimeout sets the log streaming timeout in the executor client for Linux.
func WithLogStreamingTimeout(timeout time.Duration) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring log streaming timeout in linux executor client")

		// set the maximum log size in the client
		c.logStreamingTimeout = timeout

		return nil
	}
}

// WithPrivilegedImages sets the privileged images in the executor client for Linux.
func WithPrivilegedImages(images []string) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring privileged images in linux executor client")

		// set the privileged images in the client
		c.privilegedImages = images

		return nil
	}
}

// WithEnforceTrustedRepos configures trusted repo restrictions in the executor client for Linux.
func WithEnforceTrustedRepos(enforce bool) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring trusted repo restrictions in linux executor client")

		// set trusted repo restrictions in the client
		c.enforceTrustedRepos = enforce

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

// WithOutputCtn sets the outputs container in the executor client for Linux.
func WithOutputCtn(ctn *pipeline.Container) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring output container in linux executor client")

		// set the outputs container in the client
		c.OutputCtn = ctn

		return nil
	}
}

// withStreamRequests sets the streamRequests channel in the executor client for Linux
// (primarily used for tests).
func withStreamRequests(s chan message.StreamRequest) Opt {
	return func(c *client) error {
		c.Logger.Trace("configuring stream requests in linux executor client")

		// set the streamRequests channel in the client
		c.streamRequests = s

		return nil
	}
}
