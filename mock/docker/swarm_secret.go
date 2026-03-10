// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// SecretService implements all the secret
// related functions for the Docker mock.
type SecretService struct{}

// SecretCreate is a helper function to simulate
// a mocked call to create a secret for a
// Docker swarm cluster.
func (s *SecretService) SecretCreate(_ context.Context, _ client.SecretCreateOptions) (client.SecretCreateResult, error) {
	return client.SecretCreateResult{}, nil
}

// SecretInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker secret and return
// the raw body received from the API.
func (s *SecretService) SecretInspect(_ context.Context, _ string, _ client.SecretInspectOptions) (client.SecretInspectResult, error) {
	return client.SecretInspectResult{}, nil
}

// SecretList is a helper function to simulate
// a mocked call to list secrets for a
// Docker swarm cluster.
func (s *SecretService) SecretList(_ context.Context, _ client.SecretListOptions) (client.SecretListResult, error) {
	return client.SecretListResult{}, nil
}

// SecretUpdate is a helper function to simulate
// a mocked call to update a secret for a
// Docker swarm cluster.
func (s *SecretService) SecretUpdate(_ context.Context, _ string, _ client.SecretUpdateOptions) (client.SecretUpdateResult, error) {
	return client.SecretUpdateResult{}, nil
}

// SecretRemove is a helper function to simulate
// a mocked call to remove a secret for a
// Docker swarm cluster.
func (s *SecretService) SecretRemove(_ context.Context, _ string, _ client.SecretRemoveOptions) (client.SecretRemoveResult, error) {
	return client.SecretRemoveResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// service satisfies the APIClient interface that
// the Docker client expects.
var _ client.SecretAPIClient = (*SecretService)(nil)
