// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the runtime client.
type ClientOpt func(*client) error

// WithConfigFile sets the Kubernetes config file in the runtime client.
func WithConfigFile(file string) ClientOpt {
	logrus.Trace("configuring config file in kubernetes runtime client")

	return func(c *client) error {
		// set the runtime config file in the kubernetes client
		c.config.File = file

		return nil
	}
}

// WithNamespace sets the Kubernetes namespace in the runtime client.
func WithNamespace(namespace string) ClientOpt {
	logrus.Trace("configuring namespace in kubernetes runtime client")

	return func(c *client) error {
		// check if the namespace provided is empty
		if len(namespace) == 0 {
			return fmt.Errorf("no Kubernetes namespace provided")
		}

		// set the runtime namespace in the kubernetes client
		c.config.Namespace = namespace

		return nil
	}
}

// WithPrivilegedImages sets the Kubernetes privileged images in the runtime client.
func WithPrivilegedImages(images []string) ClientOpt {
	logrus.Trace("configuring privileged images in kubernetes runtime client")

	return func(c *client) error {
		// set the runtime privileged images in the kubernetes client
		c.config.Images = images

		return nil
	}
}

// WithHostVolumes sets the Kubernetes host volumes in the runtime client.
func WithHostVolumes(volumes []string) ClientOpt {
	logrus.Trace("configuring host volumes in kubernetes runtime client")

	return func(c *client) error {
		// set the runtime host volumes in the kubernetes client
		c.config.Volumes = volumes

		return nil
	}
}
