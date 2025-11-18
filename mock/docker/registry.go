// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// RegistryService implements all the registry
// related functions for the Docker mock.
type RegistryService struct{}

// ImageSearch is a helper function to simulate
// a mocked call to search for a Docker image.
func (r *RegistryService) ImageSearch(_ context.Context, _ string, _ client.ImageSearchOptions) (client.ImageSearchResult, error) {
	return client.ImageSearchResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// RegistryService satisfies the RegistrySearchClient interface that
// the Docker client expects.
var _ client.RegistrySearchClient = (*RegistryService)(nil)
