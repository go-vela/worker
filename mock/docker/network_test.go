// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

func TestNetworkService_NetworkConnect(t *testing.T) {
	service := &NetworkService{}
	settings := &network.EndpointSettings{}

	err := service.NetworkConnect(context.Background(), "test-network", "test-container", settings)

	if err != nil {
		t.Errorf("NetworkConnect() error = %v, want nil", err)
	}
}

func TestNetworkService_NetworkCreate(t *testing.T) {
	service := &NetworkService{}
	opts := network.CreateOptions{}

	tests := []struct {
		name         string
		networkName  string
		wantErr      bool
		wantErrType  error
		wantResponse bool
	}{
		{
			name:         "valid network",
			networkName:  "test-network",
			wantErr:      false,
			wantResponse: true,
		},
		{
			name:        "empty network name",
			networkName: "",
			wantErr:     true,
		},
		{
			name:        "notfound network",
			networkName: "notfound-network",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name:         "notfound network with ignore",
			networkName:  "notfound-ignorenotfound",
			wantErr:      false,
			wantResponse: true,
		},
		{
			name:        "not-found network",
			networkName: "not-found-network",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name:         "not-found network with ignore",
			networkName:  "not-found-ignore-not-found",
			wantErr:      false,
			wantResponse: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.NetworkCreate(context.Background(), tt.networkName, opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NetworkCreate() error = nil, wantErr %v", tt.wantErr)
				}
				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("NetworkCreate() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("NetworkCreate() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantResponse {
					if response.ID == "" {
						t.Errorf("NetworkCreate() response.ID = empty, want generated ID")
					}
				}
			}
		})
	}
}

func TestNetworkService_NetworkDisconnect(t *testing.T) {
	service := &NetworkService{}

	err := service.NetworkDisconnect(context.Background(), "test-network", "test-container", false)

	if err != nil {
		t.Errorf("NetworkDisconnect() error = %v, want nil", err)
	}
}

func TestNetworkService_NetworkInspect(t *testing.T) {
	service := &NetworkService{}
	opts := network.InspectOptions{}

	tests := []struct {
		name         string
		networkID    string
		wantErr      bool
		wantErrType  error
		wantResponse bool
	}{
		{
			name:         "valid network",
			networkID:    "test-network",
			wantErr:      false,
			wantResponse: true,
		},
		{
			name:      "empty network ID",
			networkID: "",
			wantErr:   true,
		},
		{
			name:        "notfound network",
			networkID:   "notfound-network",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
		{
			name:        "not-found network",
			networkID:   "not-found-network",
			wantErr:     true,
			wantErrType: errdefs.NotFound(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.NetworkInspect(context.Background(), tt.networkID, opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NetworkInspect() error = nil, wantErr %v", tt.wantErr)
				}
				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("NetworkInspect() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("NetworkInspect() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantResponse {
					if response.Name != tt.networkID {
						t.Errorf("NetworkInspect() response.Name = %v, want %v", response.Name, tt.networkID)
					}

					if response.Driver != "host" {
						t.Errorf("NetworkInspect() response.Driver = %v, want host", response.Driver)
					}

					if response.Scope != "local" {
						t.Errorf("NetworkInspect() response.Scope = %v, want local", response.Scope)
					}

					if response.ID == "" {
						t.Errorf("NetworkInspect() response.ID = empty, want generated ID")
					}
				}
			}
		})
	}
}

func TestNetworkService_NetworkInspectWithRaw(t *testing.T) {
	service := &NetworkService{}
	opts := network.InspectOptions{}

	tests := []struct {
		name         string
		networkID    string
		wantErr      bool
		wantResponse bool
		wantRaw      bool
	}{
		{
			name:         "valid network",
			networkID:    "test-network",
			wantErr:      false,
			wantResponse: true,
			wantRaw:      true,
		},
		{
			name:      "empty network ID",
			networkID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, raw, err := service.NetworkInspectWithRaw(context.Background(), tt.networkID, opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NetworkInspectWithRaw() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("NetworkInspectWithRaw() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantResponse {
					if response.Name != tt.networkID {
						t.Errorf("NetworkInspectWithRaw() response.Name = %v, want %v", response.Name, tt.networkID)
					}

					if response.Driver != "host" {
						t.Errorf("NetworkInspectWithRaw() response.Driver = %v, want host", response.Driver)
					}
				}

				if tt.wantRaw {
					if len(raw) == 0 {
						t.Errorf("NetworkInspectWithRaw() raw = empty, want data")
					}

					var unmarshaled network.Inspect
					if err := json.Unmarshal(raw, &unmarshaled); err != nil {
						t.Errorf("NetworkInspectWithRaw() raw data invalid JSON: %v", err)
					}
				}
			}
		})
	}
}

func TestNetworkService_NetworkList(t *testing.T) {
	service := &NetworkService{}
	opts := network.ListOptions{}

	networks, err := service.NetworkList(context.Background(), opts)

	if err != nil {
		t.Errorf("NetworkList() error = %v, want nil", err)
	}

	if networks != nil {
		t.Errorf("NetworkList() = %v, want nil", networks)
	}
}

func TestNetworkService_NetworkRemove(t *testing.T) {
	service := &NetworkService{}

	tests := []struct {
		name      string
		networkID string
		wantErr   bool
	}{
		{
			name:      "valid network",
			networkID: "test-network",
			wantErr:   false,
		},
		{
			name:      "empty network ID",
			networkID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.NetworkRemove(context.Background(), tt.networkID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NetworkRemove() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("NetworkRemove() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNetworkService_NetworksPrune(t *testing.T) {
	service := &NetworkService{}
	pruneFilters := filters.Args{}

	report, err := service.NetworksPrune(context.Background(), pruneFilters)

	if err != nil {
		t.Errorf("NetworksPrune() error = %v, want nil", err)
	}

	if report.NetworksDeleted != nil {
		t.Errorf("NetworksPrune() report.NetworksDeleted = %v, want nil", report.NetworksDeleted)
	}

	// SpaceReclaimed field may not exist in this version of Docker API
}

func TestNetworkService_InterfaceCompliance(t *testing.T) {
	var _ client.NetworkAPIClient = (*NetworkService)(nil)
}