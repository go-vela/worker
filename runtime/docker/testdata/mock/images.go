package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

// helper function to specify the routes for the Docker Image APIs
func imageRoutes(r *http.Request, path string) (*http.Response, error) {

	// get image id
	image := ""
	if strings.HasPrefix(path, "/images/") {
		image = strings.TrimPrefix(path, "/images/")
	}

	switch {
	case path == "/images/create":
		return createImage(r)
	case path == "/images/json":
		return listImages(r)
	case path == fmt.Sprintf("/images/%s", image):
		if r.Method == http.MethodDelete {
			return deleteImage(r, image)
		}
		return inspectImage(r, image)
	}

	return nil, nil
}

// helper function to return the mock results from an image create
func createImage(r *http.Request) (*http.Response, error) {

	name := r.URL.Query().Get("name")
	image := r.URL.Query().Get("fromImage")
	tag := r.URL.Query().Get("tag")

	b := []byte("test:latest")

	logrus.Infof("Creating a new image %s: %s:%s", name, image, tag)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an image inspect
func inspectImage(r *http.Request, image string) (*http.Response, error) {
	logrus.Infof("Docker mock inspecting image %s", image)

	if strings.Contains(image, "notfound") {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("image not found"))),
		}, nil
	}

	// $ docker pull alpine:latest
	//
	// $ docker image inspect alpine:latest
	resp := types.ImageInspect{
		ID:              "sha256:965ea09ff2ebd2b9eeec88cd822ce156f6674c7e99be082c7efac3c62f3ff652",
		RepoTags:        []string{"alpine:latest"},
		RepoDigests:     []string{"alpine@sha256:c19173c5ada610a5989151111163d28a67368362762534d8a8121ce95cf2bd5a"},
		Parent:          "",
		Comment:         "",
		Created:         "2019-10-21T17:21:42.387111039Z",
		Container:       "baae288169b1ae2f6bd82e7b605d8eb35a79e846385800e305eccc55b9bd5986",
		ContainerConfig: nil,
		DockerVersion:   "18.06.1-ce",
		Author:          "",
		Config:          nil,
		Architecture:    "amd64",
		Os:              "linux",
		Size:            5552690,
		VirtualSize:     5552690,
		GraphDriver: types.GraphDriverData{
			Data: map[string]string{
				"MergedDir": "/var/lib/docker/overlay2/9c00b01ddf812433e804c72a139eadc69ae79396207e127737de5c917b9e89b2/merged",
				"UpperDir":  "/var/lib/docker/overlay2/9c00b01ddf812433e804c72a139eadc69ae79396207e127737de5c917b9e89b2/diff",
				"WorkDir":   "/var/lib/docker/overlay2/9c00b01ddf812433e804c72a139eadc69ae79396207e127737de5c917b9e89b2/work",
			},
			Name: "overlay2",
		},
		RootFS: types.RootFS{
			Type:   "layers",
			Layers: []string{"sha256:77cae8ab23bf486355d1b3191259705374f4a11d483b24964d2f729dd8c076a0"},
		},
		Metadata: types.ImageMetadata{LastTagTime: time.Now()},
	}

	b, _ := json.Marshal(resp)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an image list
func listImages(r *http.Request) (*http.Response, error) {

	logrus.Info("Listing all images")
	b, _ := json.Marshal([]types.ImageSummary{
		{ID: "image_id_1"},
		{ID: "image_id_2"},
	})
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

// helper function to return the mock results from an image delete
func deleteImage(r *http.Request, image string) (*http.Response, error) {

	logrus.Infof("Deleting image %s", image)
	b, _ := json.Marshal([]types.ImageDeleteResponseItem{
		{Untagged: image},
		{Deleted: image},
	})
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}
