// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the runtime client.
type ClientOpt func(*client) error

// WithPrivilegedImages sets the Docker privileged images in the runtime client.
func WithPrivilegedImages(images []string) ClientOpt {
	logrus.Trace("configuring privileged images in docker runtime client")

	return func(c *client) error {
		// set the runtime privileged images in the docker client
		c.config.Images = images

		return nil
	}
}

// WithHostVolumes sets the Docker host volumes in the runtime client.
func WithHostVolumes(volumes []string) ClientOpt {
	logrus.Trace("configuring host volumes in docker runtime client")

	return func(c *client) error {
		// set the runtime host volumes in the docker client
		c.config.Volumes = volumes

		return nil
	}
}
