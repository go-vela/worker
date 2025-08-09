// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

func TestVolumeService_VolumeCreate(t *testing.T) {
	service := &VolumeService{}

	tests := []struct {
		name        string
		options     volume.CreateOptions
		wantErr     bool
		wantErrType error
		wantVolume  bool
	}{
		{
			name: "valid volume",
			options: volume.CreateOptions{
				Name:   "test-volume",
				Driver: "local",
				Labels: map[string]string{"test": "label"},
			},
			wantErr:    false,
			wantVolume: true,
		},
		{
			name:    "empty volume name",
			options: volume.CreateOptions{},
			wantErr: true,
		},
		{
			name: "notfound volume",
			options: volume.CreateOptions{
				Name: "notfound-volume",
			},
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name: "notfound volume with ignore",
			options: volume.CreateOptions{
				Name: "notfound-ignorenotfound",
			},
			wantErr:    false,
			wantVolume: true,
		},
		{
			name: "not-found volume",
			options: volume.CreateOptions{
				Name: "not-found-volume",
			},
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name: "not-found volume with ignore",
			options: volume.CreateOptions{
				Name: "not-found-ignore-not-found",
			},
			wantErr:    false,
			wantVolume: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol, err := service.VolumeCreate(context.Background(), tt.options)

			if tt.wantErr {
				if err == nil {
					t.Errorf("VolumeCreate() error = nil, wantErr %v", tt.wantErr)
				}
				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("VolumeCreate() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("VolumeCreate() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantVolume {
					if vol.Name != tt.options.Name {
						t.Errorf("VolumeCreate() volume.Name = %v, want %v", vol.Name, tt.options.Name)
					}

					if vol.Driver != tt.options.Driver {
						t.Errorf("VolumeCreate() volume.Driver = %v, want %v", vol.Driver, tt.options.Driver)
					}

					if vol.Scope != "local" {
						t.Errorf("VolumeCreate() volume.Scope = %v, want local", vol.Scope)
					}

					if vol.Mountpoint == "" {
						t.Errorf("VolumeCreate() volume.Mountpoint = empty, want generated mountpoint")
					}
				}
			}
		})
	}
}

func TestVolumeService_VolumeInspect(t *testing.T) {
	service := &VolumeService{}

	tests := []struct {
		name        string
		volumeID    string
		wantErr     bool
		wantErrType error
		wantVolume  bool
	}{
		{
			name:       "valid volume",
			volumeID:   "test-volume",
			wantErr:    false,
			wantVolume: true,
		},
		{
			name:     "empty volume ID",
			volumeID: "",
			wantErr:  true,
		},
		{
			name:        "notfound volume",
			volumeID:    "notfound-volume",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name:        "not-found volume",
			volumeID:    "not-found-volume",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol, err := service.VolumeInspect(context.Background(), tt.volumeID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("VolumeInspect() error = nil, wantErr %v", tt.wantErr)
				}
				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("VolumeInspect() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("VolumeInspect() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantVolume {
					if vol.Name != tt.volumeID {
						t.Errorf("VolumeInspect() volume.Name = %v, want %v", vol.Name, tt.volumeID)
					}

					if vol.Driver != "local" {
						t.Errorf("VolumeInspect() volume.Driver = %v, want local", vol.Driver)
					}

					if vol.Scope != "local" {
						t.Errorf("VolumeInspect() volume.Scope = %v, want local", vol.Scope)
					}
				}
			}
		})
	}
}

func TestVolumeService_VolumeInspectWithRaw(t *testing.T) {
	service := &VolumeService{}

	tests := []struct {
		name        string
		volumeID    string
		wantErr     bool
		wantErrType error
		wantVolume  bool
		wantRaw     bool
	}{
		{
			name:       "valid volume",
			volumeID:   "test-volume",
			wantErr:    false,
			wantVolume: true,
			wantRaw:    true,
		},
		{
			name:     "empty volume ID",
			volumeID: "",
			wantErr:  true,
		},
		{
			name:        "notfound volume",
			volumeID:    "notfound-volume",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name:        "not-found volume",
			volumeID:    "not-found-volume",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol, raw, err := service.VolumeInspectWithRaw(context.Background(), tt.volumeID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("VolumeInspectWithRaw() error = nil, wantErr %v", tt.wantErr)
				}
				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("VolumeInspectWithRaw() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("VolumeInspectWithRaw() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantVolume {
					if vol.Name != tt.volumeID {
						t.Errorf("VolumeInspectWithRaw() volume.Name = %v, want %v", vol.Name, tt.volumeID)
					}

					if vol.Driver != "local" {
						t.Errorf("VolumeInspectWithRaw() volume.Driver = %v, want local", vol.Driver)
					}
				}

				if tt.wantRaw {
					if len(raw) == 0 {
						t.Errorf("VolumeInspectWithRaw() raw = empty, want data")
					}

					var unmarshaled volume.Volume
					if err := json.Unmarshal(raw, &unmarshaled); err != nil {
						t.Errorf("VolumeInspectWithRaw() raw data invalid JSON: %v", err)
					}
				}
			}
		})
	}
}

func TestVolumeService_VolumeList(t *testing.T) {
	service := &VolumeService{}
	opts := volume.ListOptions{}

	response, err := service.VolumeList(context.Background(), opts)

	if err != nil {
		t.Errorf("VolumeList() error = %v, want nil", err)
	}

	if response.Volumes != nil {
		t.Errorf("VolumeList() response.Volumes = %v, want nil", response.Volumes)
	}

	if len(response.Warnings) != 0 {
		t.Errorf("VolumeList() response.Warnings = %v, want empty", response.Warnings)
	}
}

func TestVolumeService_VolumeRemove(t *testing.T) {
	service := &VolumeService{}

	tests := []struct {
		name     string
		volumeID string
		wantErr  bool
	}{
		{
			name:     "valid volume",
			volumeID: "test-volume",
			wantErr:  false,
		},
		{
			name:     "empty volume ID",
			volumeID: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.VolumeRemove(context.Background(), tt.volumeID, false)

			if tt.wantErr {
				if err == nil {
					t.Errorf("VolumeRemove() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("VolumeRemove() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestVolumeService_VolumesPrune(t *testing.T) {
	service := &VolumeService{}
	pruneFilters := filters.Args{}

	report, err := service.VolumesPrune(context.Background(), pruneFilters)

	if err != nil {
		t.Errorf("VolumesPrune() error = %v, want nil", err)
	}

	if report.VolumesDeleted != nil {
		t.Errorf("VolumesPrune() report.VolumesDeleted = %v, want nil", report.VolumesDeleted)
	}

	if report.SpaceReclaimed != 0 {
		t.Errorf("VolumesPrune() report.SpaceReclaimed = %v, want 0", report.SpaceReclaimed)
	}
}

func TestVolumeService_VolumeUpdate(t *testing.T) {
	service := &VolumeService{}
	version := swarm.Version{}
	opts := volume.UpdateOptions{}

	err := service.VolumeUpdate(context.Background(), "test-volume", version, opts)

	if err != nil {
		t.Errorf("VolumeUpdate() error = %v, want nil", err)
	}
}

func TestVolumeService_InterfaceCompliance(t *testing.T) {
	var _ client.VolumeAPIClient = (*VolumeService)(nil)
}