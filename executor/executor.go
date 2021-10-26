// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// nolint: godot // ignore period at end for comment ending in a list
//
// New creates and returns a Vela engine capable of
// integrating with the configured executor.
//
// Currently the following executors are supported:
//
// * linux
// * local
func New(s *Setup) (Engine, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/worker/executor?tab=doc#Setup.Validate
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
		// https://pkg.go.dev/github.com/go-vela/worker/executor?tab=doc#Setup.Darwin
		return s.Darwin()
	case constants.DriverLinux:
		// handle the Linux executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor?tab=doc#Setup.Linux
		return s.Linux()
	case constants.DriverLocal:
		// handle the Local executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor?tab=doc#Setup.Local
		return s.Local()
	case constants.DriverWindows:
		// handle the Windows executor driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/worker/executor?tab=doc#Setup.Windows
		return s.Windows()
	default:
		// handle an invalid executor driver being provided
		return nil, fmt.Errorf("invalid executor driver provided: %s", s.Driver)
	}
}
