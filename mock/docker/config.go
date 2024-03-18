// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ConfigService implements all the config
// related functions for the Docker mock.
type ConfigService struct{}

// ConfigCreate is a helper function to simulate
// a mocked call to create a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigCreate(ctx context.Context, config swarm.ConfigSpec) (types.ConfigCreateResponse, error) {
	return types.ConfigCreateResponse{}, nil
}

// ConfigInspectWithRaw is a helper function to simulate
// a mocked call to inspect a config for a Docker swarm
// cluster and return the raw body received from the API.
func (c *ConfigService) ConfigInspectWithRaw(ctx context.Context, name string) (swarm.Config, []byte, error) {
	return swarm.Config{}, nil, nil
}

// ConfigList is a helper function to simulate
// a mocked call to list the configs for a
// Docker swarm cluster.
func (c *ConfigService) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) {
	return nil, nil
}

// ConfigRemove is a helper function to simulate
// a mocked call to remove a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigRemove(ctx context.Context, id string) error { return nil }

// ConfigUpdate is a helper function to simulate
// a mocked call to update a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ImageService satisfies the ImageAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ConfigAPIClient
var _ client.ConfigAPIClient = (*ConfigService)(nil)
