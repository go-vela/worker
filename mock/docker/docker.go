// SPDX-License-Identifier: Apache-2.0

package docker

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
// * the Docker version of v26.0 has a maximum API version of v1.45
// * to maintain n-1, the API version is pinned to v1.44
// .
const Version = "v1.44"

// New returns a client that is capable of handling
// Docker client calls and returning stub responses.
//
//nolint:revive // ignore unexported type as it is intentional
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
