// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/docker/docker/api/types/checkpoint"
)

// CheckpointService implements all the checkpoint
// related functions for the Docker mock.
type CheckpointService struct{}

func (cp *CheckpointService) CheckpointCreate(_ context.Context, _ string, _ checkpoint.CreateOptions) error {
	return nil
}

func (cp *CheckpointService) CheckpointDelete(_ context.Context, _ string, _ checkpoint.DeleteOptions) error {
	return nil
}

func (cp *CheckpointService) CheckpointList(_ context.Context, _ string, _ checkpoint.ListOptions) ([]checkpoint.Summary, error) {
	return nil, nil
}
