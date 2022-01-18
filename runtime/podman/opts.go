// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the runtime client for Podman.
type ClientOpt func(*client) error

// WithHostVolumes sets the host volumes in the runtime client for Podman.
func WithHostVolumes(volumes []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring host volumes in podman runtime client")

		// set the runtime host volumes in the podman client
		c.config.Volumes = volumes

		return nil
	}
}

// WithLogger sets the logger in the runtime client for Podman.
func WithLogger(logger *logrus.Entry) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring logger in podman client")

		// check if the logger provided is empty
		if logger != nil {
			// set the runtime logger in the podman client
			c.Logger = logger
		}

		return nil
	}
}

// WithPrivilegedImages sets the privileged images in the runtime client for Podman.
func WithPrivilegedImages(images []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring privileged images in podman runtime client")

		// set the runtime privileged images in the podman client
		c.config.Images = images

		return nil
	}
}
