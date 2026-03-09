// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strings"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/stringid"
)

// ExecService implements all the exec
// related functions for the Docker mock.
type ExecService struct{}

// ExecCreate is a helper function to simulate
// a mocked call to create a Docker exec instance.
func (e *ExecService) ExecCreate(_ context.Context, ctnID string, opts client.ExecCreateOptions) (client.ExecCreateResult, error) {
	// verify a container was provided
	if len(ctnID) == 0 {
		return client.ExecCreateResult{}, errors.New("no container provided")
	}

	// check if the container is not found
	if strings.Contains(ctnID, "notfound") || strings.Contains(ctnID, "not-found") {
		return client.ExecCreateResult{}, cerrdefs.ErrNotFound
	}

	// create exec ID based on command for testing scenarios
	execID := stringid.GenerateRandomID()

	// check command for specific test scenarios
	if len(opts.Cmd) > 0 {
		cmdStr := strings.Join(opts.Cmd, " ")
		if strings.Contains(cmdStr, "error") {
			execID = "exec-error-" + execID
		} else if strings.Contains(cmdStr, "multiline") {
			execID = "exec-multiline-" + execID
		} else if strings.Contains(cmdStr, "find") {
			// For artifact file search commands
			if strings.Contains(cmdStr, "/not-found") {
				execID = "exec-not-found-" + execID
			} else if strings.Contains(cmdStr, "artifacts") {
				execID = "exec-artifacts-find-" + execID
			} else if strings.Contains(cmdStr, "test-results") && strings.Contains(cmdStr, ".xml") {
				execID = "exec-test-results-xml-" + execID
			} else if strings.Contains(cmdStr, "cypress/screenshots") && strings.Contains(cmdStr, ".png") {
				execID = "exec-cypress-screenshots-" + execID
			} else if strings.Contains(cmdStr, "cypress/videos") && strings.Contains(cmdStr, ".mp4") {
				execID = "exec-cypress-videos-" + execID
			} else if strings.Contains(cmdStr, "cypress") {
				// Generic cypress pattern for combined searches
				execID = "exec-cypress-all-" + execID
			}
		}
	}

	// create response object to return
	response := client.ExecCreateResult{
		ID: execID,
	}

	return response, nil
}

// ExecInspect is a helper function to simulate
// a mocked call to inspect a Docker exec instance.
func (e *ExecService) ExecInspect(_ context.Context, _ string, _ client.ExecInspectOptions) (client.ExecInspectResult, error) {
	return client.ExecInspectResult{}, nil
}

// ExecResize is a helper function to simulate
// a mocked call to resize a Docker exec instance.
func (e *ExecService) ExecResize(_ context.Context, _ string, _ client.ExecResizeOptions) (client.ExecResizeResult, error) {
	return client.ExecResizeResult{}, nil
}

// ExecStart is a helper function to simulate
// a mocked call to start a Docker exec instance.
func (e *ExecService) ExecStart(_ context.Context, _ string, _ client.ExecStartOptions) (client.ExecStartResult, error) {
	return client.ExecStartResult{}, nil
}

// ExecAttach is a helper function to simulate
// a mocked call to attach to a Docker exec instance.
func (e *ExecService) ExecAttach(_ context.Context, execID string, _ client.ExecAttachOptions) (client.ExecAttachResult, error) {
	// create a buffer to hold mock output
	var buf bytes.Buffer

	// check for specific test scenarios based on execID
	if strings.Contains(execID, "error") {
		writeStreamString(&buf, byte(stdcopy.Stderr), "mock exec error")
	} else if strings.Contains(execID, "multiline") {
		writeStreamString(&buf, byte(stdcopy.Stdout), "line1\nline2\nline3")
	} else if strings.Contains(execID, "artifacts-find") {
		writeStreamString(&buf, byte(stdcopy.Stdout), "/vela/artifacts/test_results/alpha.txt\n/vela/artifacts/test_results/beta.txt\n/vela/artifacts/build_results/alpha.txt\n/vela/artifacts/build_results/beta.txt")
	} else if strings.Contains(execID, "test-results-xml") {
		writeStreamString(&buf, byte(stdcopy.Stdout), "/vela/workspace/test-results/junit.xml\n/vela/workspace/test-results/report.xml")
	} else if strings.Contains(execID, "cypress-screenshots") {
		writeStreamString(&buf, byte(stdcopy.Stdout), "/vela/workspace/cypress/screenshots/test1/screenshot1.png\n/vela/workspace/cypress/screenshots/test1/screenshot2.png\n/vela/workspace/cypress/screenshots/test2/error.png")
	} else if strings.Contains(execID, "cypress-videos") {
		writeStreamString(&buf, byte(stdcopy.Stdout), "/vela/workspace/cypress/videos/test1.mp4\n/vela/workspace/cypress/videos/test2.mp4")
	} else if strings.Contains(execID, "cypress-all") {
		writeStreamString(&buf, byte(stdcopy.Stdout), "/vela/workspace/cypress/screenshots/test1/screenshot1.png\n/vela/workspace/cypress/screenshots/test2/error.png\n/vela/workspace/cypress/videos/test1.mp4\n/vela/workspace/cypress/videos/test2.mp4")
	} else if strings.Contains(execID, "not-found") {
		writeStreamString(&buf, byte(stdcopy.Stderr), "find: '/not-found': No such file or directory")
	} else {
		writeStreamString(&buf, byte(stdcopy.Stdout), "mock exec output")
	}

	// create a HijackedResponse with the mock data
	response := client.ExecAttachResult{
		HijackedResponse: client.HijackedResponse{
			Reader: bufio.NewReader(&buf),
			Conn:   &mockConn{}, // Use mock connection to avoid nil pointer dereference
		},
	}

	return response, nil
}

// writeStreamString writes a docker multiplexed stream frame to the buffer.
func writeStreamString(buf *bytes.Buffer, stream byte, payload string) {
	header := []byte{stream, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(header[4:], uint32(len(payload)))
	_, _ = buf.Write(header)
	_, _ = buf.WriteString(payload)
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ExecService satisfies the ExecAPIClient interface that
// the Docker client expects.
var _ client.ExecAPIClient = (*ExecService)(nil)
