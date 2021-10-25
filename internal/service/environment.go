// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Environment attempts to update the environment variables
// for the container based off the library resources.
func Environment(c *pipeline.Container, b *library.Build, r *library.Repo, s *library.Service, version string) error {
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
		//
		// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.MergeEnv
		// ->
		// https://pkg.go.dev/github.com/go-vela/types/library#Build.Environment
		err := c.MergeEnv(b.Environment(workspace, channel))
		if err != nil {
			return err
		}
	}

	// check if the repo provided is empty
	if r != nil {
		// populate environment variables from repo library
		//
		// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.MergeEnv
		// ->
		// https://pkg.go.dev/github.com/go-vela/types/library#Repo.Environment
		err := c.MergeEnv(r.Environment())
		if err != nil {
			return err
		}
	}

	// check if the service provided is empty
	if s != nil {
		// populate environment variables from service library
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
