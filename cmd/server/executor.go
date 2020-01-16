// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"

	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/executor/linux"
	"github.com/go-vela/worker/runtime"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// helper function to setup the queue from the CLI arguments.
func setupExecutor(c *cli.Context, client *vela.Client, runtime runtime.Engine) (executor.Engine, error) {
	logrus.Debug("Creating executor clients from CLI configuration")

	switch c.String("executor-driver") {
	case constants.DriverDarwin:
		return setupDarwin(c, client, runtime)
	case constants.DriverLinux:
		return setupLinux(c, client, runtime)
	case constants.DriverWindows:
		return setupWindows(c, client, runtime)
	default:
		return nil, fmt.Errorf("invalid executor driver: %s", c.String("executor-driver"))
	}
}

// helper function to setup the Darwin executor from the CLI arguments.
func setupDarwin(c *cli.Context, client *vela.Client, runtime runtime.Engine) (executor.Engine, error) {
	logrus.Tracef("Creating %s executor client from CLI configuration", constants.DriverDarwin)
	// return darwin.New(client, runtime)
	return nil, fmt.Errorf("unsupported executor driver: %s", constants.DriverDarwin)
}

// helper function to setup the Linux executor from the CLI arguments.
func setupLinux(c *cli.Context, client *vela.Client, runtime runtime.Engine) (executor.Engine, error) {
	logrus.Tracef("Creating %s executor client from CLI configuration", constants.DriverLinux)
	return linux.New(client, runtime)
}

// helper function to setup the Windows executor from the CLI arguments.
func setupWindows(c *cli.Context, client *vela.Client, runtime runtime.Engine) (executor.Engine, error) {
	logrus.Tracef("Creating %s executor client from CLI configuration", constants.DriverWindows)
	// return windows.New(client, runtime)
	return nil, fmt.Errorf("unsupported executor driver: %s", constants.DriverWindows)
}
