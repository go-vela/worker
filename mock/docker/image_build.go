// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"io"

	"github.com/moby/moby/client"
)

// ImageBuildService implements all the image build
// related functions for the Docker mock.
type ImageBuildService struct{}

// ExecCreate is a helper function to simulate
// a mocked call to create a Docker exec instance.
func (ib *ImageBuildService) ImageBuild(_ context.Context, _ io.Reader, _ client.ImageBuildOptions) (client.ImageBuildResult, error) {
	return client.ImageBuildResult{}, nil
}

// ExecInspect is a helper function to simulate
// a mocked call to inspect a Docker exec instance.
func (ib *ImageBuildService) BuildCachePrune(_ context.Context, _ client.BuildCachePruneOptions) (client.BuildCachePruneResult, error) {
	return client.BuildCachePruneResult{}, nil
}

// ExecResize is a helper function to simulate
// a mocked call to resize a Docker exec instance.
func (ib *ImageBuildService) BuildCancel(_ context.Context, _ string, _ client.BuildCancelOptions) (client.BuildCancelResult, error) {
	return client.BuildCancelResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ExecService satisfies the ExecAPIClient interface that
// the Docker client expects.
var _ client.ImageBuildAPIClient = (*ImageBuildService)(nil)
