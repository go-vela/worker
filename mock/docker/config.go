// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// ConfigService implements all the config
// related functions for the Docker mock.
type ConfigService struct{}

// ConfigCreate is a helper function to simulate
// a mocked call to create a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigCreate(_ context.Context, _ client.ConfigCreateOptions) (client.ConfigCreateResult, error) {
	return client.ConfigCreateResult{}, nil
}

// ConfigInspectWithRaw is a helper function to simulate
// a mocked call to inspect a config for a Docker swarm
// cluster and return the raw body received from the API.
func (c *ConfigService) ConfigInspect(_ context.Context, _ string, _ client.ConfigInspectOptions) (client.ConfigInspectResult, error) {
	return client.ConfigInspectResult{}, nil
}

// ConfigList is a helper function to simulate
// a mocked call to list the configs for a
// Docker swarm cluster.
func (c *ConfigService) ConfigList(_ context.Context, _ client.ConfigListOptions) (client.ConfigListResult, error) {
	return client.ConfigListResult{}, nil
}

// ConfigRemove is a helper function to simulate
// a mocked call to remove a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigRemove(_ context.Context, _ string, _ client.ConfigRemoveOptions) (client.ConfigRemoveResult, error) {
	return client.ConfigRemoveResult{}, nil
}

// ConfigUpdate is a helper function to simulate
// a mocked call to update a config for a
// Docker swarm cluster.
func (c *ConfigService) ConfigUpdate(_ context.Context, _ string, _ client.ConfigUpdateOptions) (client.ConfigUpdateResult, error) {
	return client.ConfigUpdateResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ImageService satisfies the ImageAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ConfigAPIClient
var _ client.ConfigAPIClient = (*ConfigService)(nil)
