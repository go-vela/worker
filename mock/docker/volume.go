// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/volume"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/stringid"
)

// VolumeService implements all the volume
// related functions for the Docker mock.
type VolumeService struct{}

// VolumeCreate is a helper function to simulate
// a mocked call to create a Docker volume.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeCreate
func (v *VolumeService) VolumeCreate(_ context.Context, options client.VolumeCreateOptions) (client.VolumeCreateResult, error) {
	// verify a volume was provided
	if len(options.Name) == 0 {
		return client.VolumeCreateResult{}, errors.New("no volume provided")
	}

	// check if the volume is notfound and
	// check if the notfound should be ignored
	if strings.Contains(options.Name, "notfound") &&
		!strings.Contains(options.Name, "ignorenotfound") {
		return client.VolumeCreateResult{}, cerrdefs.ErrNotFound
	}

	// check if the volume is not-found and
	// check if the not-found should be ignored
	if strings.Contains(options.Name, "not-found") &&
		!strings.Contains(options.Name, "ignore-not-found") {
		return client.VolumeCreateResult{}, cerrdefs.ErrNotFound
	}

	// create response object to return
	response := client.VolumeCreateResult{
		Volume: volume.Volume{
			CreatedAt:  time.Now().String(),
			Driver:     options.Driver,
			Labels:     options.Labels,
			Mountpoint: fmt.Sprintf("/var/lib/docker/volumes/%s/_data", stringid.GenerateRandomID()),
			Name:       options.Name,
			Options:    options.DriverOpts,
			Scope:      "local",
		},
	}

	return response, nil
}

// VolumeInspect is a helper function to simulate
// a mocked call to inspect a Docker volume.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeInspect
func (v *VolumeService) VolumeInspect(_ context.Context, volumeID string, _ client.VolumeInspectOptions) (client.VolumeInspectResult, error) {
	// verify a volume was provided
	if len(volumeID) == 0 {
		return client.VolumeInspectResult{}, errors.New("no volume provided")
	}

	// check if the volume is notfound
	if strings.Contains(volumeID, "notfound") {
		return client.VolumeInspectResult{}, cerrdefs.ErrNotFound
	}

	// check if the volume is not-found
	if strings.Contains(volumeID, "not-found") {
		return client.VolumeInspectResult{}, cerrdefs.ErrNotFound
	}

	// create response object to return
	response := client.VolumeInspectResult{
		Volume: volume.Volume{
			CreatedAt:  time.Now().String(),
			Driver:     "local",
			Mountpoint: fmt.Sprintf("/var/lib/docker/volumes/%s/_data", stringid.GenerateRandomID()),
			Name:       volumeID,
			Scope:      "local",
		},
	}

	return response, nil
}

// VolumeList is a helper function to simulate
// a mocked call to list Docker volumes.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeList
func (v *VolumeService) VolumeList(_ context.Context, _ client.VolumeListOptions) (client.VolumeListResult, error) {
	return client.VolumeListResult{}, nil
}

// VolumeUpdate is a helper function to simulate
// a mocked call to update Docker volumes.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeUpdate
func (v *VolumeService) VolumeUpdate(_ context.Context, _ string, _ client.VolumeUpdateOptions) (client.VolumeUpdateResult, error) {
	return client.VolumeUpdateResult{}, nil
}

// VolumeRemove is a helper function to simulate
// a mocked call to remove Docker a volume.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeRemove
func (v *VolumeService) VolumeRemove(_ context.Context, volumeID string, _ client.VolumeRemoveOptions) (client.VolumeRemoveResult, error) {
	// verify a volume was provided
	if len(volumeID) == 0 {
		return client.VolumeRemoveResult{}, errors.New("no volume provided")
	}

	return client.VolumeRemoveResult{}, nil
}

// VolumesPrune is a helper function to simulate
// a mocked call to prune Docker volumes.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumesPrune
func (v *VolumeService) VolumePrune(_ context.Context, _ client.VolumePruneOptions) (client.VolumePruneResult, error) {
	return client.VolumePruneResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// VolumeService satisfies the VolumeAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#VolumeAPIClient
var _ client.VolumeAPIClient = (*VolumeService)(nil)
