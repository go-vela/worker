// SPDX-License-Identifier: Apache-2.0

package step

import (
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/types/constants"
)

// Environment attempts to update the environment variables
// for the container based off the library resources.
func Environment(c *pipeline.Container, b *api.Build, s *api.Step, version, reqToken string) error {
	// check if container or container environment are empty
	if c == nil || c.Environment == nil {
		return fmt.Errorf("empty container provided for environment")
	}

	// check if the build provided is empty
	if b != nil {
		// check if the channel exists in the environment
		channel, ok := c.Environment["VELA_CHANNEL"]
		if !ok {
			// set default for channel
			channel = constants.DefaultRoute
		}

		// check if the workspace exists in the environment
		workspace, ok := c.Environment["VELA_WORKSPACE"]
		if !ok {
			// set default for workspace
			workspace = constants.WorkspaceDefault
		}

		// update environment variables
		c.Environment["VELA_DISTRIBUTION"] = b.GetDistribution()
		c.Environment["VELA_HOST"] = b.GetHost()
		c.Environment["VELA_RUNTIME"] = b.GetRuntime()
		c.Environment["VELA_VERSION"] = version
		c.Environment["VELA_ID_TOKEN_REQUEST_TOKEN"] = reqToken
		c.Environment["VELA_OUTPUTS"] = "/vela/outputs/.env"
		c.Environment["VELA_MASKED_OUTPUTS"] = "/vela/outputs/masked.env"

		// populate environment variables from build library
		//
		// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.MergeEnv
		// ->
		// https://pkg.go.dev/github.com/go-vela/types/library#Build.Environment
		err := c.MergeEnv(b.Environment(workspace, channel))
		if err != nil {
			return err
		}
	}

	// populate environment variables from build library
	//
	// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.MergeEnv
	err := c.MergeEnv(b.GetRepo().Environment())
	if err != nil {
		return err
	}

	// check if the step provided is empty
	if s != nil {
		// populate environment variables from step library
		//
		// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.MergeEnv
		// ->
		// https://pkg.go.dev/github.com/go-vela/types/library#Service.Environment
		err := c.MergeEnv(s.Environment())
		if err != nil {
			return err
		}
	}

	return nil
}
