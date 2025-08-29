// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/go-units"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	vol "github.com/go-vela/worker/internal/volume"
)

// CreateVolume creates the pipeline volume.
func (c *client) CreateVolume(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("creating volume for pipeline %s", b.ID)

	// create options for creating volume
	//
	// https://pkg.go.dev/github.com/docker/docker/api/types/volume#CreateOptions
	opts := volume.CreateOptions{
		Name:   b.ID,
		Driver: "local",
	}

	// send API call to create the volume
	//
	// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeCreate
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
	// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeInspect
	v, err := c.Docker.VolumeInspect(ctx, b.ID)
	if err != nil {
		return output, err
	}

	// convert volume type Volume to bytes with pretty print
	//
	// https://pkg.go.dev/github.com/docker/docker/api/types#Volume
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
	// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeRemove
	err := c.Docker.VolumeRemove(ctx, b.ID, true)
	if err != nil {
		return err
	}

	return nil
}

// ResourceLimits represents configurable resource limits for containers.
type ResourceLimits struct {
	Memory    int64 // Memory limit in bytes
	CPUQuota  int64 // CPU quota in millicores * 1000
	CPUPeriod int64 // CPU period
	PidsLimit int64 // Process limit
}

// hostConfig is a helper function to generate the host config
// with Ulimit and volume specifications for a container.
func hostConfig(logger *logrus.Entry, id string, ulimits pipeline.UlimitSlice, volumes []string, dropCaps []string, resourceLimits *ResourceLimits) *container.HostConfig {
	logger.Tracef("creating mount for default volume %s", id)

	// create default mount for pipeline volume
	mounts := []mount.Mount{
		{
			Type:   mount.TypeVolume,
			Source: id,
			Target: constants.WorkspaceMount,
		},
	}

	// Security hardening: Apply container resource limits and security constraints
	// Use provided resource limits or fallback to secure defaults
	var memory, cpuQuota, cpuPeriod, pidsLimit int64
	if resourceLimits != nil {
		memory = resourceLimits.Memory
		cpuQuota = resourceLimits.CPUQuota
		cpuPeriod = resourceLimits.CPUPeriod
		pidsLimit = resourceLimits.PidsLimit
	} else {
		// Secure defaults
		memory = int64(4) * 1024 * 1024 * 1024 // 4GB limit
		cpuQuota = int64(1.2 * 100000)         // 1.2 CPU cores
		cpuPeriod = 100000                     // Standard period
		pidsLimit = 1024                       // Prevent fork bombs
	}

	resources := container.Resources{
		Memory:    memory,
		CPUQuota:  cpuQuota,
		CPUPeriod: cpuPeriod,
		PidsLimit: &pidsLimit,
	}

	// Apply default security ulimits if none provided
	if len(ulimits) == 0 {
		resources.Ulimits = []*units.Ulimit{
			{Name: "nofile", Hard: 1024, Soft: 1024}, // File descriptors
			{Name: "nproc", Hard: 512, Soft: 512},    // Process limit
		}
	} else {
		// iterate through all ulimits provided
		for _, v := range ulimits {
			resources.Ulimits = append(resources.Ulimits, &units.Ulimit{
				Name: v.Name,
				Hard: v.Hard,
				Soft: v.Soft,
			})
		}
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

	// Ensure dropCaps includes ALL capabilities if empty (security hardening)
	if len(dropCaps) == 0 {
		dropCaps = []string{"ALL"}
	}

	// https://pkg.go.dev/github.com/docker/docker/api/types/container#HostConfig
	return &container.HostConfig{
		// https://pkg.go.dev/github.com/docker/docker/api/types/container#LogConfig
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: false,
		// https://pkg.go.dev/github.com/docker/docker/api/types/mount#Mount
		Mounts: mounts,
		// https://pkg.go.dev/github.com/docker/docker/api/types/container#Resources.Ulimits
		Resources: resources,
		// Security hardening: Drop all capabilities by default, add only essential ones
		CapDrop: dropCaps,
		CapAdd:  []string{"CHOWN", "SETUID", "SETGID"}, // Essential capabilities only
		// Security options to prevent privilege escalation
		SecurityOpt: []string{
			"no-new-privileges:true", // Prevent privilege escalation
			"seccomp=docker/default", // Apply seccomp filtering
		},
		ReadonlyRootfs: false, // Start with false, enable per-container as feasible
	}
}
