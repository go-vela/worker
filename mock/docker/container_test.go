// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
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

func TestContainerCreate(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		config   *container.Config
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			config:  &container.Config{Image: "alpine:latest"},
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful creation",
			config:  &container.Config{Image: "alpine:latest"},
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			config:  &container.Config{Image: "alpine:latest"},
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			config:  &container.Config{Image: "alpine:latest"},
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "ignorenotfound",
			config:  &container.Config{Image: "alpine:latest"},
			ctnName: "ignorenotfound",
			wantErr: false,
		},
		{
			name:    "ignore-not-found",
			config:  &container.Config{Image: "alpine:latest"},
			ctnName: "ignore-not-found",
			wantErr: false,
		},
		{
			name:    "image notfound",
			config:  &container.Config{Image: "notfound:latest"},
			ctnName: "test-container",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "not found")
			},
		},
		{
			name:    "image not-found",
			config:  &container.Config{Image: "not-found:latest"},
			ctnName: "test-container",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ContainerCreate(
				context.Background(),
				tt.config,
				&container.HostConfig{},
				&network.NetworkingConfig{},
				&v1.Platform{},
				tt.ctnName,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerCreate() error check failed for error: %v", err)
			}

			if !tt.wantErr && resp.ID == "" {
				t.Error("ContainerCreate() returned empty ID")
			}
		})
	}
}

func TestContainerInspect(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful inspect",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "ignorenotfound",
			ctnName: "ignorenotfound",
			wantErr: false,
		},
		{
			name:    "ignore-not-found",
			ctnName: "ignore-not-found",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ContainerInspect(context.Background(), tt.ctnName)

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerInspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerInspect() error check failed for error: %v", err)
			}

			if !tt.wantErr {
				if resp.ID == "" {
					t.Error("ContainerInspect() returned empty ID")
				}

				if resp.Name != tt.ctnName {
					t.Errorf("ContainerInspect() Name = %v, want %v", resp.Name, tt.ctnName)
				}
			}
		})
	}
}

func TestContainerInspectWithRaw(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful inspect",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, raw, err := service.ContainerInspectWithRaw(context.Background(), tt.ctnName, false)

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerInspectWithRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerInspectWithRaw() error check failed for error: %v", err)
			}

			if !tt.wantErr {
				if resp.ID == "" {
					t.Error("ContainerInspectWithRaw() returned empty ID")
				}

				if len(raw) == 0 {
					t.Error("ContainerInspectWithRaw() returned empty raw bytes")
				}
			}
		})
	}
}

func TestContainerKill(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful kill",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ContainerKill(context.Background(), tt.ctnName, "SIGTERM")

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerKill() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerKill() error check failed for error: %v", err)
			}
		})
	}
}

func TestContainerLogs(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful logs",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := service.ContainerLogs(context.Background(), tt.ctnName, container.LogsOptions{})

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerLogs() error check failed for error: %v", err)
			}

			if !tt.wantErr && reader != nil {
				defer reader.Close()
				// Read logs to verify content
				logs, _ := io.ReadAll(reader)
				if len(logs) == 0 {
					t.Error("ContainerLogs() returned empty logs")
				}
				// Verify logs contain expected content
				logsStr := string(logs)
				if !strings.Contains(logsStr, "stdout") && !strings.Contains(logsStr, "stderr") {
					t.Error("ContainerLogs() logs don't contain expected output")
				}
			}
		})
	}
}

func TestContainerRemove(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful remove",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ContainerRemove(context.Background(), tt.ctnName, container.RemoveOptions{})

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerRemove() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerRemove() error check failed for error: %v", err)
			}
		})
	}
}

func TestContainerStart(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful start",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ContainerStart(context.Background(), tt.ctnName, container.StartOptions{})

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerStart() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerStart() error check failed for error: %v", err)
			}
		})
	}
}

func TestContainerStop(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful stop",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ContainerStop(context.Background(), tt.ctnName, container.StopOptions{})

			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerStop() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
				t.Errorf("ContainerStop() error check failed for error: %v", err)
			}
		})
	}
}

func TestContainerWait(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name     string
		ctnName  string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "empty container name",
			ctnName: "",
			wantErr: true,
		},
		{
			name:    "successful wait",
			ctnName: "test-container",
			wantErr: false,
		},
		{
			name:    "container notfound",
			ctnName: "notfound",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
		{
			name:    "container not-found",
			ctnName: "not-found",
			wantErr: true,
			errCheck: func(err error) bool {
				return strings.Contains(err.Error(), "No such container")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respCh, errCh := service.ContainerWait(context.Background(), tt.ctnName, container.WaitCondition(""))

			select {
			case resp := <-respCh:
				if tt.wantErr {
					t.Errorf("ContainerWait() expected error but got response: %v", resp)
				}

				if resp.StatusCode != 15 {
					t.Errorf("ContainerWait() StatusCode = %v, want 15", resp.StatusCode)
				}
			case err := <-errCh:
				if (err != nil) != tt.wantErr {
					t.Errorf("ContainerWait() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.errCheck != nil && err != nil && !tt.errCheck(err) {
					t.Errorf("ContainerWait() error check failed for error: %v", err)
				}
			case <-time.After(2 * time.Second):
				if tt.wantErr {
					t.Error("ContainerWait() timeout waiting for error")
				} else {
					t.Error("ContainerWait() timeout waiting for response")
				}
			}
		})
	}
}
