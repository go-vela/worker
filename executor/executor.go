// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// New creates and returns a Vela engine capable of
// integrating with the configured executor.
//
// Currently the following executors are supported:
//
// * linux
// * local
// .
func New(s *Setup) (Engine, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/worker/executor#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating executor engine from setup")
	// process the executor driver being provided
	switch s.Driver {
	case constants.DriverDarwin:
		// handle the Darwin executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor#Setup.Darwin
		return s.Darwin()
	case constants.DriverLinux:
		// handle the Linux executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor#Setup.Linux
		return s.Linux()
	case constants.DriverLocal:
		// handle the Local executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor#Setup.Local
		return s.Local()
	case constants.DriverWindows:
		// handle the Windows executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor#Setup.Windows
		return s.Windows()
	default:
		// handle an invalid executor driver being provided
		return nil, fmt.Errorf("invalid executor driver provided: %s", s.Driver)
	}
}
