// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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
func (c *client) CreateVolume(ctx context.Context, b *pipeline.Build) (string, error) {
	logrus.Tracef("Creating volume for pipeline %s", b.ID)

	// create options for creating volume
	opts := types.VolumeCreateBody{
		Name:   b.ID,
		Driver: "local",
	}

	// send API call to create the volume
	v, err := c.Runtime.VolumeCreate(ctx, opts)
	if err != nil {
		return "", err
	}

	return v.Name, nil
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
