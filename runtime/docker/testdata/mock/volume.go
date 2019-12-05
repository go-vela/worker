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
func volumeRoutes(r *http.Request, path string) (*http.Response, error) {

	// get volume id
	vol := ""
	if strings.HasPrefix(path, "/volumes/") {
		vol = strings.TrimPrefix(path, "/volumes/")
	}

	switch {
	case strings.EqualFold(r.Method, http.MethodPost): // Path: /volumes/create
		return createVolume(r)
	case strings.EqualFold(r.Method, http.MethodGet): // Path: /volumes/:volume_id
		return inspectVolume(r, vol)
	case strings.EqualFold(r.Method, http.MethodDelete): // Path: /volumes/:volume_id
		return removeVolume(r, vol)
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

// helper function to return the mock results from a network inspect
func inspectVolume(r *http.Request, vol string) (*http.Response, error) {
	logrus.Infof("Inspecting volume %s", vol)

	if strings.Contains(vol, "notfound") {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("volume not found"))),
		}, nil
	}

	resp := types.Volume{
		CreatedAt:  "2019-12-05T12:00:00Z",
		Driver:     "local",
		Labels:     map[string]string{},
		Mountpoint: "/var/lib/docker/volumes/9c00b01ddf812433e804c72a139eadc69ae79396207e127737de5c917b9e89b2/_data",
		Name:       "9c00b01ddf812433e804c72a139eadc69ae79396207e127737de5c917b9e89b2",
		Options:    map[string]string{},
		Scope:      "local",
	}

	b, _ := json.Marshal(resp)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an volume remove
func removeVolume(r *http.Request, vol string) (*http.Response, error) {

	logrus.Infof("Deleting a volume %s", vol)
	var b []byte

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)), // empty body
	}, nil
}
