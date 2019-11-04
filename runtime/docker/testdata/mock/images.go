package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

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

	logrus.Infof("Deleting a image %s", image)
	b, _ := json.Marshal([]types.ImageDeleteResponseItem{
		{Untagged: image},
		{Deleted: image},
	})
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}
