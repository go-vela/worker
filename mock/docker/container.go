// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/stringid"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// ContainerService implements all the container
// related functions for the Docker mock.
type ContainerService struct{}

// ContainerAttach is a helper function to simulate
// a mocked call to attach a connection to a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerAttach
func (c *ContainerService) ContainerAttach(ctx context.Context, ctn string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerCommit is a helper function to simulate
// a mocked call to apply changes to a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerCommit
func (c *ContainerService) ContainerCommit(ctx context.Context, ctn string, options types.ContainerCommitOptions) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

// ContainerCreate is a helper function to simulate
// a mocked call to create a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerCreate
func (c *ContainerService) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, p *v1.Platform, ctn string) (container.ContainerCreateCreatedBody, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return container.ContainerCreateCreatedBody{},
			errors.New("no container provided")
	}

	// check if the container is notfound and
	// check if the notfound should be ignored
	if strings.Contains(ctn, "notfound") &&
		!strings.Contains(ctn, "ignorenotfound") {
		return container.ContainerCreateCreatedBody{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// check if the container is not-found and
	// check if the not-found should be ignored
	if strings.Contains(ctn, "not-found") &&
		!strings.Contains(ctn, "ignore-not-found") {
		return container.ContainerCreateCreatedBody{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// check if the image is not found
	if strings.Contains(config.Image, "notfound") ||
		strings.Contains(config.Image, "not-found") {
		return container.ContainerCreateCreatedBody{},
			errdefs.NotFound(
				//nolint:stylecheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", config.Image),
			)
	}

	// create response object to return
	response := container.ContainerCreateCreatedBody{
		ID: stringid.GenerateRandomID(),
	}

	return response, nil
}

// ContainerDiff is a helper function to simulate
// a mocked call to show the differences in the
// filesystem between two Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerDiff
func (c *ContainerService) ContainerDiff(ctx context.Context, ctn string) ([]container.ContainerChangeResponseItem, error) {
	return nil, nil
}

// ContainerExecAttach is a helper function to simulate
// a mocked call to attach a connection to a process
// running inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerExecAttach
func (c *ContainerService) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerExecCreate is a helper function to simulate
// a mocked call to create a process to run inside a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerExecCreate
func (c *ContainerService) ContainerExecCreate(ctx context.Context, ctn string, config types.ExecConfig) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

// ContainerExecInspect is a helper function to simulate
// a mocked call to inspect a process running inside a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerExecInspect
func (c *ContainerService) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	return types.ContainerExecInspect{}, nil
}

// ContainerExecResize is a helper function to simulate
// a mocked call to resize the tty for a process running
// inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerExecResize
func (c *ContainerService) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	return nil
}

// ContainerExecStart is a helper function to simulate
// a mocked call to start a process inside a Docker
// container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerExecStart
func (c *ContainerService) ContainerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	return nil
}

// ContainerExport is a helper function to simulate
// a mocked call to expore the contents of a Docker
// container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerExport
func (c *ContainerService) ContainerExport(ctx context.Context, ctn string) (io.ReadCloser, error) {
	return nil, nil
}

// ContainerInspect is a helper function to simulate
// a mocked call to inspect a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerInspect
func (c *ContainerService) ContainerInspect(ctx context.Context, ctn string) (types.ContainerJSON, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return types.ContainerJSON{}, errors.New("no container provided")
	}

	// check if the container is notfound and
	// check if the notfound should be ignored
	if strings.Contains(ctn, "notfound") &&
		!strings.Contains(ctn, "ignorenotfound") {
		return types.ContainerJSON{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// check if the container is not-found and
	// check if the not-found should be ignored
	if strings.Contains(ctn, "not-found") &&
		!strings.Contains(ctn, "ignore-not-found") {
		return types.ContainerJSON{},
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// create response object to return
	response := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    stringid.GenerateRandomID(),
			Image: "alpine:latest",
			Name:  ctn,
			State: &types.ContainerState{Running: true},
		},
		Config: &container.Config{
			Image: "alpine:latest",
		},
	}

	return response, nil
}

// ContainerInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker container and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerInspectWithRaw
func (c *ContainerService) ContainerInspectWithRaw(ctx context.Context, ctn string, getSize bool) (types.ContainerJSON, []byte, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return types.ContainerJSON{}, nil, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		return types.ContainerJSON{},
			nil,
			//nolint:stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// create response object to return
	response := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    stringid.GenerateRandomID(),
			Image: "alpine:latest",
			Name:  ctn,
			State: &types.ContainerState{Running: true},
		},
		Config: &container.Config{
			Image: "alpine:latest",
		},
	}

	// marshal response into raw bytes
	b, err := json.Marshal(response)
	if err != nil {
		return types.ContainerJSON{}, nil, err
	}

	return response, b, nil
}

// ContainerKill is a helper function to simulate
// a mocked call to kill a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerKill
func (c *ContainerService) ContainerKill(ctx context.Context, ctn, signal string) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		//nolint:stylecheck // messsage is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerList is a helper function to simulate
// a mocked call to list Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerList
func (c *ContainerService) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return nil, nil
}

// ContainerLogs is a helper function to simulate
// a mocked call to capture the logs from a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerLogs
func (c *ContainerService) ContainerLogs(ctx context.Context, ctn string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return nil, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		//nolint:stylecheck // messsage is capitalized to match Docker messages
		return nil, errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// create response object to return
	response := new(bytes.Buffer)

	// write stdout logs to response buffer
	_, err := stdcopy.
		NewStdWriter(response, stdcopy.Stdout).
		Write([]byte("hello to stdout from github.com/go-vela/worker/mock/docker"))
	if err != nil {
		return nil, err
	}

	// write stderr logs to response buffer
	_, err = stdcopy.
		NewStdWriter(response, stdcopy.Stderr).
		Write([]byte("hello to stderr from github.com/go-vela/worker/mock/docker"))
	if err != nil {
		return nil, err
	}

	return io.NopCloser(response), nil
}

// ContainerPause is a helper function to simulate
// a mocked call to pause a running Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerPause
func (c *ContainerService) ContainerPause(ctx context.Context, ctn string) error {
	return nil
}

// ContainerRemove is a helper function to simulate
// a mocked call to remove a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerRemove
func (c *ContainerService) ContainerRemove(ctx context.Context, ctn string, options types.ContainerRemoveOptions) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		//nolint:stylecheck // messsage is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerRename is a helper function to simulate
// a mocked call to rename a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerRename
func (c *ContainerService) ContainerRename(ctx context.Context, container, newContainerName string) error {
	return nil
}

// ContainerResize is a helper function to simulate
// a mocked call to resize a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerResize
func (c *ContainerService) ContainerResize(ctx context.Context, ctn string, options types.ResizeOptions) error {
	return nil
}

// ContainerRestart is a helper function to simulate
// a mocked call to restart a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerRestart
func (c *ContainerService) ContainerRestart(ctx context.Context, ctn string, timeout *time.Duration) error {
	return nil
}

// ContainerStart is a helper function to simulate
// a mocked call to start a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerStart
func (c *ContainerService) ContainerStart(ctx context.Context, ctn string, options types.ContainerStartOptions) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		//nolint:stylecheck // messsage is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerStatPath is a helper function to simulate
// a mocked call to capture information about a path
// inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerStatPath
func (c *ContainerService) ContainerStatPath(ctx context.Context, container, path string) (types.ContainerPathStat, error) {
	return types.ContainerPathStat{}, nil
}

// ContainerStats is a helper function to simulate
// a mocked call to capture information about a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerStats
func (c *ContainerService) ContainerStats(ctx context.Context, ctn string, stream bool) (types.ContainerStats, error) {
	return types.ContainerStats{}, nil
}

// ContainerStop is a helper function to simulate
// a mocked call to stop a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerStop
func (c *ContainerService) ContainerStop(ctx context.Context, ctn string, timeout *time.Duration) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		//nolint:stylecheck // messsage is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerTop is a helper function to simulate
// a mocked call to show running processes inside
// a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerTop
func (c *ContainerService) ContainerTop(ctx context.Context, ctn string, arguments []string) (container.ContainerTopOKBody, error) {
	return container.ContainerTopOKBody{}, nil
}

// ContainerUnpause is a helper function to simulate
// a mocked call to unpause a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerUnpause
func (c *ContainerService) ContainerUnpause(ctx context.Context, ctn string) error {
	return nil
}

// ContainerUpdate is a helper function to simulate
// a mocked call to update a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerUpdate
func (c *ContainerService) ContainerUpdate(ctx context.Context, ctn string, updateConfig container.UpdateConfig) (container.ContainerUpdateOKBody, error) {
	return container.ContainerUpdateOKBody{}, nil
}

// ContainerWait is a helper function to simulate
// a mocked call to wait for a running Docker
// container to finish.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainerWait
func (c *ContainerService) ContainerWait(ctx context.Context, ctn string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	ctnCh := make(chan container.ContainerWaitOKBody, 1)
	errCh := make(chan error, 1)

	// verify a container was provided
	if len(ctn) == 0 {
		// propagate the error to the error channel
		errCh <- errors.New("no container provided")

		return ctnCh, errCh
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		// propagate the error to the error channel
		//nolint:stylecheck // messsage is capitalized to match Docker messages
		errCh <- errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))

		return ctnCh, errCh
	}

	// create goroutine for responding to call
	go func() {
		// create response object to return
		response := container.ContainerWaitOKBody{
			StatusCode: 15,
		}

		// sleep for 1 second to simulate waiting for the container
		time.Sleep(1 * time.Second)

		// propagate the response to the container channel
		ctnCh <- response

		// propagate nil to the error channel
		errCh <- nil
	}()

	return ctnCh, errCh
}

// ContainersPrune is a helper function to simulate
// a mocked call to prune Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ContainersPrune
func (c *ContainerService) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (types.ContainersPruneReport, error) {
	return types.ContainersPruneReport{}, nil
}

// ContainerStatsOneShot is a helper function to simulate
// a mocked call to return near realtime stats for a given container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.CopyFromContainer
func (c *ContainerService) ContainerStatsOneShot(ctx context.Context, containerID string) (types.ContainerStats, error) {
	return types.ContainerStats{}, nil
}

// CopyFromContainer is a helper function to simulate
// a mocked call to copy content from a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.CopyFromContainer
func (c *ContainerService) CopyFromContainer(ctx context.Context, container, srcPath string) (io.ReadCloser, types.ContainerPathStat, error) {
	return nil, types.ContainerPathStat{}, nil
}

// CopyToContainer is a helper function to simulate
// a mocked call to copy content to a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.CopyToContainer
func (c *ContainerService) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ContainerService satisfies the ContainerAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#ContainerAPIClient
var _ client.ContainerAPIClient = (*ContainerService)(nil)
