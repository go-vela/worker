// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"time"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
)

// Flags represents all supported command line
// interface (CLI) flags for the executor.
//
// https://pkg.go.dev/github.com/urfave/cli#Flag
var Flags = []cli.Flag{
	// Executor Flags

	&cli.StringFlag{
		Name:  "executor.driver",
		Usage: "driver to be used for the executor",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_EXECUTOR_DRIVER"),
			cli.EnvVar("EXECUTOR_DRIVER"),
			cli.File("/vela/executor/driver"),
		),
		Value: constants.DriverLinux,
	},
	&cli.UintFlag{
		Name:  "executor.max_log_size",
		Usage: "maximum log size (in bytes)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_EXECUTOR_MAX_LOG_SIZE"),
			cli.EnvVar("EXECUTOR_MAX_LOG_SIZE"),
			cli.File("/vela/executor/max_log_size"),
		),
	},
	&cli.DurationFlag{
		Name:    "executor.log_streaming_timeout",
		Usage:   "maximum amount of time to wait for log streaming after build completes",
		Sources: cli.EnvVars("WORKER_LOG_STREAMING_TIMEOUT", "VELA_LOG_STREAMING_TIMEOUT", "LOG_STREAMING_TIMEOUT"),
		Value:   1 * time.Minute,
	},
	&cli.BoolFlag{
		Name:  "executor.enforce-trusted-repos",
		Usage: "enforce trusted repo restrictions for privileged images",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_EXECUTOR_ENFORCE_TRUSTED_REPOS"),
			cli.EnvVar("EXECUTOR_ENFORCE_TRUSTED_REPOS"),
			cli.File("/vela/executor/enforce_trusted_repos"),
		),
		Value: true,
	},
	&cli.StringFlag{
		Name:  "executor.outputs-image",
		Usage: "image used for the outputs container sidecar",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_EXECUTOR_OUTPUTS_IMAGE"),
			cli.EnvVar("EXECUTOR_OUTPUTS_IMAGE"),
			cli.File("/vela/executor/outputs_image"),
		),
	},
}
