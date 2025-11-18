// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// SystemService implements all the system
// related functions for the Docker mock.
type SystemService struct{}

// Events is a helper function to simulate
// a mocked call to capture the events
// from the Docker daemon.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.Events
func (s *SystemService) Events(_ context.Context, _ client.EventsListOptions) client.EventsResult {
	return client.EventsResult{}
}

// Info is a helper function to simulate
// a mocked call to capture the system
// information from the Docker daemon.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.Info
func (s *SystemService) Info(_ context.Context, _ client.InfoOptions) (client.SystemInfoResult, error) {
	return client.SystemInfoResult{}, nil
}

// RegistryLogin is a helper function to simulate
// a mocked call to authenticate the Docker
// daemon against a Docker registry.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.RegistryLogin
func (s *SystemService) RegistryLogin(_ context.Context, _ client.RegistryLoginOptions) (client.RegistryLoginResult, error) {
	return client.RegistryLoginResult{}, nil
}

// DiskUsage is a helper function to simulate
// a mocked call to capture the data usage
// from the Docker daemon.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.DiskUsage
func (s *SystemService) DiskUsage(_ context.Context, _ client.DiskUsageOptions) (client.DiskUsageResult, error) {
	return client.DiskUsageResult{}, nil
}

// Ping is a helper function to simulate
// a mocked call to ping the Docker
// daemon and return version information.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.Ping
func (s *SystemService) Ping(_ context.Context, _ client.PingOptions) (client.PingResult, error) {
	return client.PingResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// SystemService satisfies the SystemAPIClient interface that
// the Docker client expects.
//
// hhttps://pkg.go.dev/github.com/docker/docker/client#SystemAPIClient
var _ client.NetworkAPIClient = (*NetworkService)(nil)
