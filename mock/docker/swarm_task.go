// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// TaskService implements all the task
// related functions for the Docker mock.
type TaskService struct{}

// TaskInspectWithRaw is a helper function to simulate
// a mocked call to inspect a task for a Docker swarm
// cluster and return the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.TaskInspectWithRaw
func (t *TaskService) TaskInspect(_ context.Context, _ string, _ client.TaskInspectOptions) (client.TaskInspectResult, error) {
	return client.TaskInspectResult{}, nil
}

// TaskList is a helper function to simulate
// a mocked call to list tasks for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.TaskList
func (t *TaskService) TaskList(_ context.Context, _ client.TaskListOptions) (client.TaskListResult, error) {
	return client.TaskListResult{}, nil
}

// TaskLogs is a helper function to simulate
// a mocked call to capture the logs from a
// task for a Docker swarm cluster.
func (t *TaskService) TaskLogs(_ context.Context, _ string, _ client.TaskLogsOptions) (client.TaskLogsResult, error) {
	return nil, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ServiceService satisfies the ServiceAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ServiceAPIClient
var _ client.TaskAPIClient = (*TaskService)(nil)
