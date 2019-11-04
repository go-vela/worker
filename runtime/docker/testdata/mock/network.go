package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

// helper function to specify the routes for the Docker Network APIs
func networkRoutes(r *http.Request, p string) (*http.Response, error) {

	switch {
	case strings.EqualFold(r.Method, http.MethodPost): // Path: /networks/create
		return createNetwork(r)
	case strings.EqualFold(r.Method, http.MethodDelete): // Path: /networks/:network_id
		id := strings.Split(p, "/")[1]
		return removeNetwork(r, id)
	}

	return nil, nil
}

// helper function to return the mock results from an network creates
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

// helper function to return the mock results from an network remove
func removeNetwork(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Deleting a network %s", id)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}
