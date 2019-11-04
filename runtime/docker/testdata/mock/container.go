package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stringid"
	"github.com/sirupsen/logrus"
)

// helper function to return the mock results from an container creates
func createContainer(r *http.Request) (*http.Response, error) {

	logrus.Infof("Creating a new container: %s", r.URL.Query().Get("name"))
	b, _ := json.Marshal(container.ContainerCreateCreatedBody{
		ID: stringid.GenerateRandomID(),
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an listing containers
func listContainers(r *http.Request) (*http.Response, error) {

	logrus.Info("Listing all containers")
	b, _ := json.Marshal([]types.Container{
		{
			ID:      stringid.GenerateRandomID(),
			Names:   []string{"hello docker"},
			Image:   "test:image",
			ImageID: stringid.GenerateRandomID(),
			Command: "top",
		},
		{
			ID:      stringid.GenerateRandomID(),
			Names:   []string{"hello docker 2"},
			Image:   "test:image",
			ImageID: stringid.GenerateRandomID(),
			Command: "top",
		},
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an inspecting running containers
func getContainer(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Getting container with ID: %s", id)
	b, _ := json.Marshal(types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    id,
			Image: "test:image",
			Name:  "name",
			State: &types.ContainerState{
				Running: true,
			},
			HostConfig: &container.HostConfig{
				Resources: container.Resources{
					CPUQuota:  9999,
					CPUPeriod: 9999,
					CPUShares: 999,
					Memory:    99999,
				},
			},
		},
		Config: &container.Config{
			Labels: map[string]string{
				"ERU": "1",
			},
			Image: "test:image",
		},
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an removing a running containers
func removeContainer(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Removing container with ID: %s", id)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}

// helper function to return the mock results from starting a running containers
func startContainer(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Starting container with ID: %s", id)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}

// helper function to return the mock results from stopping a running containers
func stopContainer(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Stopping container with ID: %s", id)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}

// helper function to return the mock results from stopping a running containers
func waitContainer(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Waiting container with ID: %s", id)
	b, _ := json.Marshal(container.ContainerWaitOKBody{
		StatusCode: 15,
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}

// helper function to return the mock results from stopping a running containers
func killContainer(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Killing container with ID: %s", id)
	b, _ := json.Marshal(container.ContainerWaitOKBody{
		StatusCode: 15,
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}

// helper function to return the mock results from logs on a running containers
func logsContainer(r *http.Request, id string) (*http.Response, error) {

	b := []byte("Hello, Docker")

	logrus.Infof("Getting logs from container with ID: %s", id)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}
