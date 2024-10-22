// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/constants"
)

// Flags represents all supported command line
// interface (CLI) flags for the executor.
//
// https://pkg.go.dev/github.com/urfave/cli#Flag
var Flags = []cli.Flag{
	// Executor Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_EXECUTOR_DRIVER", "EXECUTOR_DRIVER"},
		FilePath: "/vela/executor/driver",
		Name:     "executor.driver",
		Usage:    "driver to be used for the executor",
		Value:    constants.DriverLinux,
	},
	&cli.UintFlag{
		EnvVars:  []string{"VELA_EXECUTOR_MAX_LOG_SIZE", "EXECUTOR_MAX_LOG_SIZE"},
		FilePath: "/vela/executor/max_log_size",
		Name:     "executor.max_log_size",
		Usage:    "maximum log size (in bytes)",
	},
	&cli.DurationFlag{
		EnvVars: []string{"WORKER_LOG_STREAMING_TIMEOUT", "VELA_LOG_STREAMING_TIMEOUT", "LOG_STREAMING_TIMEOUT"},
		Name:    "executor.log_streaming_timeout",
		Usage:   "maximum amount of time to wait for log streaming after build completes",
		Value:   5 * time.Minute,
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_EXECUTOR_ENFORCE_TRUSTED_REPOS", "EXECUTOR_ENFORCE_TRUSTED_REPOS"},
		FilePath: "/vela/executor/enforce_trusted_repos",
		Name:     "executor.enforce-trusted-repos",
		Usage:    "enforce trusted repo restrictions for privileged images",
		Value:    true,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_EXECUTOR_OUTPUTS_IMAGE", "EXECUTOR_OUTPUTS_IMAGE"},
		FilePath: "/vela/executor/outputs_image",
		Name:     "executor.outputs-image",
		Usage:    "image used for the outputs container sidecar",
	},
}
