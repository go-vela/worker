// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// SwarmService implements all the swarm
// related functions for the Docker mock.
type SwarmService struct{}

// SwarmInit is a helper function to simulate
// a mocked call to initialize the Docker
// swarm cluster.
func (s *SwarmService) SwarmInit(_ context.Context, _ client.SwarmInitOptions) (client.SwarmInitResult, error) {
	return client.SwarmInitResult{}, nil
}

// SwarmJoin is a helper function to simulate
// a mocked call to join the Docker swarm
// cluster.
func (s *SwarmService) SwarmJoin(_ context.Context, _ client.SwarmJoinOptions) (client.SwarmJoinResult, error) {
	return client.SwarmJoinResult{}, nil
}

// SwarmInspect is a helper function to simulate
// a mocked call to inspect the Docker swarm
// cluster.
func (s *SwarmService) SwarmInspect(_ context.Context, _ client.SwarmInspectOptions) (client.SwarmInspectResult, error) {
	return client.SwarmInspectResult{}, nil
}

// SwarmUpdate is a helper function to simulate
// a mocked call to update the Docker swarm
// cluster.
func (s *SwarmService) SwarmUpdate(_ context.Context, _ client.SwarmUpdateOptions) (client.SwarmUpdateResult, error) {
	return client.SwarmUpdateResult{}, nil
}

// SwarmLeave is a helper function to simulate
// a mocked call to leave the Docker swarm
// cluster.
func (s *SwarmService) SwarmLeave(_ context.Context, _ client.SwarmLeaveOptions) (client.SwarmLeaveResult, error) {
	return client.SwarmLeaveResult{}, nil
}

// SwarmGetUnlockKey is a helper function to simulate
// a mocked call to capture the unlock key for a
// Docker swarm cluster.
func (s *SwarmService) SwarmGetUnlockKey(_ context.Context) (client.SwarmGetUnlockKeyResult, error) {
	return client.SwarmGetUnlockKeyResult{}, nil
}

// SwarmUnlock is a helper function to simulate
// a mocked call to unlock the Docker swarm
// cluster.
func (s *SwarmService) SwarmUnlock(_ context.Context, _ client.SwarmUnlockOptions) (client.SwarmUnlockResult, error) {
	return client.SwarmUnlockResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// service satisfies the APIClient interface that
// the Docker client expects.
var _ client.SwarmAPIClient = (*SwarmService)(nil)
