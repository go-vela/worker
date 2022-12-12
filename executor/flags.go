// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	// Executor Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_EXECUTOR_DRIVER", "EXECUTOR_DRIVER"},
		FilePath: "/vela/executor/driver",
		Name:     "executor.driver",
		Usage:    "driver to be used for the executor",
		Value:    constants.DriverLinux,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_EXECUTOR_LOG_METHOD", "EXECUTOR_LOG_METHOD"},
		FilePath: "/vela/executor/log_method",
		Name:     "executor.log_method",
		Usage:    "method used to publish logs to the server - options: (byte-chunks|time-chunks)",
		Value:    "byte-chunks",
	},
	&cli.UintFlag{
		EnvVars:  []string{"VELA_EXECUTOR_MAX_LOG_SIZE", "EXECUTOR_MAX_LOG_SIZE"},
		FilePath: "/vela/executor/max_log_size",
		Name:     "executor.max_log_size",
		Usage:    "maximum log size (in bytes)",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_EXECUTOR_ENFORCE_TRUSTED_REPOS", "EXECUTOR_ENFORCE_TRUSTED_REPOS"},
		FilePath: "/vela/executor/enforce_trusted_repos",
		Name:     "executor.enforce-trusted-repos",
		Usage:    "enforce trusted repo restrictions for privileged images",
		Value:    true,
	},
}
