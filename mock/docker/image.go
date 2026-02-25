// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/storage"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stringid"
)

// ImageService implements all the image
// related functions for the Docker mock.
type ImageService struct{}

// BuildCachePrune is a helper function to simulate
// a mocked call to prune the build cache for the
// Docker daemon.
func (i *ImageService) BuildCachePrune(_ context.Context, _ build.CachePruneOptions) (*build.CachePruneReport, error) {
	return nil, nil
}

// BuildCancel is a helper function to simulate
// a mocked call to cancel building a Docker image.
func (i *ImageService) BuildCancel(_ context.Context, _ string) error {
	return nil
}

// ImageBuild is a helper function to simulate
// a mocked call to build a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageBuild
func (i *ImageService) ImageBuild(_ context.Context, _ io.Reader, _ build.ImageBuildOptions) (build.ImageBuildResponse, error) {
	return build.ImageBuildResponse{}, nil
}

// ImageCreate is a helper function to simulate
// a mocked call to create a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageCreate
func (i *ImageService) ImageCreate(_ context.Context, _ string, _ image.CreateOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageHistory is a helper function to simulate
// a mocked call to inspect the history for a
// Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageHistory
func (i *ImageService) ImageHistory(_ context.Context, _ string, _ ...client.ImageHistoryOption) ([]image.HistoryResponseItem, error) {
	return nil, nil
}

// ImageImport is a helper function to simulate
// a mocked call to import a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageImport
func (i *ImageService) ImageImport(_ context.Context, _ image.ImportSource, _ string, _ image.ImportOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageInspect is a helper function to simulate
// a mocked call to inspect a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageInspect
func (i *ImageService) ImageInspectWithRaw(_ context.Context, _ string) (image.InspectResponse, []byte, error) {
	return image.InspectResponse{}, nil, nil
}

// ImageInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker image and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageInspectWithRaw
func (i *ImageService) ImageInspect(_ context.Context, img string, _ ...client.ImageInspectOption) (image.InspectResponse, error) {
	// verify an image was provided
	if len(img) == 0 {
		return image.InspectResponse{}, errors.New("no image provided")
	}

	// check if the image is not found
	if strings.Contains(img, "notfound") || strings.Contains(img, "not-found") {
		return image.InspectResponse{},
			errdefs.NotFound(
				//nolint:staticcheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", img),
			)
	}

	path := fmt.Sprintf("/var/lib/docker/overlay2/%s", stringid.GenerateRandomID())

	// create response object to return
	response := image.InspectResponse{
		ID:            fmt.Sprintf("sha256:%s", stringid.GenerateRandomID()),
		RepoTags:      []string{"alpine:latest"},
		RepoDigests:   []string{fmt.Sprintf("alpine@sha256:%s", stringid.GenerateRandomID())},
		Created:       time.Now().String(),
		Container:     fmt.Sprintf("sha256:%s", stringid.GenerateRandomID()),
		DockerVersion: "19.03.1-ce",
		Architecture:  "amd64",
		Os:            "linux",
		Size:          5552690,
		VirtualSize:   5552690,
		GraphDriver: storage.DriverData{
			Data: map[string]string{
				"MergedDir": fmt.Sprintf("%s/merged", path),
				"UpperDir":  fmt.Sprintf("%s/diff", path),
				"WorkDir":   fmt.Sprintf("%s/work", path),
			},
			Name: "overlay2",
		},
		RootFS: image.RootFS{
			Type:   "layers",
			Layers: []string{fmt.Sprintf("sha256:%s", stringid.GenerateRandomID())},
		},
		Metadata: image.Metadata{LastTagTime: time.Now()},
	}

	return response, nil
}

// ImageList is a helper function to simulate
// a mocked call to list Docker images.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageList
func (i *ImageService) ImageList(_ context.Context, _ image.ListOptions) ([]image.Summary, error) {
	return nil, nil
}

// ImageLoad is a helper function to simulate
// a mocked call to load a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageLoad
func (i *ImageService) ImageLoad(_ context.Context, _ io.Reader, _ ...client.ImageLoadOption) (image.LoadResponse, error) {
	return image.LoadResponse{}, nil
}

// ImagePull is a helper function to simulate
// a mocked call to pull a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagePull
func (i *ImageService) ImagePull(_ context.Context, image string, _ image.PullOptions) (io.ReadCloser, error) {
	// verify an image was provided
	if len(image) == 0 {
		return nil, errors.New("no container provided")
	}

	// check if the image is notfound and
	// check if the notfound should be ignored
	if strings.Contains(image, "notfound") &&
		!strings.Contains(image, "ignorenotfound") {
		return nil,
			errdefs.NotFound(
				//nolint:staticcheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", image),
			)
	}

	// check if the image is not-found and
	// check if the not-found should be ignored
	if strings.Contains(image, "not-found") &&
		!strings.Contains(image, "ignore-not-found") {
		return nil,
			errdefs.NotFound(
				//nolint:staticcheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", image),
			)
	}

	// create response object to return
	response := io.NopCloser(
		bytes.NewReader(

			fmt.Appendf(nil, "%s\n%s\n%s\n%s\n",
				fmt.Sprintf("latest: Pulling from %s", image),
				fmt.Sprintf("Digest: sha256:%s", stringid.GenerateRandomID()),
				fmt.Sprintf("Status: Image is up to date for %s", image),
				image,
			),
		),
	)

	return response, nil
}

// ImagePush is a helper function to simulate
// a mocked call to push a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagePush
func (i *ImageService) ImagePush(_ context.Context, _ string, _ image.PushOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageRemove is a helper function to simulate
// a mocked call to remove a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageRemove
func (i *ImageService) ImageRemove(_ context.Context, _ string, _ image.RemoveOptions) ([]image.DeleteResponse, error) {
	return nil, nil
}

// ImageSave is a helper function to simulate
// a mocked call to save a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageSave
func (i *ImageService) ImageSave(_ context.Context, _ []string, _ ...client.ImageSaveOption) (io.ReadCloser, error) {
	return nil, nil
}

// ImageSearch is a helper function to simulate
// a mocked call to search for a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageSearch
func (i *ImageService) ImageSearch(_ context.Context, _ string, _ registry.SearchOptions) ([]registry.SearchResult, error) {
	return nil, nil
}

// ImageTag is a helper function to simulate
// a mocked call to tag a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageTag
func (i *ImageService) ImageTag(_ context.Context, _ string, _ string) error {
	return nil
}

// ImagesPrune is a helper function to simulate
// a mocked call to prune Docker images.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagesPrune
func (i *ImageService) ImagesPrune(_ context.Context, _ filters.Args) (image.PruneReport, error) {
	return image.PruneReport{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ImageService satisfies the ImageAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ImageAPIClient
var _ client.ImageAPIClient = (*ImageService)(nil)
