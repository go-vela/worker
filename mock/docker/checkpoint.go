// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// CheckpointService implements all the checkpoint
// related functions for the Docker mock.
type CheckpointService struct{}

func (cp *CheckpointService) CheckpointCreate(_ context.Context, _ string, _ client.CheckpointCreateOptions) (client.CheckpointCreateResult, error) {
	return client.CheckpointCreateResult{}, nil
}

func (cp *CheckpointService) CheckpointRemove(_ context.Context, _ string, _ client.CheckpointRemoveOptions) (client.CheckpointRemoveResult, error) {
	return client.CheckpointRemoveResult{}, nil
}

func (cp *CheckpointService) CheckpointList(_ context.Context, _ string, _ client.CheckpointListOptions) (client.CheckpointListResult, error) {
	return client.CheckpointListResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ImageService satisfies the ImageAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ConfigAPIClient
var _ client.CheckpointAPIClient = (*CheckpointService)(nil)
