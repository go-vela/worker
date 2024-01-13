// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// New creates and returns a Vela engine capable of
// integrating with the configured runtime.
//
// Currently the following runtimes are supported:
//
// * docker
// * kubernetes
// .
func New(s *Setup) (Engine, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime#Setup.Validate
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
		// https://pkg.go.dev/github.com/go-vela/worker/runtime#Setup.Docker
		return s.Docker()
	case constants.DriverKubernetes:
		// handle the Kubernetes runtime driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime#Setup.Kubernetes
		return s.Kubernetes()
	default:
		// handle an invalid runtime driver being provided
		return nil, fmt.Errorf("invalid runtime driver provided: %s", s.Driver)
	}
}
