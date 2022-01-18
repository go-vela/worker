// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// nolint: godot // ignore period at end for comment ending in a list
//
// New creates and returns a Vela engine capable of
// integrating with the configured runtime.
//
// Currently the following runtimes are supported:
//
// * docker
// * kubernetes
func New(s *Setup) (Engine, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating runtime engine from setup")
	// process the runtime driver being provided
	switch s.Driver {
	case constants.DriverDocker:
		// handle the Docker runtime driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime?tab=doc#Setup.Docker
		return s.Docker()
	case constants.DriverPodman:
		return s.Podman()
	case constants.DriverKubernetes:
		// handle the Kubernetes runtime driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime?tab=doc#Setup.Kubernetes
		return s.Kubernetes()
	default:
		// handle an invalid runtime driver being provided
		return nil, fmt.Errorf("invalid runtime driver provided: %s", s.Driver)
	}
}
