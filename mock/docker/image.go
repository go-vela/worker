// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
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
func (i *ImageService) BuildCachePrune(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error) {
	return nil, nil
}

// BuildCancel is a helper function to simulate
// a mocked call to cancel building a Docker image.
func (i *ImageService) BuildCancel(ctx context.Context, id string) error {
	return nil
}

// ImageBuild is a helper function to simulate
// a mocked call to build a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageBuild
func (i *ImageService) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{}, nil
}

// ImageCreate is a helper function to simulate
// a mocked call to create a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageCreate
func (i *ImageService) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageHistory is a helper function to simulate
// a mocked call to inspect the history for a
// Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageHistory
func (i *ImageService) ImageHistory(ctx context.Context, image string) ([]image.HistoryResponseItem, error) {
	return nil, nil
}

// ImageImport is a helper function to simulate
// a mocked call to import a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageImport
func (i *ImageService) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker image and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageInspectWithRaw
func (i *ImageService) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	// verify an image was provided
	if len(image) == 0 {
		return types.ImageInspect{}, nil, errors.New("no image provided")
	}

	// check if the image is not found
	if strings.Contains(image, "notfound") || strings.Contains(image, "not-found") {
		return types.ImageInspect{},
			nil,
			errdefs.NotFound(
				// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", image),
			)
	}

	path := fmt.Sprintf("/var/lib/docker/overlay2/%s", stringid.GenerateRandomID())

	// create response object to return
	response := types.ImageInspect{
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
		GraphDriver: types.GraphDriverData{
			Data: map[string]string{
				"MergedDir": fmt.Sprintf("%s/merged", path),
				"UpperDir":  fmt.Sprintf("%s/diff", path),
				"WorkDir":   fmt.Sprintf("%s/work", path),
			},
			Name: "overlay2",
		},
		RootFS: types.RootFS{
			Type:   "layers",
			Layers: []string{fmt.Sprintf("sha256:%s", stringid.GenerateRandomID())},
		},
		Metadata: types.ImageMetadata{LastTagTime: time.Now()},
	}

	// marshal response into raw bytes
	b, err := json.Marshal(response)
	if err != nil {
		return types.ImageInspect{}, nil, err
	}

	return response, b, nil
}

// ImageList is a helper function to simulate
// a mocked call to list Docker images.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageList
func (i *ImageService) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	return nil, nil
}

// ImageLoad is a helper function to simulate
// a mocked call to load a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageLoad
func (i *ImageService) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	return types.ImageLoadResponse{}, nil
}

// ImagePull is a helper function to simulate
// a mocked call to pull a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImagePull
func (i *ImageService) ImagePull(ctx context.Context, image string, options types.ImagePullOptions) (io.ReadCloser, error) {
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
				// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", image),
			)
	}

	// check if the image is not-found and
	// check if the not-found should be ignored
	if strings.Contains(image, "not-found") &&
		!strings.Contains(image, "ignore-not-found") {
		return nil,
			errdefs.NotFound(
				// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
				fmt.Errorf("Error response from daemon: manifest for %s not found: manifest unknown", image),
			)
	}

	// create response object to return
	response := ioutil.NopCloser(
		bytes.NewReader(
			[]byte(
				fmt.Sprintf("%s\n%s\n%s\n%s\n",
					fmt.Sprintf("latest: Pulling from %s", image),
					fmt.Sprintf("Digest: sha256:%s", stringid.GenerateRandomID()),
					fmt.Sprintf("Status: Image is up to date for %s", image),
					image,
				),
			),
		),
	)

	return response, nil
}

// ImagePush is a helper function to simulate
// a mocked call to push a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImagePush
func (i *ImageService) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageRemove is a helper function to simulate
// a mocked call to remove a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageRemove
func (i *ImageService) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	return nil, nil
}

// ImageSave is a helper function to simulate
// a mocked call to save a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageSave
func (i *ImageService) ImageSave(ctx context.Context, images []string) (io.ReadCloser, error) {
	return nil, nil
}

// ImageSearch is a helper function to simulate
// a mocked call to search for a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageSearch
func (i *ImageService) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	return nil, nil
}

// ImageTag is a helper function to simulate
// a mocked call to tag a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImageTag
func (i *ImageService) ImageTag(ctx context.Context, image, ref string) error {
	return nil
}

// ImagesPrune is a helper function to simulate
// a mocked call to prune Docker images.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ImagesPrune
func (i *ImageService) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	return types.ImagesPruneReport{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ImageService satisfies the ImageAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#ImageAPIClient
var _ client.ImageAPIClient = (*ImageService)(nil)
