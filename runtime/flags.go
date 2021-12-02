// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"github.com/go-vela/types/constants"

	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the runtime.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	// Runtime Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_RUNTIME_DRIVER", "RUNTIME_DRIVER"},
		FilePath: "/vela/runtime/driver",
		Name:     "runtime.driver",
		Usage:    "driver to be used for the runtime",
		Value:    constants.DriverDocker,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_RUNTIME_CONFIG", "RUNTIME_CONFIG"},
		FilePath: "/vela/runtime/config",
		Name:     "runtime.config",
		Usage:    "path to configuration file for the runtime",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_RUNTIME_NAMESPACE", "RUNTIME_NAMESPACE"},
		FilePath: "/vela/runtime/namespace",
		Name:     "runtime.namespace",
		Usage:    "namespace to use for the runtime (only used by kubernetes)",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_RUNTIME_PRIVILEGED_IMAGES", "RUNTIME_PRIVILEGED_IMAGES"},
		FilePath: "/vela/runtime/privileged_images",
		Name:     "runtime.privileged-images",
		Usage:    "list of images allowed to run in privileged mode for the runtime",
		Value:    cli.NewStringSlice("target/vela-docker"),
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_RUNTIME_VOLUMES", "RUNTIME_VOLUMES"},
		FilePath: "/vela/runtime/volumes",
		Name:     "runtime.volumes",
		Usage:    "list of host volumes to mount for the runtime",
	},
}
