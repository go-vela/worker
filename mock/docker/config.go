// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ConfigService implements all the config
// related functions for the Docker mock.
type ConfigService struct{}

// ConfigCreate is a helper function to simulate
// a mocked call to create a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigCreate(_ context.Context, _ swarm.ConfigSpec) (swarm.ConfigCreateResponse, error) {
	return swarm.ConfigCreateResponse{}, nil
}

// ConfigInspectWithRaw is a helper function to simulate
// a mocked call to inspect a config for a Docker swarm
// cluster and return the raw body received from the API.
func (c *ConfigService) ConfigInspectWithRaw(_ context.Context, _ string) (swarm.Config, []byte, error) {
	return swarm.Config{}, nil, nil
}

// ConfigList is a helper function to simulate
// a mocked call to list the configs for a
// Docker swarm cluster.
func (c *ConfigService) ConfigList(_ context.Context, _ swarm.ConfigListOptions) ([]swarm.Config, error) {
	return nil, nil
}

// ConfigRemove is a helper function to simulate
// a mocked call to remove a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigRemove(_ context.Context, _ string) error { return nil }

// ConfigUpdate is a helper function to simulate
// a mocked call to update a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigUpdate(_ context.Context, _ string, _ swarm.Version, _ swarm.ConfigSpec) error {
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
