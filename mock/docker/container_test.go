// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestContainerServiceSimple(t *testing.T) {
	service := &ContainerService{}

	// Test ContainerAttach
	_, err := service.ContainerAttach(context.Background(), "test", container.AttachOptions{})
	if err != nil {
		t.Errorf("ContainerAttach() error = %v, want nil", err)
	}

	// Test ContainerCommit
	_, err = service.ContainerCommit(context.Background(), "test", container.CommitOptions{})
	if err != nil {
		t.Errorf("ContainerCommit() error = %v, want nil", err)
	}

	// Test ContainerDiff
	_, err = service.ContainerDiff(context.Background(), "test")
	if err != nil {
		t.Errorf("ContainerDiff() error = %v, want nil", err)
	}

	// Test ContainerExecAttach
	_, err = service.ContainerExecAttach(context.Background(), "test", container.ExecAttachOptions{})
	if err != nil {
		t.Errorf("ContainerExecAttach() error = %v, want nil", err)
	}

	// Test ContainerExecCreate
	_, err = service.ContainerExecCreate(context.Background(), "test", container.ExecOptions{})
	if err != nil {
		t.Errorf("ContainerExecCreate() error = %v, want nil", err)
	}

	// Test ContainerExecInspect
	_, err = service.ContainerExecInspect(context.Background(), "test")
	if err != nil {
		t.Errorf("ContainerExecInspect() error = %v, want nil", err)
	}

	// Test ContainerExecResize
	err = service.ContainerExecResize(context.Background(), "test", container.ResizeOptions{})
	if err != nil {
		t.Errorf("ContainerExecResize() error = %v, want nil", err)
	}

	// Test ContainerExecStart
	err = service.ContainerExecStart(context.Background(), "test", container.ExecStartOptions{})
	if err != nil {
		t.Errorf("ContainerExecStart() error = %v, want nil", err)
	}

	// Test ContainerExport
	_, err = service.ContainerExport(context.Background(), "test")
	if err != nil {
		t.Errorf("ContainerExport() error = %v, want nil", err)
	}

	// Test ContainerList
	_, err = service.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		t.Errorf("ContainerList() error = %v, want nil", err)
	}

	// Test ContainerPause
	err = service.ContainerPause(context.Background(), "test")
	if err != nil {
		t.Errorf("ContainerPause() error = %v, want nil", err)
	}

	// Test ContainerRename
	err = service.ContainerRename(context.Background(), "old", "new")
	if err != nil {
		t.Errorf("ContainerRename() error = %v, want nil", err)
	}

	// Test ContainerResize
	err = service.ContainerResize(context.Background(), "test", container.ResizeOptions{})
	if err != nil {
		t.Errorf("ContainerResize() error = %v, want nil", err)
	}

	// Test ContainerRestart
	err = service.ContainerRestart(context.Background(), "test", container.StopOptions{})
	if err != nil {
		t.Errorf("ContainerRestart() error = %v, want nil", err)
	}

	// Test ContainerStatPath
	_, err = service.ContainerStatPath(context.Background(), "test", "/path")
	if err != nil {
		t.Errorf("ContainerStatPath() error = %v, want nil", err)
	}

	// Test ContainerStats
	_, err = service.ContainerStats(context.Background(), "test", true)
	if err != nil {
		t.Errorf("ContainerStats() error = %v, want nil", err)
	}

	// Test ContainerTop
	_, err = service.ContainerTop(context.Background(), "test", []string{"aux"})
	if err != nil {
		t.Errorf("ContainerTop() error = %v, want nil", err)
	}

	// Test ContainerUnpause
	err = service.ContainerUnpause(context.Background(), "test")
	if err != nil {
		t.Errorf("ContainerUnpause() error = %v, want nil", err)
	}

	// Test ContainerUpdate
	_, err = service.ContainerUpdate(context.Background(), "test", container.UpdateConfig{})
	if err != nil {
		t.Errorf("ContainerUpdate() error = %v, want nil", err)
	}

	// Test ContainersPrune
	_, err = service.ContainersPrune(context.Background(), filters.Args{})
	if err != nil {
		t.Errorf("ContainersPrune() error = %v, want nil", err)
	}

	// Test ContainerStatsOneShot
	_, err = service.ContainerStatsOneShot(context.Background(), "test")
	if err != nil {
		t.Errorf("ContainerStatsOneShot() error = %v, want nil", err)
	}

	// Test CopyFromContainer
	_, _, err = service.CopyFromContainer(context.Background(), "test", "/path")
	if err != nil {
		t.Errorf("CopyFromContainer() error = %v, want nil", err)
	}

	// Test CopyToContainer
	err = service.CopyToContainer(context.Background(), "test", "/path", nil, container.CopyToContainerOptions{})
	if err != nil {
		t.Errorf("CopyToContainer() error = %v, want nil", err)
	}

	// Test interface compliance
	var _ client.ContainerAPIClient = (*ContainerService)(nil)
}