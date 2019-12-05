package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/sirupsen/logrus"
)

// helper function to specify the routes for the Docker Network APIs
func networkRoutes(r *http.Request, path string) (*http.Response, error) {

	// get network id
	net := ""
	if strings.HasPrefix(path, "/networks/") {
		net = strings.TrimPrefix(path, "/networks/")
	}

	switch {
	case strings.EqualFold(r.Method, http.MethodPost): // Path: /networks/create
		return createNetwork(r)
	case strings.EqualFold(r.Method, http.MethodGet): // Path: /networks/:network_id
		return inspectNetwork(r, net)
	case strings.EqualFold(r.Method, http.MethodDelete): // Path: /networks/:network_id
		// id := strings.Split(p, "/")[1]
		return removeNetwork(r, net)
	}

	return nil, nil
}

// helper function to return the mock results from a network creates
func createNetwork(r *http.Request) (*http.Response, error) {

	logrus.Infof("Creating a new network")
	b, _ := json.Marshal(types.NetworkCreateResponse{
		ID:      "network_id",
		Warning: "warning",
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from a network inspect
func inspectNetwork(r *http.Request, net string) (*http.Response, error) {
	logrus.Infof("Inspecting network %s", net)

	if strings.Contains(net, "notfound") {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("network not found"))),
		}, nil
	}

	resp := types.NetworkResource{
		Name:       "host",
		ID:         "d0097728e3575854f5d7e5704304aa0c3afefcdfbf0f037d19d48afa2f1cabeb",
		Created:    time.Now(),
		Scope:      "local",
		Driver:     "host",
		EnableIPv6: false,
		IPAM: network.IPAM{
			Driver:  "default",
			Options: map[string]string{},
			Config:  []network.IPAMConfig{},
		},
		Internal:   false,
		Attachable: false,
		Ingress:    false,
		ConfigFrom: network.ConfigReference{
			Network: "",
		},
		ConfigOnly: false,
		Containers: map[string]types.EndpointResource{},
		Options:    map[string]string{},
		Labels:     map[string]string{},
	}

	b, _ := json.Marshal(resp)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from a network remove
func removeNetwork(r *http.Request, net string) (*http.Response, error) {

	logrus.Infof("Deleting network %s", net)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}
