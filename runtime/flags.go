// SPDX-License-Identifier: Apache-2.0

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
	&cli.StringFlag{
		EnvVars:  []string{"VELA_RUNTIME_PODS_TEMPLATE_NAME", "RUNTIME_PODS_TEMPLATE_NAME"},
		FilePath: "/vela/runtime/pods_template_name",
		Name:     "runtime.pods-template-name",
		Usage:    "name of the PipelinePodsTemplate to retrieve from the runtime.namespace (only used by kubernetes)",
	},
	&cli.PathFlag{
		EnvVars:  []string{"VELA_RUNTIME_PODS_TEMPLATE_FILE", "RUNTIME_PODS_TEMPLATE_FILE"},
		FilePath: "/vela/runtime/pods_template_file",
		Name:     "runtime.pods-template-file",
		Usage:    "path to local fallback file containing a PipelinePodsTemplate in YAML (only used by kubernetes; only used if runtime.pods-template-name is not defined)",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_RUNTIME_PRIVILEGED_IMAGES", "RUNTIME_PRIVILEGED_IMAGES"},
		FilePath: "/vela/runtime/privileged_images",
		Name:     "runtime.privileged-images",
		Usage:    "list of images allowed to run in privileged mode for the runtime",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_RUNTIME_VOLUMES", "RUNTIME_VOLUMES"},
		FilePath: "/vela/runtime/volumes",
		Name:     "runtime.volumes",
		Usage:    "list of host volumes to mount for the runtime",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_RUNTIME_DROP_CAPABILITIES", "RUNTIME_DROP_CAPABILITIES"},
		FilePath: "/vela/runtime/drop_capabilities",
		Name:     "runtime.drop-capabilities",
		Usage:    "list of kernel capabilities to drop from container privileges (only used by Docker)",
	},
}
