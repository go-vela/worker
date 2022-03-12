// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

// nolint: godot // ignore comment ending in a list
//
// Version represents the supported Docker API version for the mock.
//
// The Docker API version is pinned to ensure compatibility between the
// Docker API and client. The goal is to maintain n-1 compatibility.
//
// The maximum supported Docker API version for the client is here:
//
// https://docs.docker.com/engine/api/#api-version-matrix
//
// For example (use the compatibility matrix above for reference):
//
// * the Docker version of v20.10 has a maximum API version of v1.41
// * to maintain n-1, the API version is pinned to v1.40
const Version = "v1.40"

// New returns a client that is capable of handling
// Docker client calls and returning stub responses.
func New() (*mock, error) {
	return &mock{
		ConfigService:       ConfigService{},
		ContainerService:    ContainerService{},
		DistributionService: DistributionService{},
		ImageService:        ImageService{},
		NetworkService:      NetworkService{},
		NodeService:         NodeService{},
		PluginService:       PluginService{},
		SecretService:       SecretService{},
		ServiceService:      ServiceService{},
		SystemService:       SystemService{},
		SwarmService:        SwarmService{},
		VolumeService:       VolumeService{},
		Version:             Version,
	}, nil
}
