// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/stringid"
)

// ContainerService implements all the container
// related functions for the Docker mock.
type ContainerService struct{}

// ContainerCreate is a helper function to simulate
// a mocked call to create a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerCreate
func (c *ContainerService) ContainerCreate(_ context.Context, opts client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	// verify a container was provided
	if len(opts.Name) == 0 {
		return client.ContainerCreateResult{},
			errors.New("no container provided")
	}

	// check if the container is not-found and
	// check if the not-found should be ignored
	if strings.Contains(opts.Name, "not-found") &&
		!strings.Contains(opts.Name, "ignore-not-found") {
		return client.ContainerCreateResult{}, cerrdefs.ErrNotFound
	}

	// check if the image is not found
	if strings.Contains(opts.Config.Image, "notfound") ||
		strings.Contains(opts.Config.Image, "not-found") {
		return client.ContainerCreateResult{}, cerrdefs.ErrNotFound
	}

	// create response object to return
	response := client.ContainerCreateResult{
		ID: stringid.GenerateRandomID(),
	}

	return response, nil
}

// ContainerInspect is a helper function to simulate
// a mocked call to inspect a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerInspect
func (c *ContainerService) ContainerInspect(_ context.Context, ctn string, _ client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return client.ContainerInspectResult{}, errors.New("no container provided")
	}

	// check if the container is notfound and
	// check if the notfound should be ignored
	if strings.Contains(ctn, "notfound") &&
		!strings.Contains(ctn, "ignorenotfound") {
		return client.ContainerInspectResult{}, cerrdefs.ErrNotFound
	}

	// check if the container is not-found and
	// check if the not-found should be ignored
	if strings.Contains(ctn, "not-found") &&
		!strings.Contains(ctn, "ignore-not-found") {
		return client.ContainerInspectResult{}, cerrdefs.ErrNotFound
	}

	// create response object to return
	response := client.ContainerInspectResult{
		Container: container.InspectResponse{
			ID:    stringid.GenerateRandomID(),
			Image: "alpine:latest",
			Name:  ctn,
			State: &container.State{Running: true},
		},
	}

	return response, nil
}

// ContainerList is a helper function to simulate
// a mocked call to list Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerList
func (c *ContainerService) ContainerList(_ context.Context, _ client.ContainerListOptions) (client.ContainerListResult, error) {
	return client.ContainerListResult{}, nil
}

// ContainerUpdate is a helper function to simulate
// a mocked call to update a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerUpdate
func (c *ContainerService) ContainerUpdate(_ context.Context, _ string, _ client.ContainerUpdateOptions) (client.ContainerUpdateResult, error) {
	return client.ContainerUpdateResult{}, nil
}

// ContainerRemove is a helper function to simulate
// a mocked call to remove a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerRemove
func (c *ContainerService) ContainerRemove(_ context.Context, ctn string, _ client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return client.ContainerRemoveResult{}, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		return client.ContainerRemoveResult{}, cerrdefs.ErrNotFound
	}

	return client.ContainerRemoveResult{}, nil
}

// ContainersPrune is a helper function to simulate
// a mocked call to prune Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainersPrune
func (c *ContainerService) ContainerPrune(_ context.Context, _ client.ContainerPruneOptions) (client.ContainerPruneResult, error) {
	return client.ContainerPruneResult{}, nil
}

// ContainerLogs is a helper function to simulate
// a mocked call to capture the logs from a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerLogs
func (c *ContainerService) ContainerLogs(_ context.Context, ctn string, _ client.ContainerLogsOptions) (client.ContainerLogsResult, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return nil, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		return nil, cerrdefs.ErrNotFound
	}

	var buf bytes.Buffer
	buf.WriteString("hello to stdout from github.com/go-vela/worker/mock/docker\n")
	buf.WriteString("hello to stderr from github.com/go-vela/worker/mock/docker\n")

	return io.NopCloser(&buf), nil
}

// ContainerStart is a helper function to simulate
// a mocked call to start a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStart
func (c *ContainerService) ContainerStart(_ context.Context, ctn string, _ client.ContainerStartOptions) (client.ContainerStartResult, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return client.ContainerStartResult{}, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		return client.ContainerStartResult{}, cerrdefs.ErrNotFound
	}

	return client.ContainerStartResult{}, nil
}

// ContainerStop is a helper function to simulate
// a mocked call to stop a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStop
func (c *ContainerService) ContainerStop(_ context.Context, ctn string, _ client.ContainerStopOptions) (client.ContainerStopResult, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return client.ContainerStopResult{}, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		return client.ContainerStopResult{}, cerrdefs.ErrNotFound
	}

	return client.ContainerStopResult{}, nil
}

// ContainerRestart is a helper function to simulate
// a mocked call to restart a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerRestart
func (c *ContainerService) ContainerRestart(_ context.Context, _ string, _ client.ContainerRestartOptions) (client.ContainerRestartResult, error) {
	return client.ContainerRestartResult{}, nil
}

// ContainerPause is a helper function to simulate
// a mocked call to pause a running Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerPause
func (c *ContainerService) ContainerPause(_ context.Context, _ string, _ client.ContainerPauseOptions) (client.ContainerPauseResult, error) {
	return client.ContainerPauseResult{}, nil
}

// ContainerUnpause is a helper function to simulate
// a mocked call to unpause a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerUnpause
func (c *ContainerService) ContainerUnpause(_ context.Context, _ string, _ client.ContainerUnpauseOptions) (client.ContainerUnpauseResult, error) {
	return client.ContainerUnpauseResult{}, nil
}

// ContainerWait is a helper function to simulate
// a mocked call to wait for a running Docker
// container to finish.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerWait
func (c *ContainerService) ContainerWait(_ context.Context, ctn string, _ client.ContainerWaitOptions) client.ContainerWaitResult {
	ctnCh := make(chan container.WaitResponse, 1)
	errCh := make(chan error, 1)

	result := client.ContainerWaitResult{
		Result: ctnCh,
		Error:  errCh,
	}

	// verify a container was provided
	if len(ctn) == 0 {
		// propagate the error to the error channel
		errCh <- errors.New("no container provided")

		return result
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") || strings.Contains(ctn, "not-found") {
		// propagate the error to the error channel
		errCh <- cerrdefs.ErrNotFound

		return result
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

	return result
}

// ContainerKill is a helper function to simulate
// a mocked call to kill a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerKill
func (c *ContainerService) ContainerKill(_ context.Context, ctn string, _ client.ContainerKillOptions) (client.ContainerKillResult, error) {
	// verify a container was provided
	if len(ctn) == 0 {
		return client.ContainerKillResult{}, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctn, "notfound") ||
		strings.Contains(ctn, "not-found") {
		return client.ContainerKillResult{}, cerrdefs.ErrNotFound
	}

	return client.ContainerKillResult{}, nil
}

// ContainerRename is a helper function to simulate
// a mocked call to rename a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerRename
func (c *ContainerService) ContainerRename(_ context.Context, _ string, _ client.ContainerRenameOptions) (client.ContainerRenameResult, error) {
	return client.ContainerRenameResult{}, nil
}

// ContainerResize is a helper function to simulate
// a mocked call to resize a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerResize
func (c *ContainerService) ContainerResize(_ context.Context, _ string, _ client.ContainerResizeOptions) (client.ContainerResizeResult, error) {
	return client.ContainerResizeResult{}, nil
}

// ContainerAttach is a helper function to simulate
// a mocked call to attach a connection to a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerAttach
func (c *ContainerService) ContainerAttach(_ context.Context, _ string, _ client.ContainerAttachOptions) (client.ContainerAttachResult, error) {
	return client.ContainerAttachResult{}, nil
}

// ContainerCommit is a helper function to simulate
// a mocked call to apply changes to a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerCommit
func (c *ContainerService) ContainerCommit(_ context.Context, _ string, _ client.ContainerCommitOptions) (client.ContainerCommitResult, error) {
	return client.ContainerCommitResult{}, nil
}

// ContainerDiff is a helper function to simulate
// a mocked call to show the differences in the
// filesystem between two Docker containers.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerDiff
func (c *ContainerService) ContainerDiff(_ context.Context, _ string, _ client.ContainerDiffOptions) (client.ContainerDiffResult, error) {
	return client.ContainerDiffResult{}, nil
}

// ContainerExport is a helper function to simulate
// a mocked call to expore the contents of a Docker
// container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerExport
func (c *ContainerService) ContainerExport(_ context.Context, _ string, _ client.ContainerExportOptions) (client.ContainerExportResult, error) {
	return nil, nil
}

// ContainerStats is a helper function to simulate
// a mocked call to capture information about a
// Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStats
func (c *ContainerService) ContainerStats(_ context.Context, _ string, _ client.ContainerStatsOptions) (client.ContainerStatsResult, error) {
	return client.ContainerStatsResult{}, nil
}

// ContainerTop is a helper function to simulate
// a mocked call to show running processes inside
// a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerTop
func (c *ContainerService) ContainerTop(_ context.Context, _ string, _ client.ContainerTopOptions) (client.ContainerTopResult, error) {
	return client.ContainerTopResult{}, nil
}

// ContainerStatPath is a helper function to simulate
// a mocked call to capture information about a path
// inside a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerStatPath
func (c *ContainerService) ContainerStatPath(_ context.Context, _ string, _ client.ContainerStatPathOptions) (client.ContainerStatPathResult, error) {
	return client.ContainerStatPathResult{}, nil
}

// CopyFromContainer is a helper function to simulate
// a mocked call to copy content from a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.CopyFromContainer
func (c *ContainerService) CopyFromContainer(_ context.Context, _ string, _ client.CopyFromContainerOptions) (client.CopyFromContainerResult, error) {
	return client.CopyFromContainerResult{}, nil
}

// CopyToContainer is a helper function to simulate
// a mocked call to copy content to a Docker container.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.CopyToContainer
func (c *ContainerService) CopyToContainer(_ context.Context, _ string, _ client.CopyToContainerOptions) (client.CopyToContainerResult, error) {
	return client.CopyToContainerResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ContainerService satisfies the ContainerAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ContainerAPIClient
var _ client.ContainerAPIClient = (*ContainerService)(nil)
