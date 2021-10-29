// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
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
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmGetUnlockKey
func (s *SwarmService) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error) {
	return types.SwarmUnlockKeyResponse{}, nil
}

// SwarmInit is a helper function to simulate
// a mocked call to initialize the Docker
// swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmInit
func (s *SwarmService) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	return "", nil
}

// SwarmInspect is a helper function to simulate
// a mocked call to inspect the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmInspect
func (s *SwarmService) SwarmInspect(ctx context.Context) (swarm.Swarm, error) {
	return swarm.Swarm{}, nil
}

// SwarmJoin is a helper function to simulate
// a mocked call to join the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmJoin
func (s *SwarmService) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error {
	return nil
}

// SwarmLeave is a helper function to simulate
// a mocked call to leave the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmLeave
func (s *SwarmService) SwarmLeave(ctx context.Context, force bool) error {
	return nil
}

// SwarmUnlock is a helper function to simulate
// a mocked call to unlock the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmUnlock
func (s *SwarmService) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error {
	return nil
}

// SwarmUpdate is a helper function to simulate
// a mocked call to update the Docker swarm
// cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SwarmUpdate
func (s *SwarmService) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, flags swarm.UpdateFlags) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// SwarmService satisfies the SwarmAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#SwarmAPIClient
var _ client.SwarmAPIClient = (*SwarmService)(nil)
