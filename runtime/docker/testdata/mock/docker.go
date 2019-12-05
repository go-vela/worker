// Package mock is a collection of functions from moby/moby
// to mock hitting the Docker API
// https://github.com/moby/moby/tree/master/client
package mock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
)

const mockAPIVersion = "v1.38"

var (
	mockContextWithTimeout bool
	err                    error
	ctx                    = context.Background()
)

// Client to create the mock client used for http requests.
// This function exists in moby/moby: https://github.com/moby/moby/blob/master/client/client_mock_test.go#L20
func Client(doer func(*http.Request) (*http.Response, error)) *http.Client {

	return &http.Client{
		Transport: transportFunc(doer),
	}
}

// Router to emulate the Docker router to return mock results
// for Docker API requests
func Router(r *http.Request) (*http.Response, error) {
	prefix := fmt.Sprintf("/%s", mockAPIVersion)
	path := strings.TrimPrefix(r.URL.Path, prefix)

	// get container id
	containerID := ""
	if strings.HasPrefix(path, "/containers/") {
		cid := strings.TrimPrefix(path, "/containers/")
		containerID = strings.Split(cid, "/")[0]
		if containerID == "" {
			containerID = "_"
		}
	}

	switch {

	// Image endpoints
	case strings.HasPrefix(path, "/images/"):
		return imageRoutes(r, path)

	// Container endpoints
	case path == "/containers/create":
		return createContainer(r)
	case path == "/containers/json":
		return listContainers(r)
	case path == fmt.Sprintf("/containers/%s/json", containerID):
		return getContainer(r, containerID)
	case path == fmt.Sprintf("/containers/%s", containerID):
		return removeContainer(r, containerID)
	case path == fmt.Sprintf("/containers/%s/start", containerID):
		return startContainer(r, containerID)
	case path == fmt.Sprintf("/containers/%s/stop", containerID):
		return stopContainer(r, containerID)
	case path == fmt.Sprintf("/containers/%s/wait", containerID):
		return waitContainer(r, containerID)
	case path == fmt.Sprintf("/containers/%s/kill", containerID):
		return killContainer(r, containerID)
	case path == fmt.Sprintf("/containers/%s/logs", containerID):
		return logsContainer(r, containerID)

	// Network endpoints
	case strings.HasPrefix(path, "/networks/"):
		return networkRoutes(r, path)

	// Volume endpoints
	case strings.HasPrefix(path, "/volumes/"):
		return volumeRoutes(r, path)
	}

	return errorMock(500, fmt.Sprintf("Server Error, unknown path: %s", path))
}

type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

// helper function to return the error mocks in Docker. This function exists in moby/moby
// https://github.com/moby/moby/blob/master/client/client_mock_test.go#L26
func errorMock(statusCode int, message string) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	body, err := json.Marshal(&types.ErrorResponse{
		Message: message,
	})
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
		Header:     header,
	}, fmt.Errorf(message)
}
