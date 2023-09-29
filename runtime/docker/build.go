// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/go-vela/types/pipeline"
)

// InspectBuild displays details about the pod for the init step.
// This is a no-op for docker.
func (c *client) InspectBuild(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("no-op: inspecting build for pipeline %s", b.ID)

	return []byte{}, nil
}

// SetupBuild prepares the pipeline build.
// This is a no-op for docker.
func (c *client) SetupBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: setting up for build %s", b.ID)

	return nil
}

// StreamBuild initializes log/event streaming for build.
// This is a no-op for docker.
func (c *client) StreamBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: streaming build %s", b.ID)

	return nil
}

// AssembleBuild finalizes pipeline build setup.
// This is a no-op for docker.
func (c *client) AssembleBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: assembling build %s", b.ID)

	return nil
}

// RemoveBuild deletes (kill, remove) the pipeline build metadata.
// This is a no-op for docker.
func (c *client) RemoveBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: removing build %s", b.ID)

	return nil
}
