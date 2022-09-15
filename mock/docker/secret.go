// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code
package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SecretService implements all the secret
// related functions for the Docker mock.
type SecretService struct{}

// SecretCreate is a helper function to simulate
// a mocked call to create a secret for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SecretCreate
func (s *SecretService) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (types.SecretCreateResponse, error) {
	return types.SecretCreateResponse{}, nil
}

// SecretInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker secret and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SecretInspectWithRaw
func (s *SecretService) SecretInspectWithRaw(ctx context.Context, name string) (swarm.Secret, []byte, error) {
	return swarm.Secret{}, nil, nil
}

// SecretList is a helper function to simulate
// a mocked call to list secrets for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SecretList
func (s *SecretService) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
	return nil, nil
}

// SecretRemove is a helper function to simulate
// a mocked call to remove a secret for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SecretRemove
func (s *SecretService) SecretRemove(ctx context.Context, id string) error {
	return nil
}

// SecretUpdate is a helper function to simulate
// a mocked call to update a secret for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.SecretUpdate
func (s *SecretService) SecretUpdate(ctx context.Context, id string, version swarm.Version, secret swarm.SecretSpec) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// SecretService satisfies the SecretAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#SecretAPIClient
var _ client.SecretAPIClient = (*SecretService)(nil)
