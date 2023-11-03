// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stringid"
)

// VolumeService implements all the volume
// related functions for the Docker mock.
type VolumeService struct{}

// VolumeCreate is a helper function to simulate
// a mocked call to create a Docker volume.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeCreate
func (v *VolumeService) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error) {
	// verify a volume was provided
	if len(options.Name) == 0 {
		return volume.Volume{}, errors.New("no volume provided")
	}

	// check if the volume is notfound and
	// check if the notfound should be ignored
	if strings.Contains(options.Name, "notfound") &&
		!strings.Contains(options.Name, "ignorenotfound") {
		return volume.Volume{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such volume: %s", options.Name))
	}

	// check if the volume is not-found and
	// check if the not-found should be ignored
	if strings.Contains(options.Name, "not-found") &&
		!strings.Contains(options.Name, "ignore-not-found") {
		return volume.Volume{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such volume: %s", options.Name))
	}

	// create response object to return
	response := volume.Volume{
		CreatedAt:  time.Now().String(),
		Driver:     options.Driver,
		Labels:     options.Labels,
		Mountpoint: fmt.Sprintf("/var/lib/docker/volumes/%s/_data", stringid.GenerateRandomID()),
		Name:       options.Name,
		Options:    options.DriverOpts,
		Scope:      "local",
	}

	return response, nil
}

// VolumeInspect is a helper function to simulate
// a mocked call to inspect a Docker volume.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeInspect
func (v *VolumeService) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error) {
	// verify a volume was provided
	if len(volumeID) == 0 {
		return volume.Volume{}, errors.New("no volume provided")
	}

	// check if the volume is notfound
	if strings.Contains(volumeID, "notfound") {
		return volume.Volume{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such volume: %s", volumeID))
	}

	// check if the volume is not-found
	if strings.Contains(volumeID, "not-found") {
		return volume.Volume{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such volume: %s", volumeID))
	}

	// create response object to return
	response := volume.Volume{
		CreatedAt:  time.Now().String(),
		Driver:     "local",
		Mountpoint: fmt.Sprintf("/var/lib/docker/volumes/%s/_data", stringid.GenerateRandomID()),
		Name:       volumeID,
		Scope:      "local",
	}

	return response, nil
}

// VolumeInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker volume and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeInspectWithRaw
func (v *VolumeService) VolumeInspectWithRaw(ctx context.Context, volumeID string) (volume.Volume, []byte, error) {
	// verify a volume was provided
	if len(volumeID) == 0 {
		return volume.Volume{}, nil, errors.New("no volume provided")
	}

	// check if the volume is notfound
	if strings.Contains(volumeID, "notfound") {
		return volume.Volume{}, nil,
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such volume: %s", volumeID))
	}

	// check if the volume is not-found
	if strings.Contains(volumeID, "not-found") {
		return volume.Volume{}, nil,
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such volume: %s", volumeID))
	}

	// create response object to return
	response := volume.Volume{
		CreatedAt:  time.Now().String(),
		Driver:     "local",
		Mountpoint: fmt.Sprintf("/var/lib/docker/volumes/%s/_data", stringid.GenerateRandomID()),
		Name:       volumeID,
		Scope:      "local",
	}

	// marshal response into raw bytes
	b, err := json.Marshal(response)
	if err != nil {
		return volume.Volume{}, nil, err
	}

	return response, b, nil
}

// VolumeList is a helper function to simulate
// a mocked call to list Docker volumes.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeList
func (v *VolumeService) VolumeList(ctx context.Context, opts volume.ListOptions) (volume.ListResponse, error) {
	return volume.ListResponse{}, nil
}

// VolumeRemove is a helper function to simulate
// a mocked call to remove Docker a volume.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeRemove
func (v *VolumeService) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	// verify a volume was provided
	if len(volumeID) == 0 {
		return errors.New("no volume provided")
	}

	return nil
}

// VolumesPrune is a helper function to simulate
// a mocked call to prune Docker volumes.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumesPrune
func (v *VolumeService) VolumesPrune(ctx context.Context, pruneFilter filters.Args) (types.VolumesPruneReport, error) {
	return types.VolumesPruneReport{}, nil
}

// VolumesUpdate is a helper function to simulate
// a mocked call to update Docker volumes.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.VolumeUpdate
func (v *VolumeService) VolumeUpdate(ctx context.Context, volumeID string, version swarm.Version, options volume.UpdateOptions) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// VolumeService satisfies the VolumeAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#VolumeAPIClient
var _ client.VolumeAPIClient = (*VolumeService)(nil)
