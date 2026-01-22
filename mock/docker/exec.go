// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// ExecService implements all the exec
// related functions for the Docker mock.
type ExecService struct{}

// ExecCreate is a helper function to simulate
// a mocked call to create a Docker exec instance.
func (e *ExecService) ExecCreate(_ context.Context, _ string, _ client.ExecCreateOptions) (client.ExecCreateResult, error) {
	return client.ExecCreateResult{}, nil
}

// ExecInspect is a helper function to simulate
// a mocked call to inspect a Docker exec instance.
func (e *ExecService) ExecInspect(_ context.Context, _ string, _ client.ExecInspectOptions) (client.ExecInspectResult, error) {
	return client.ExecInspectResult{}, nil
}

// ExecResize is a helper function to simulate
// a mocked call to resize a Docker exec instance.
func (e *ExecService) ExecResize(_ context.Context, _ string, _ client.ExecResizeOptions) (client.ExecResizeResult, error) {
	return client.ExecResizeResult{}, nil
}

// ExecStart is a helper function to simulate
// a mocked call to start a Docker exec instance.
func (e *ExecService) ExecStart(_ context.Context, _ string, _ client.ExecStartOptions) (client.ExecStartResult, error) {
	return client.ExecStartResult{}, nil
}

// ExecAttach is a helper function to simulate
// a mocked call to attach to a Docker exec instance.
func (e *ExecService) ExecAttach(_ context.Context, _ string, _ client.ExecAttachOptions) (client.ExecAttachResult, error) {
	return client.ExecAttachResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ExecService satisfies the ExecAPIClient interface that
// the Docker client expects.
var _ client.ExecAPIClient = (*ExecService)(nil)
