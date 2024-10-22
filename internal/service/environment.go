// SPDX-License-Identifier: Apache-2.0

package service

import (
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Environment attempts to update the environment variables
// for the container based off the library resources.
func Environment(c *pipeline.Container, b *api.Build, s *api.Service, version string) error {
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

		// populate environment variables from build library
		err := c.MergeEnv(b.Environment(workspace, channel))
		if err != nil {
			return err
		}
	}

	// populate environment variables from repo library
	err := c.MergeEnv(b.GetRepo().Environment())
	if err != nil {
		return err
	}

	// check if the service provided is empty
	if s != nil {
		// populate environment variables from service library
		err := c.MergeEnv(s.Environment())
		if err != nil {
			return err
		}
	}

	return nil
}
