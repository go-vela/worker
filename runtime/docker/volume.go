// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"

	types "github.com/docker/docker/api/types/volume"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// CreateVolume creates the pipeline volume.
func (c *client) CreateVolume(ctx context.Context, b *pipeline.Build) error {
	logrus.Tracef("Creating volume for pipeline %s", b.ID)

	// create options for creating volume
	opts := types.VolumeCreateBody{
		Name:   b.ID,
		Driver: "local",
	}

	// send API call to create the volume
	_, err := c.Runtime.VolumeCreate(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}

// InspectVolume inspects the pipeline volume.
func (c *client) InspectVolume(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	logrus.Tracef("Inspecting volume for pipeline %s", b.ID)

	// send API call to inspect the volume
	v, err := c.Runtime.VolumeInspect(ctx, b.ID)
	if err != nil {
		return nil, err
	}

	return []byte(v.Name + "\n"), nil
}

// RemoveVolume deletes the pipeline volume.
func (c *client) RemoveVolume(ctx context.Context, b *pipeline.Build) error {
	logrus.Tracef("Removing volume for pipeline %s", b.ID)

	// send API call to remove the volume
	err := c.Runtime.VolumeRemove(ctx, b.ID, true)
	if err != nil {
		return err
	}

	return nil
}
