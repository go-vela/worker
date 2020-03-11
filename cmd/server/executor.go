// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/constants"

	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/executor/linux"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// helper function to setup the queue from the CLI arguments.
func setupExecutor(c *cli.Context, client *vela.Client, r runtime.Engine) (executor.Engine, error) {
	logrus.Debug("Creating executor clients from CLI configuration")

	switch c.String("executor-driver") {
	case constants.DriverDarwin:
		return setupDarwin(client, r)
	case constants.DriverLinux:
		return setupLinux(client, r)
	case constants.DriverWindows:
		return setupWindows(client, r)
	default:
		return nil, fmt.Errorf("invalid executor driver: %s", c.String("executor-driver"))
	}
}

// helper function to setup the Darwin executor from the CLI arguments.
func setupDarwin(client *vela.Client, r runtime.Engine) (executor.Engine, error) {
	logrus.Tracef("Creating %s executor client from CLI configuration", constants.DriverDarwin)
	// return darwin.New(client, r)
	return nil, fmt.Errorf("unsupported executor driver: %s", constants.DriverDarwin)
}

// helper function to setup the Linux executor from the CLI arguments.
func setupLinux(client *vela.Client, r runtime.Engine) (executor.Engine, error) {
	logrus.Tracef("Creating %s executor client from CLI configuration", constants.DriverLinux)
	return linux.New(client, r)
}

// helper function to setup the Windows executor from the CLI arguments.
func setupWindows(client *vela.Client, r runtime.Engine) (executor.Engine, error) {
	logrus.Tracef("Creating %s executor client from CLI configuration", constants.DriverWindows)
	// return windows.New(client, r)
	return nil, fmt.Errorf("unsupported executor driver: %s", constants.DriverWindows)
}
