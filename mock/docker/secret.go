// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package docker

import (
	"context"

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
// https://pkg.go.dev/github.com/docker/docker/client#Client.SecretCreate
func (s *SecretService) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (swarm.SecretCreateResponse, error) {
	return swarm.SecretCreateResponse{}, nil
}

// SecretInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker secret and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SecretInspectWithRaw
func (s *SecretService) SecretInspectWithRaw(ctx context.Context, name string) (swarm.Secret, []byte, error) {
	return swarm.Secret{}, nil, nil
}

// SecretList is a helper function to simulate
// a mocked call to list secrets for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SecretList
func (s *SecretService) SecretList(ctx context.Context, options swarm.SecretListOptions) ([]swarm.Secret, error) {
	return nil, nil
}

// SecretRemove is a helper function to simulate
// a mocked call to remove a secret for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SecretRemove
func (s *SecretService) SecretRemove(ctx context.Context, id string) error {
	return nil
}

// SecretUpdate is a helper function to simulate
// a mocked call to update a secret for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.SecretUpdate
func (s *SecretService) SecretUpdate(ctx context.Context, id string, version swarm.Version, secret swarm.SecretSpec) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// SecretService satisfies the SecretAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#SecretAPIClient
var _ client.SecretAPIClient = (*SecretService)(nil)
