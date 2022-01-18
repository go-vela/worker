// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/containers/podman/v3/pkg/bindings/network"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/go-vela/types/pipeline"
)

// CreateNetwork creates the pipeline network.
func (c *client) CreateNetwork(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("creating network for pipeline %s", b.ID)

	// create options for network
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/network#CreateOptions
	netOpts := new(network.CreateOptions).
		WithName(b.ID)

	// send API call to create the network
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/network#Create
	_, err := network.Create(c.Podman, netOpts)
	if err != nil {
		return err
	}

	return nil
}

// InspectNetwork inspects the pipeline network.
func (c *client) InspectNetwork(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("inspecting network for pipeline %s", b.ID)

	// create output for inspecting network
	output := []byte(
		fmt.Sprintf("$ podman network inspect %s\n", b.ID),
	)

	// send the API call to inspect the network
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/network#Inspect
	n, err := network.Inspect(c.Podman, b.ID, &network.InspectOptions{})
	if err != nil {
		return []byte{}, err
	}

	// convert network type NetworkInspectReport to bytes with pretty print
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/domain/entities#NetworkInspectReport
	network, err := json.MarshalIndent(n, "", " ")
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, append(network, "\n"...)...), nil
}

// RemoveNetwork deletes the pipeline network.
func (c *client) RemoveNetwork(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("removing network for pipeline %s", b.ID)

	// remove options for removing the network
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/network#RemoveOptions
	rmOpts := new(network.RemoveOptions).WithForce(true)

	// send the API call to remove the network
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/network#Remove
	_, err := network.Remove(c.Podman, b.ID, rmOpts)
	if err != nil {
		return err
	}

	return nil
}

// netConfig is a helper function to generate
// the network config for a container.
func netConfig(id, alias string) specgen.ContainerNetworkConfig {
	netConf := specgen.ContainerNetworkConfig{}

	// the network to join
	netConf.CNINetworks = []string{id}

	// alias to assign for the given network
	netConf.Aliases = map[string][]string{
		id: []string{alias},
	}

	return netConf
}
