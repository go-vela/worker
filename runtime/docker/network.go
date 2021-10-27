// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"

	"github.com/go-vela/types/pipeline"

	"github.com/sirupsen/logrus"
)

// CreateNetwork creates the pipeline network.
func (c *client) CreateNetwork(ctx context.Context, b *pipeline.Build) error {
	logrus.Tracef("creating network for pipeline %s", b.ID)

	// create options for creating network
	//
	// https://godoc.org/github.com/docker/docker/api/types#NetworkCreate
	opts := types.NetworkCreate{
		Driver: "bridge",
	}

	// send API call to create the network
	//
	// https://godoc.org/github.com/docker/docker/client#Client.NetworkCreate
	_, err := c.Docker.NetworkCreate(ctx, b.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

// InspectNetwork inspects the pipeline network.
func (c *client) InspectNetwork(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	logrus.Tracef("inspecting network for pipeline %s", b.ID)

	// create options for inspecting network
	//
	// https://godoc.org/github.com/docker/docker/api/types#NetworkInspectOptions
	opts := types.NetworkInspectOptions{}

	// create output for inspecting network
	output := []byte(
		fmt.Sprintf("$ docker network inspect %s\n", b.ID),
	)

	// send API call to inspect the network
	//
	// https://godoc.org/github.com/docker/docker/client#Client.NetworkInspect
	n, err := c.Docker.NetworkInspect(ctx, b.ID, opts)
	if err != nil {
		return output, err
	}

	// convert network type NetworkResource to bytes with pretty print
	//
	// https://godoc.org/github.com/docker/docker/api/types#NetworkResource
	network, err := json.MarshalIndent(n, "", " ")
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, append(network, "\n"...)...), nil
}

// RemoveNetwork deletes the pipeline network.
func (c *client) RemoveNetwork(ctx context.Context, b *pipeline.Build) error {
	logrus.Tracef("removing network for pipeline %s", b.ID)

	// send API call to remove the network
	//
	// https://godoc.org/github.com/docker/docker/client#Client.NetworkRemove
	err := c.Docker.NetworkRemove(ctx, b.ID)
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
	//
	// https://godoc.org/github.com/docker/docker/api/types/network#EndpointSettings
	endpoints[id] = &network.EndpointSettings{
		NetworkID: id,
		Aliases:   []string{alias},
	}

	// return network config with configured endpoints
	//
	// https://godoc.org/github.com/docker/docker/api/types/network#NetworkingConfig
	return &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}
