// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ServiceService implements all the service
// related functions for the Docker mock.
type ServiceService struct{}

// ServiceCreate is a helper function to simulate
// a mocked call to create a service for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceCreate
func (s *ServiceService) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options swarm.ServiceCreateOptions) (swarm.ServiceCreateResponse, error) {
	return swarm.ServiceCreateResponse{}, nil
}

// ServiceInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker service and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceInspectWithRaw
func (s *ServiceService) ServiceInspectWithRaw(ctx context.Context, serviceID string, options swarm.ServiceInspectOptions) (swarm.Service, []byte, error) {
	return swarm.Service{}, nil, nil
}

// ServiceList is a helper function to simulate
// a mocked call to list services for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceList
func (s *ServiceService) ServiceList(ctx context.Context, options swarm.ServiceListOptions) ([]swarm.Service, error) {
	return nil, nil
}

// ServiceLogs is a helper function to simulate
// a mocked call to capture the logs from a
// service for a Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceLogs
func (s *ServiceService) ServiceLogs(ctx context.Context, serviceID string, options container.LogsOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ServiceRemove is a helper function to simulate
// a mocked call to remove a service for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceRemove
func (s *ServiceService) ServiceRemove(ctx context.Context, serviceID string) error {
	return nil
}

// ServiceUpdate is a helper function to simulate
// a mocked call to update a service for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServiceUpdate
func (s *ServiceService) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options swarm.ServiceUpdateOptions) (swarm.ServiceUpdateResponse, error) {
	return swarm.ServiceUpdateResponse{}, nil
}

// TaskInspectWithRaw is a helper function to simulate
// a mocked call to inspect a task for a Docker swarm
// cluster and return the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.TaskInspectWithRaw
func (s *ServiceService) TaskInspectWithRaw(ctx context.Context, taskID string) (swarm.Task, []byte, error) {
	return swarm.Task{}, nil, nil
}

// TaskList is a helper function to simulate
// a mocked call to list tasks for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.TaskList
func (s *ServiceService) TaskList(ctx context.Context, options swarm.TaskListOptions) ([]swarm.Task, error) {
	return nil, nil
}

// TaskLogs is a helper function to simulate
// a mocked call to capture the logs from a
// task for a Docker swarm cluster.
func (s *ServiceService) TaskLogs(ctx context.Context, taskID string, options container.LogsOptions) (io.ReadCloser, error) {
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
