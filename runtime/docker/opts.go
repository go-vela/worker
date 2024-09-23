// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the runtime client for Docker.
type ClientOpt func(*client) error

// WithHostVolumes sets the host volumes in the runtime client for Docker.
func WithHostVolumes(volumes []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring host volumes in docker runtime client")

		// set the runtime host volumes in the docker client
		c.config.Volumes = volumes

		return nil
	}
}

// WithLogger sets the logger in the runtime client for Docker.
func WithLogger(logger *logrus.Entry) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring logger in docker runtime client")

		// check if the logger provided is empty
		if logger != nil {
			// set the runtime logger in the docker client
			c.Logger = logger
		}

		return nil
	}
}

// WithPrivilegedImages sets the privileged images in the runtime client for Docker.
func WithPrivilegedImages(images []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring privileged images in docker runtime client")

		// set the runtime privileged images in the docker client
		c.config.Images = images

		return nil
	}
}

// WithDropCapabilities sets the kernel capabilities to drop from each container in the runtime client for Docker.
func WithDropCapabilities(caps []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring dropped capabilities in docker runtime client")

		// set the runtime dropped kernel capabilities in the docker client
		c.config.DropCapabilities = caps

		return nil
	}
}

func WithContainerPlatform(platform string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring container platform in docker runtime client")
		// set the runtime container platform in the docker client
		c.config.ContainerPlatform = platform
		return nil
	}
}
