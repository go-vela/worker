// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
)

// Flags represents all supported command line
// interface (CLI) flags for the runtime.
//
// https://pkg.go.dev/github.com/urfave/cli#Flag
var Flags = []cli.Flag{
	// Runtime Flags

	&cli.StringFlag{
		Name:  "runtime.driver",
		Usage: "driver to be used for the runtime",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_DRIVER"),
			cli.EnvVar("RUNTIME_DRIVER"),
			cli.File("/vela/runtime/driver"),
		),
		Value: constants.DriverDocker,
	},
	&cli.StringFlag{
		Name:  "runtime.config",
		Usage: "path to configuration file for the runtime",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_CONFIG"),
			cli.EnvVar("RUNTIME_CONFIG"),
			cli.File("/vela/runtime/config"),
		),
	},
	&cli.StringFlag{
		Name:  "runtime.namespace",
		Usage: "namespace to use for the runtime (only used by kubernetes)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_NAMESPACE"),
			cli.EnvVar("RUNTIME_NAMESPACE"),
			cli.File("/vela/runtime/namespace"),
		),
	},
	&cli.StringFlag{
		Name:  "runtime.pods-template-name",
		Usage: "name of the PipelinePodsTemplate to retrieve from the runtime.namespace (only used by kubernetes)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_PODS_TEMPLATE_NAME"),
			cli.EnvVar("RUNTIME_PODS_TEMPLATE_NAME"),
			cli.File("/vela/runtime/pods_template_name"),
		),
	},
	&cli.StringFlag{
		Name:  "runtime.pods-template-file",
		Usage: "path to local fallback file containing a PipelinePodsTemplate in YAML (only used by kubernetes; only used if runtime.pods-template-name is not defined)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_PODS_TEMPLATE_FILE"),
			cli.EnvVar("RUNTIME_PODS_TEMPLATE_FILE"),
			cli.File("/vela/runtime/pods_template_file"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "runtime.privileged-images",
		Usage: "list of images allowed to run in privileged mode for the runtime",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_PRIVILEGED_IMAGES"),
			cli.EnvVar("RUNTIME_PRIVILEGED_IMAGES"),
			cli.File("/vela/runtime/privileged_images"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "runtime.volumes",
		Usage: "list of host volumes to mount for the runtime",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_VOLUMES"),
			cli.EnvVar("RUNTIME_VOLUMES"),
			cli.File("/vela/runtime/volumes"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "runtime.drop-capabilities",
		Usage: "list of kernel capabilities to drop from container privileges (only used by Docker)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_RUNTIME_DROP_CAPABILITIES"),
			cli.EnvVar("RUNTIME_DROP_CAPABILITIES"),
			cli.File("/vela/runtime/drop_capabilities"),
		),
	},
}
