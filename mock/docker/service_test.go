// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func TestServiceService_ServiceCreate(t *testing.T) {
	service := &ServiceService{}
	spec := swarm.ServiceSpec{}
	opts := swarm.ServiceCreateOptions{}

	response, err := service.ServiceCreate(context.Background(), spec, opts)

	if err != nil {
		t.Errorf("ServiceCreate() error = %v, want nil", err)
	}

	if response.ID != "" {
		t.Errorf("ServiceCreate() response.ID = %v, want empty", response.ID)
	}

	if len(response.Warnings) != 0 {
		t.Errorf("ServiceCreate() response.Warnings = %v, want empty", response.Warnings)
	}
}

func TestServiceService_ServiceInspectWithRaw(t *testing.T) {
	service := &ServiceService{}
	opts := swarm.ServiceInspectOptions{}

	svc, raw, err := service.ServiceInspectWithRaw(context.Background(), "test-service", opts)

	if err != nil {
		t.Errorf("ServiceInspectWithRaw() error = %v, want nil", err)
	}

	if svc.ID != "" {
		t.Errorf("ServiceInspectWithRaw() service.ID = %v, want empty", svc.ID)
	}

	if raw != nil {
		t.Errorf("ServiceInspectWithRaw() raw = %v, want nil", raw)
	}
}

func TestServiceService_ServiceList(t *testing.T) {
	service := &ServiceService{}
	opts := swarm.ServiceListOptions{}

	services, err := service.ServiceList(context.Background(), opts)

	if err != nil {
		t.Errorf("ServiceList() error = %v, want nil", err)
	}

	if services != nil {
		t.Errorf("ServiceList() = %v, want nil", services)
	}
}

func TestServiceService_ServiceLogs(t *testing.T) {
	service := &ServiceService{}
	opts := container.LogsOptions{}

	logs, err := service.ServiceLogs(context.Background(), "test-service", opts)

	if err != nil {
		t.Errorf("ServiceLogs() error = %v, want nil", err)
	}

	if logs != nil {
		t.Errorf("ServiceLogs() = %v, want nil", logs)
	}
}

func TestServiceService_ServiceRemove(t *testing.T) {
	service := &ServiceService{}

	err := service.ServiceRemove(context.Background(), "test-service")

	if err != nil {
		t.Errorf("ServiceRemove() error = %v, want nil", err)
	}
}

func TestServiceService_ServiceUpdate(t *testing.T) {
	service := &ServiceService{}
	version := swarm.Version{}
	spec := swarm.ServiceSpec{}
	opts := swarm.ServiceUpdateOptions{}

	response, err := service.ServiceUpdate(context.Background(), "test-service", version, spec, opts)

	if err != nil {
		t.Errorf("ServiceUpdate() error = %v, want nil", err)
	}

	if len(response.Warnings) != 0 {
		t.Errorf("ServiceUpdate() response.Warnings = %v, want empty", response.Warnings)
	}
}

func TestServiceService_TaskInspectWithRaw(t *testing.T) {
	service := &ServiceService{}

	task, raw, err := service.TaskInspectWithRaw(context.Background(), "test-task")

	if err != nil {
		t.Errorf("TaskInspectWithRaw() error = %v, want nil", err)
	}

	if task.ID != "" {
		t.Errorf("TaskInspectWithRaw() task.ID = %v, want empty", task.ID)
	}

	if raw != nil {
		t.Errorf("TaskInspectWithRaw() raw = %v, want nil", raw)
	}
}

func TestServiceService_TaskList(t *testing.T) {
	service := &ServiceService{}
	opts := swarm.TaskListOptions{}

	tasks, err := service.TaskList(context.Background(), opts)

	if err != nil {
		t.Errorf("TaskList() error = %v, want nil", err)
	}

	if tasks != nil {
		t.Errorf("TaskList() = %v, want nil", tasks)
	}
}

func TestServiceService_TaskLogs(t *testing.T) {
	service := &ServiceService{}
	opts := container.LogsOptions{}

	logs, err := service.TaskLogs(context.Background(), "test-task", opts)

	if err != nil {
		t.Errorf("TaskLogs() error = %v, want nil", err)
	}

	if logs != nil {
		t.Errorf("TaskLogs() = %v, want nil", logs)
	}
}

func TestServiceService_InterfaceCompliance(t *testing.T) {
	var _ client.ServiceAPIClient = (*ServiceService)(nil)
}