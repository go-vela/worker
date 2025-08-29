// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

// DistributionService implements all the distribution
// related functions for the Docker mock.
type DistributionService struct{}

// DistributionInspect is a helper function to simulate
// a mocked call to inspect a Docker image and return
// the digest and manifest.
func (d *DistributionService) DistributionInspect(_ context.Context, _ string, _ string) (registry.DistributionInspect, error) {
	return registry.DistributionInspect{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// DistributionService satisfies the DistributionAPIClient interface that
// the Docker client expects.
var _ client.DistributionAPIClient = (*DistributionService)(nil)
