// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/go-units"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	vol "github.com/go-vela/worker/internal/volume"

	"github.com/sirupsen/logrus"
)

// CreateVolume creates the pipeline volume.
func (c *client) CreateVolume(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("creating volume for pipeline %s", b.ID)

	// create options for creating volume
	//
	// https://godoc.org/github.com/docker/docker/api/types/volume#VolumeCreateBody
	opts := volume.VolumeCreateBody{
		Name:   b.ID,
		Driver: "local",
	}

	// send API call to create the volume
	//
	// https://godoc.org/github.com/docker/docker/client#Client.VolumeCreate
	_, err := c.Docker.VolumeCreate(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}

// InspectVolume inspects the pipeline volume.
func (c *client) InspectVolume(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("inspecting volume for pipeline %s", b.ID)

	// create output for inspecting volume
	output := []byte(
		fmt.Sprintf("$ docker volume inspect %s\n", b.ID),
	)

	// send API call to inspect the volume
	//
	// https://godoc.org/github.com/docker/docker/client#Client.VolumeInspect
	v, err := c.Docker.VolumeInspect(ctx, b.ID)
	if err != nil {
		return output, err
	}

	// convert volume type Volume to bytes with pretty print
	//
	// https://godoc.org/github.com/docker/docker/api/types#Volume
	volume, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, append(volume, "\n"...)...), nil
}

// RemoveVolume deletes the pipeline volume.
func (c *client) RemoveVolume(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("removing volume for pipeline %s", b.ID)

	// send API call to remove the volume
	//
	// https://godoc.org/github.com/docker/docker/client#Client.VolumeRemove
	err := c.Docker.VolumeRemove(ctx, b.ID, true)
	if err != nil {
		return err
	}

	return nil
}

// hostConfig is a helper function to generate the host config
// with Ulimit and volume specifications for a container.
func hostConfig(logger *logrus.Entry, id string, ulimits pipeline.UlimitSlice, volumes []string, dropCaps []string) *container.HostConfig {
	logger.Tracef("creating mount for default volume %s", id)

	// create default mount for pipeline volume
	mounts := []mount.Mount{
		{
			Type:   mount.TypeVolume,
			Source: id,
			Target: constants.WorkspaceMount,
		},
	}

	resources := container.Resources{}
	// iterate through all ulimits provided

	for _, v := range ulimits {
		resources.Ulimits = append(resources.Ulimits, &units.Ulimit{
			Name: v.Name,
			Hard: v.Hard,
			Soft: v.Soft,
		})
	}

	// check if other volumes were provided
	if len(volumes) > 0 {
		// iterate through all volumes provided
		for _, v := range volumes {
			logger.Tracef("creating mount for volume %s", v)

			// parse the volume provided
			_volume, err := vol.ParseWithError(v)
			if err != nil {
				logger.Error(err)
			}

			// add the volume to the set of mounts
			mounts = append(mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   _volume.Source,
				Target:   _volume.Destination,
				ReadOnly: _volume.AccessMode == "ro",
			})
		}
	}

	// https://godoc.org/github.com/docker/docker/api/types/container#HostConfig
	return &container.HostConfig{
		// https://godoc.org/github.com/docker/docker/api/types/container#LogConfig
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: false,
		// https://godoc.org/github.com/docker/docker/api/types/mount#Mount
		Mounts: mounts,
		// https://pkg.go.dev/github.com/docker/docker/api/types/container#Resources.Ulimits
		Resources: resources,
		CapDrop:   dropCaps,
	}
}
