// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/bindings/volumes"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	vol "github.com/go-vela/worker/internal/volume"
)

// CreateVolume creates the pipeline volume.
func (c *client) CreateVolume(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("creating volume for pipeline %s", b.ID)

	// create options for creating volume
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/domain/entities#VolumeCreateOptions
	opts := entities.VolumeCreateOptions{
		Name:   b.ID,
		Driver: define.VolumeDriverLocal,
	}

	// send API call to create the volume
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/volumes#Create
	_, err := volumes.Create(c.Podman, opts, &volumes.CreateOptions{})
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
		fmt.Sprintf("$ podman volume inspect %s\n", b.ID),
	)

	// send API call to inspect the volume
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/volumes#Inspect
	v, err := volumes.Inspect(c.Podman, b.ID, &volumes.InspectOptions{})
	if err != nil {
		return output, err
	}

	// convert volume type Volume to bytes with pretty print
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/domain/entities#VolumeConfigResponse
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

	// remove options to use for the call to remove the volume
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/volumes#RemoveOptions
	rmOpts := new(volumes.RemoveOptions).WithForce(true)

	// send API call to remove the volume
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/volumes#Remove
	err := volumes.Remove(c.Podman, b.ID, rmOpts)
	if err != nil {
		return err
	}

	return nil
}

// storageConfig is a helper function to configure
// volumes and mounts for a container
func storageConfig(ctn *pipeline.Container, id string, volumes []string, logger *logrus.Entry) specgen.ContainerStorageConfig {
	storageConfig := specgen.ContainerStorageConfig{}

	storageConfig.Image = ctn.Image
	storageConfig.CreateWorkingDir = true
	storageConfig.WorkDir = ctn.Directory

	// create default volume/mount for pipeline volume
	storageConfig.Volumes = append(storageConfig.Volumes, &specgen.NamedVolume{
		Name: id,
		Dest: constants.WorkspaceMount,
	})

	// check if volumes were provided
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
			storageConfig.Mounts = append(storageConfig.Mounts, specs.Mount{
				Type:        define.TypeBind,
				Source:      _volume.Source,
				Destination: _volume.Destination,
				Options:     strings.Split(_volume.AccessMode, ","),
			})
		}
	}

	return storageConfig
}
