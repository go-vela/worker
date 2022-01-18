// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"

	"github.com/go-vela/types/pipeline"
)

// CreateNetwork creates the pipeline network.
// This is a no-op for podman. Using pod-scoped networking.
func (c *client) CreateNetwork(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: creating network for pipeline %s", b.ID)

	return nil
}

// InspectNetwork inspects the pipeline network.
// This is a no-op for podman.
func (c *client) InspectNetwork(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("no-op: inspecting network for pipeline %s", b.ID)

	return []byte{}, nil
}

// RemoveNetwork deletes the pipeline network.
// This is a no-op for podman.
func (c *client) RemoveNetwork(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: removing network for pipeline %s", b.ID)

	return nil
}
