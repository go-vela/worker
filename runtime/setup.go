// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"fmt"

	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela engine capable of integrating
// with a configured runtime environment.
type Setup struct {
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry

	// Runtime Configuration

	// specifies the driver to use for the runtime client
	Driver string
	// specifies the path to a configuration file to use for the runtime client
	ConfigFile string
	// specifies a list of host volumes to use for the runtime client
	HostVolumes []string
	// specifies the namespace to use for the runtime client (only used by kubernetes)
	Namespace string
	// specifies a list of privileged images to use for the runtime client
	PrivilegedImages []string
}

// Docker creates and returns a Vela engine capable of
// integrating with a Docker runtime environment.
func (s *Setup) Docker() (Engine, error) {
	logrus.Trace("creating docker runtime client from setup")

	// create new Docker runtime engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/docker?tab=doc#New
	return docker.New(
		docker.WithHostVolumes(s.HostVolumes),
		docker.WithPrivilegedImages(s.PrivilegedImages),
		docker.WithLogger(s.Logger),
	)
}

// Kubernetes creates and returns a Vela engine capable of
// integrating with a Kubernetes runtime environment.
func (s *Setup) Kubernetes() (Engine, error) {
	logrus.Trace("creating kubernetes runtime client from setup")

	// create new Kubernetes runtime engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/kubernetes?tab=doc#New
	return kubernetes.New(
		kubernetes.WithConfigFile(s.ConfigFile),
		kubernetes.WithHostVolumes(s.HostVolumes),
		kubernetes.WithNamespace(s.Namespace),
		kubernetes.WithPrivilegedImages(s.PrivilegedImages),
		kubernetes.WithLogger(s.Logger),
	)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating runtime setup for client")

	// check if a runtime driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no runtime driver provided")
	}

	// process the secret driver being provided
	switch s.Driver {
	case constants.DriverDocker:
		break
	case constants.DriverKubernetes:
		// check if a runtime namespace was provided
		if len(s.Namespace) == 0 {
			return fmt.Errorf("no runtime namespace provided")
		}
	}

	// setup is valid
	return nil
}
