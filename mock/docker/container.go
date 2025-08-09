// SPDX-License-Identifier: Apache-2.0

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
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerAttach
func (c *ContainerService) ContainerAttach(ctx context.Context, ctn string, options container.AttachOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerCommit is a helper function to simulate
// a mocked call to apply changes to a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerCommit
func (c *ContainerService) ContainerCommit(ctx context.Context, ctn string, options container.CommitOptions) (container.CommitResponse, error) {
	return container.CommitResponse{}, nil
}

// ContainerCreate is a helper function to simulate
// a mocked call to create a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerCreate
func (c *ContainerService) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, p *v1.Platform, ctn string) (container.CreateResponse, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return container.CreateResponse{},
			errors.New("no container provided")
	}

	// check if the container is notfound and
	// check if the notfound should be ignored
	if strings.Contains(ctn, "notfound") &&
		!strings.Contains(ctn, "ignorenotfound") {
		return container.CreateResponse{},
			//nolint:staticcheck // message is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// check if the container is not-found and
	// check if the not-found should be ignored
	if strings.Contains(ctn, "not-found") &&
		!strings.Contains(ctn, "ignore-not-found") {
		return container.CreateResponse{},
			//nolint:staticcheck // message is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// check if the image is not found
	if strings.Contains(config.Image, "notfound") ||
		strings.Contains(config.Image, "not-found") {
		return container.CreateResponse{},
			errdefs.NotFound(
				//nolint:staticcheck // message is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", config.Image),
			)
	}

	// create response object to return
	response := container.CreateResponse{
		ID: stringid.GenerateRandomID(),
	}

	return response, nil
}

// ContainerDiff is a helper function to simulate
// a mocked call to show the differences in the
// filesystem between two Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerDiff
func (c *ContainerService) ContainerDiff(ctx context.Context, ctn string) ([]container.FilesystemChange, error) {
	return nil, nil
}

// ContainerExecAttach is a helper function to simulate
// a mocked call to attach a connection to a process
// running inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExecAttach
func (c *ContainerService) ContainerExecAttach(ctx context.Context, execID string, config container.ExecAttachOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerExecCreate is a helper function to simulate
// a mocked call to create a process to run inside a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExecCreate
func (c *ContainerService) ContainerExecCreate(ctx context.Context, ctn string, config container.ExecOptions) (container.ExecCreateResponse, error) {
	return container.ExecCreateResponse{}, nil
}

// ContainerExecInspect is a helper function to simulate
// a mocked call to inspect a process running inside a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExecInspect
func (c *ContainerService) ContainerExecInspect(ctx context.Context, execID string) (container.ExecInspect, error) {
	return container.ExecInspect{}, nil
}

// ContainerExecResize is a helper function to simulate
// a mocked call to resize the tty for a process running
// inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExecResize
func (c *ContainerService) ContainerExecResize(ctx context.Context, execID string, options container.ResizeOptions) error {
	return nil
}

// ContainerExecStart is a helper function to simulate
// a mocked call to start a process inside a Docker
// container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExecStart
func (c *ContainerService) ContainerExecStart(ctx context.Context, execID string, config container.ExecStartOptions) error {
	return nil
}

// ContainerExport is a helper function to simulate
// a mocked call to expore the contents of a Docker
// container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExport
func (c *ContainerService) ContainerExport(ctx context.Context, ctn string) (io.ReadCloser, error) {
	return nil, nil
}

// ContainerInspect is a helper function to simulate
// a mocked call to inspect a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerInspect
func (c *ContainerService) ContainerInspect(ctx context.Context, ctn string) (container.InspectResponse, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return container.InspectResponse{}, errors.New("no container provided")
	}

	// check if the container is notfound and
	// check if the notfound should be ignored
	if strings.Contains(ctn, "notfound") &&
		!strings.Contains(ctn, "ignorenotfound") {
		return container.InspectResponse{},
			//nolint:staticcheck // message is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// check if the container is not-found and
	// check if the not-found should be ignored
	if strings.Contains(ctn, "not-found") &&
		!strings.Contains(ctn, "ignore-not-found") {
		return container.InspectResponse{},
			//nolint:staticcheck // message is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// create response object to return
	response := container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:    stringid.GenerateRandomID(),
			Image: "alpine:latest",
			Name:  ctn,
			State: &container.State{Running: true},
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
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerInspectWithRaw
func (c *ContainerService) ContainerInspectWithRaw(ctx context.Context, ctn string, getSize bool) (container.InspectResponse, []byte, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return container.InspectResponse{}, nil, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		return container.InspectResponse{},
			nil,
			//nolint:staticcheck // message is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	// create response object to return
	response := container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:    stringid.GenerateRandomID(),
			Image: "alpine:latest",
			Name:  ctn,
			State: &container.State{Running: true},
		},
		Config: &container.Config{
			Image: "alpine:latest",
		},
	}

	// marshal response into raw bytes
	b, err := json.Marshal(response)
	if err != nil {
		return container.InspectResponse{}, nil, err
	}

	return response, b, nil
}

// ContainerKill is a helper function to simulate
// a mocked call to kill a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerKill
func (c *ContainerService) ContainerKill(ctx context.Context, ctn, signal string) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		//nolint:staticcheck // message is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerList is a helper function to simulate
// a mocked call to list Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerList
func (c *ContainerService) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	return nil, nil
}

// ContainerLogs is a helper function to simulate
// a mocked call to capture the logs from a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerLogs
func (c *ContainerService) ContainerLogs(ctx context.Context, ctn string, options container.LogsOptions) (io.ReadCloser, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return nil, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		//nolint:staticcheck // message is capitalized to match Docker messages
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
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerPause
func (c *ContainerService) ContainerPause(ctx context.Context, ctn string) error {
	return nil
}

// ContainerRemove is a helper function to simulate
// a mocked call to remove a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerRemove
func (c *ContainerService) ContainerRemove(ctx context.Context, ctn string, options container.RemoveOptions) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		//nolint:staticcheck // message is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerRename is a helper function to simulate
// a mocked call to rename a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerRename
func (c *ContainerService) ContainerRename(ctx context.Context, container, newContainerName string) error {
	return nil
}

// ContainerResize is a helper function to simulate
// a mocked call to resize a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerResize
func (c *ContainerService) ContainerResize(ctx context.Context, ctn string, options container.ResizeOptions) error {
	return nil
}

// ContainerRestart is a helper function to simulate
// a mocked call to restart a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerRestart
func (c *ContainerService) ContainerRestart(ctx context.Context, ctn string, options container.StopOptions) error {
	return nil
}

// ContainerStart is a helper function to simulate
// a mocked call to start a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStart
func (c *ContainerService) ContainerStart(ctx context.Context, ctn string, options container.StartOptions) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		//nolint:staticcheck // message is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerStatPath is a helper function to simulate
// a mocked call to capture information about a path
// inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStatPath
func (c *ContainerService) ContainerStatPath(ctx context.Context, containerID, path string) (container.PathStat, error) {
	return container.PathStat{}, nil
}

// ContainerStats is a helper function to simulate
// a mocked call to capture information about a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStats
func (c *ContainerService) ContainerStats(ctx context.Context, ctn string, stream bool) (container.StatsResponseReader, error) {
	return container.StatsResponseReader{}, nil
}

// ContainerStop is a helper function to simulate
// a mocked call to stop a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStop
func (c *ContainerService) ContainerStop(ctx context.Context, ctn string, options container.StopOptions) error {
	// verify a container was provided
	if len(ctn) == 0 {
		return errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		//nolint:staticcheck // message is capitalized to match Docker messages
		return errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))
	}

	return nil
}

// ContainerTop is a helper function to simulate
// a mocked call to show running processes inside
// a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerTop
func (c *ContainerService) ContainerTop(ctx context.Context, ctn string, arguments []string) (container.TopResponse, error) {
	return container.TopResponse{}, nil
}

// ContainerUnpause is a helper function to simulate
// a mocked call to unpause a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerUnpause
func (c *ContainerService) ContainerUnpause(ctx context.Context, ctn string) error {
	return nil
}

// ContainerUpdate is a helper function to simulate
// a mocked call to update a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerUpdate
func (c *ContainerService) ContainerUpdate(ctx context.Context, ctn string, updateConfig container.UpdateConfig) (container.UpdateResponse, error) {
	return container.UpdateResponse{}, nil
}

// ContainerWait is a helper function to simulate
// a mocked call to wait for a running Docker
// container to finish.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerWait
func (c *ContainerService) ContainerWait(ctx context.Context, ctn string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	ctnCh := make(chan container.WaitResponse, 1)
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
		//nolint:staticcheck // message is capitalized to match Docker messages
		errCh <- errdefs.NotFound(fmt.Errorf("Error: No such container: %s", ctn))

		return ctnCh, errCh
	}

	// create goroutine for responding to call
	go func() {
		// create response object to return
		response := container.WaitResponse{
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
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainersPrune
func (c *ContainerService) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (container.PruneReport, error) {
	return container.PruneReport{}, nil
}

// ContainerStatsOneShot is a helper function to simulate
// a mocked call to return near realtime stats for a given container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.CoontainerStatsOneShot
func (c *ContainerService) ContainerStatsOneShot(ctx context.Context, containerID string) (container.StatsResponseReader, error) {
	return container.StatsResponseReader{}, nil
}

// CopyFromContainer is a helper function to simulate
// a mocked call to copy content from a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.CopyFromContainer
func (c *ContainerService) CopyFromContainer(ctx context.Context, containerID, srcPath string) (io.ReadCloser, container.PathStat, error) {
	return nil, container.PathStat{}, nil
}

// CopyToContainer is a helper function to simulate
// a mocked call to copy content to a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.CopyToContainer
func (c *ContainerService) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options container.CopyToContainerOptions) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ContainerService satisfies the ContainerAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ContainerAPIClient
var _ client.ContainerAPIClient = (*ContainerService)(nil)
