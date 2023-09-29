// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

// SystemService implements all the system
// related functions for the Docker mock.
type SystemService struct{}

// DiskUsage is a helper function to simulate
// a mocked call to capture the data usage
// from the Docker daemon.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.DiskUsage
func (s *SystemService) DiskUsage(ctx context.Context) (types.DiskUsage, error) {
	return types.DiskUsage{}, nil
}

// Events is a helper function to simulate
// a mocked call to capture the events
// from the Docker daemon.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.Events
func (s *SystemService) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	return nil, nil
}

// Info is a helper function to simulate
// a mocked call to capture the system
// information from the Docker daemon.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.Info
func (s *SystemService) Info(ctx context.Context) (types.Info, error) {
	return types.Info{}, nil
}

// Ping is a helper function to simulate
// a mocked call to ping the Docker
// daemon and return version information.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.Ping
func (s *SystemService) Ping(ctx context.Context) (types.Ping, error) {
	return types.Ping{}, nil
}

// RegistryLogin is a helper function to simulate
// a mocked call to authenticate the Docker
// daemon against a Docker registry.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.RegistryLogin
func (s *SystemService) RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error) {
	return registry.AuthenticateOKBody{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// SystemService satisfies the SystemAPIClient interface that
// the Docker client expects.
//
// hhttps://pkg.go.dev/github.com/docker/docker/client?tab=doc#SystemAPIClient
var _ client.NetworkAPIClient = (*NetworkService)(nil)
