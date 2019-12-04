// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// CreateNetwork creates the pipeline network.
func (c *client) CreateNetwork(ctx context.Context, b *pipeline.Build) error {
	logrus.Tracef("Creating network for pipeline %s", b.ID)

	// create options for creating network
	opts := types.NetworkCreate{
		Driver: "bridge",
	}

	// send API call to create the network
	_, err := c.Runtime.NetworkCreate(ctx, b.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

// InspectNetwork inspects the pipeline network.
func (c *client) InspectNetwork(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	logrus.Tracef("Inspecting network for pipeline %s", b.ID)

	// create options for inspecting network
	opts := types.NetworkInspectOptions{}

	// send API call to inspect the network
	n, err := c.Runtime.NetworkInspect(ctx, b.ID, opts)
	if err != nil {
		return nil, err
	}

	return []byte(n.ID + "\n"), nil
}

// RemoveNetwork deletes the pipeline network.
func (c *client) RemoveNetwork(ctx context.Context, b *pipeline.Build) error {
	logrus.Tracef("Removing volume for pipeline %s", b.ID)

	// send API call to remove the network
	err := c.Runtime.NetworkRemove(ctx, b.ID)
	if err != nil {
		return err
	}

	return nil
}

// netConfig is a helper function to generate
// the network config for a container.
func netConfig(id, alias string) *network.NetworkingConfig {
	endpoints := make(map[string]*network.EndpointSettings)

	// set pipeline id for endpoint with alias
	endpoints[id] = &network.EndpointSettings{
		NetworkID: id,
		Aliases:   []string{alias},
	}

	return &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}
