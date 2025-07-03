// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SwarmService implements all the swarm
// related functions for the Docker mock.
type SwarmService struct{}

// SwarmGetUnlockKey is a helper function to simulate
// a mocked call to capture the unlock key for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmGetUnlockKey
func (s *SwarmService) SwarmGetUnlockKey(ctx context.Context) (swarm.UnlockKeyResponse, error) {
	return swarm.UnlockKeyResponse{}, nil
}

// SwarmInit is a helper function to simulate
// a mocked call to initialize the Docker
// swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmInit
func (s *SwarmService) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	return "", nil
}

// SwarmInspect is a helper function to simulate
// a mocked call to inspect the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmInspect
func (s *SwarmService) SwarmInspect(ctx context.Context) (swarm.Swarm, error) {
	return swarm.Swarm{}, nil
}

// SwarmJoin is a helper function to simulate
// a mocked call to join the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmJoin
func (s *SwarmService) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error {
	return nil
}

// SwarmLeave is a helper function to simulate
// a mocked call to leave the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmLeave
func (s *SwarmService) SwarmLeave(ctx context.Context, force bool) error {
	return nil
}

// SwarmUnlock is a helper function to simulate
// a mocked call to unlock the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmUnlock
func (s *SwarmService) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error {
	return nil
}

// SwarmUpdate is a helper function to simulate
// a mocked call to update the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SwarmUpdate
func (s *SwarmService) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, flags swarm.UpdateFlags) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// SwarmService satisfies the SwarmAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#SwarmAPIClient
var _ client.SwarmAPIClient = (*SwarmService)(nil)
