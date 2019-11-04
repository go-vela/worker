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

// helper function to specify the routes for the Docker Volume APIs
func volumeRoutes(r *http.Request, p string) (*http.Response, error) {

	switch {
	case strings.EqualFold(r.Method, http.MethodPost): // Path: /Volumes/create
		return createVolume(r)
	case strings.EqualFold(r.Method, http.MethodDelete): // Path: /Volumes/:Volume_id
		id := strings.Split(p, "/")[1]
		return removeVolume(r, id)
	}

	return nil, nil
}

// helper function to return the mock results from an volume creates
func createVolume(r *http.Request) (*http.Response, error) {

	logrus.Infof("Creating a new volume")
	b, _ := json.Marshal(types.Volume{
		Name:       "volume",
		Driver:     "local",
		Mountpoint: "mountpoint",
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an volume remove
func removeVolume(r *http.Request, id string) (*http.Response, error) {

	logrus.Infof("Deleting a volume %s", id)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}
