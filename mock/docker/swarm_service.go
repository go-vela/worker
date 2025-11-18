// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// ServiceService implements all the service
// related functions for the Docker mock.
type ServiceService struct{}

// ServiceCreate is a helper function to simulate
// a mocked call to create a service for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceCreate
func (s *ServiceService) ServiceCreate(_ context.Context, _ client.ServiceCreateOptions) (client.ServiceCreateResult, error) {
	return client.ServiceCreateResult{}, nil
}

// ServiceInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker service and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceInspectWithRaw
func (s *ServiceService) ServiceInspect(_ context.Context, _ string, _ client.ServiceInspectOptions) (client.ServiceInspectResult, error) {
	return client.ServiceInspectResult{}, nil
}

// ServiceList is a helper function to simulate
// a mocked call to list services for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceList
func (s *ServiceService) ServiceList(_ context.Context, _ client.ServiceListOptions) (client.ServiceListResult, error) {
	return client.ServiceListResult{}, nil
}

// ServiceUpdate is a helper function to simulate
// a mocked call to update a service for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceUpdate
func (s *ServiceService) ServiceUpdate(_ context.Context, _ string, _ client.ServiceUpdateOptions) (client.ServiceUpdateResult, error) {
	return client.ServiceUpdateResult{}, nil
}

// ServiceRemove is a helper function to simulate
// a mocked call to remove a service for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceRemove
func (s *ServiceService) ServiceRemove(_ context.Context, _ string, _ client.ServiceRemoveOptions) (client.ServiceRemoveResult, error) {
	return client.ServiceRemoveResult{}, nil
}

// ServiceLogs is a helper function to simulate
// a mocked call to capture the logs from a
// service for a Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceLogs
func (s *ServiceService) ServiceLogs(_ context.Context, _ string, _ client.ServiceLogsOptions) (client.ServiceLogsResult, error) {
	return nil, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ServiceService satisfies the ServiceAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ServiceAPIClient
var _ client.ServiceAPIClient = (*ServiceService)(nil)
