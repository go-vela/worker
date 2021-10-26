// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"github.com/go-vela/types/constants"

	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the executor.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{

	// Logging Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_LOG_FORMAT", "EXECUTOR_LOG_FORMAT"},
		FilePath: "/vela/executor/log_format",
		Name:     "executor.log.format",
		Usage:    "format of logs to output",
		Value:    "json",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_LOG_LEVEL", "EXECUTOR_LOG_LEVEL"},
		FilePath: "/vela/executor/log_level",
		Name:     "executor.log.level",
		Usage:    "level of logs to output",
		Value:    "info",
	},

	// Executor Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_EXECUTOR_DRIVER", "EXECUTOR_DRIVER"},
		FilePath: "/vela/executor/driver",
		Name:     "executor.driver",
		Usage:    "driver to be used for the executor",
		Value:    constants.DriverLinux,
	},
}
