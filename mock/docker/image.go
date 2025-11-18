// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"strings"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/jsonstream"
	"github.com/moby/moby/api/types/storage"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/stringid"
)

// ImageService implements all the image
// related functions for the Docker mock.
type ImageService struct{}

// ImageImport is a helper function to simulate
// a mocked call to import a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageImport
func (i *ImageService) ImageImport(_ context.Context, _ client.ImageImportSource, _ string, _ client.ImageImportOptions) (client.ImageImportResult, error) {
	return nil, nil
}

// ImageList is a helper function to simulate
// a mocked call to list Docker images.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageList
func (i *ImageService) ImageList(_ context.Context, _ client.ImageListOptions) (client.ImageListResult, error) {
	return client.ImageListResult{}, nil
}

type mockImagePullResponse struct {
	r io.ReadCloser
}

func (m *mockImagePullResponse) Read(p []byte) (int, error) {
	return m.r.Read(p)
}

func (m *mockImagePullResponse) Close() error {
	return m.r.Close()
}

func (m *mockImagePullResponse) JSONMessages(_ context.Context) iter.Seq2[jsonstream.Message, error] {
	return nil
}

func (m *mockImagePullResponse) Wait(_ context.Context) error {
	return nil
}

// ImagePull is a helper function to simulate
// a mocked call to pull a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagePull
func (i *ImageService) ImagePull(_ context.Context, image string, _ client.ImagePullOptions) (client.ImagePullResponse, error) {
	// verify an image was provided
	if len(image) == 0 {
		return nil, errors.New("no container provided")
	}

	// check if the image is notfound and
	// check if the notfound should be ignored
	if strings.Contains(image, "notfound") &&
		!strings.Contains(image, "ignorenotfound") {
		return nil, cerrdefs.ErrNotFound
	}

	// check if the image is not-found and
	// check if the not-found should be ignored
	if strings.Contains(image, "not-found") &&
		!strings.Contains(image, "ignore-not-found") {
		return nil, cerrdefs.ErrNotFound
	}

	payload := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		fmt.Sprintf("latest: Pulling from %s", image),
		fmt.Sprintf("Digest: sha256:%s", stringid.GenerateRandomID()),
		fmt.Sprintf("Status: Image is up to date for %s", image),
		image,
	)

	// simple non-streaming variant: all data available immediately
	reader := io.NopCloser(bytes.NewReader([]byte(payload)))

	return &mockImagePullResponse{r: reader}, nil
}

// ImagePush is a helper function to simulate
// a mocked call to push a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagePush
func (i *ImageService) ImagePush(_ context.Context, _ string, _ client.ImagePushOptions) (client.ImagePushResponse, error) {
	return nil, nil
}

// ImageRemove is a helper function to simulate
// a mocked call to remove a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageRemove
func (i *ImageService) ImageRemove(_ context.Context, _ string, _ client.ImageRemoveOptions) (client.ImageRemoveResult, error) {
	return client.ImageRemoveResult{}, nil
}

// ImageTag is a helper function to simulate
// a mocked call to tag a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageTag
func (i *ImageService) ImageTag(_ context.Context, _ client.ImageTagOptions) (client.ImageTagResult, error) {
	return client.ImageTagResult{}, nil
}

// ImagesPrune is a helper function to simulate
// a mocked call to prune Docker images.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagesPrune
func (i *ImageService) ImagePrune(_ context.Context, _ client.ImagePruneOptions) (client.ImagePruneResult, error) {
	return client.ImagePruneResult{}, nil
}

// ImageInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker image and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageInspectWithRaw
func (i *ImageService) ImageInspect(_ context.Context, img string, _ ...client.ImageInspectOption) (client.ImageInspectResult, error) {
	// verify an image was provided
	if len(img) == 0 {
		return client.ImageInspectResult{}, errors.New("no image provided")
	}

	// check if the image is not found
	if strings.Contains(img, "notfound") || strings.Contains(img, "not-found") {
		return client.ImageInspectResult{}, cerrdefs.ErrNotFound
	}

	path := fmt.Sprintf("/var/lib/docker/overlay2/%s", stringid.GenerateRandomID())

	// create response object to return
	response := image.InspectResponse{
		ID:           fmt.Sprintf("sha256:%s", stringid.GenerateRandomID()),
		RepoTags:     []string{"alpine:latest"},
		RepoDigests:  []string{fmt.Sprintf("alpine@sha256:%s", stringid.GenerateRandomID())},
		Created:      time.Now().String(),
		Architecture: "amd64",
		Os:           "linux",
		Size:         5552690,
		GraphDriver: &storage.DriverData{
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

	return client.ImageInspectResult{InspectResponse: response}, nil
}

// ImageHistory is a helper function to simulate
// a mocked call to inspect the history for a
// Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageHistory
func (i *ImageService) ImageHistory(_ context.Context, _ string, _ ...client.ImageHistoryOption) (client.ImageHistoryResult, error) {
	return client.ImageHistoryResult{}, nil
}

// ImageLoad is a helper function to simulate
// a mocked call to load a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageLoad
func (i *ImageService) ImageLoad(_ context.Context, _ io.Reader, _ ...client.ImageLoadOption) (client.ImageLoadResult, error) {
	return nil, nil
}

// ImageSave is a helper function to simulate
// a mocked call to save a Docker image.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageSave
func (i *ImageService) ImageSave(_ context.Context, _ []string, _ ...client.ImageSaveOption) (client.ImageSaveResult, error) {
	return nil, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// ImageService satisfies the ImageAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#ImageAPIClient
var _ client.ImageAPIClient = (*ImageService)(nil)
