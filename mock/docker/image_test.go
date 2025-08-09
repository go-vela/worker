// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func TestImageService_BuildCachePrune(t *testing.T) {
	service := &ImageService{}
	opts := build.CachePruneOptions{}

	report, err := service.BuildCachePrune(context.Background(), opts)
	if err != nil {
		t.Errorf("BuildCachePrune() error = %v, want nil", err)
	}

	if report != nil {
		t.Errorf("BuildCachePrune() = %v, want nil", report)
	}
}

func TestImageService_BuildCancel(t *testing.T) {
	service := &ImageService{}

	err := service.BuildCancel(context.Background(), "test-id")
	if err != nil {
		t.Errorf("BuildCancel() error = %v, want nil", err)
	}
}

func TestImageService_ImageBuild(t *testing.T) {
	service := &ImageService{}
	opts := build.ImageBuildOptions{}

	response, err := service.ImageBuild(context.Background(), nil, opts)
	if err != nil {
		t.Errorf("ImageBuild() error = %v, want nil", err)
	}

	if response.Body != nil {
		t.Errorf("ImageBuild() response.Body = %v, want nil", response.Body)
	}
}

func TestImageService_ImageCreate(t *testing.T) {
	service := &ImageService{}
	opts := image.CreateOptions{}

	response, err := service.ImageCreate(context.Background(), "test-image", opts)
	if err != nil {
		t.Errorf("ImageCreate() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("ImageCreate() = %v, want nil", response)
	}
}

func TestImageService_ImageHistory(t *testing.T) {
	service := &ImageService{}

	history, err := service.ImageHistory(context.Background(), "test-image")
	if err != nil {
		t.Errorf("ImageHistory() error = %v, want nil", err)
	}

	if history != nil {
		t.Errorf("ImageHistory() = %v, want nil", history)
	}
}

func TestImageService_ImageImport(t *testing.T) {
	service := &ImageService{}
	source := image.ImportSource{}
	opts := image.ImportOptions{}

	response, err := service.ImageImport(context.Background(), source, "test-ref", opts)
	if err != nil {
		t.Errorf("ImageImport() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("ImageImport() = %v, want nil", response)
	}
}

func TestImageService_ImageInspect(t *testing.T) {
	service := &ImageService{}

	tests := []struct {
		name        string
		imageName   string
		wantErr     bool
		wantErrType error
	}{
		{
			name:      "valid image",
			imageName: "alpine:latest",
			wantErr:   false,
		},
		{
			name:      "empty image",
			imageName: "",
			wantErr:   true,
		},
		{
			name:        "notfound image",
			imageName:   "notfound:latest",
			wantErr:     true,
			wantErrType: errdefs.ErrNotFound,
		},
		{
			name:      "notfound image with ignore",
			imageName: "notfound-ignorenotfound:latest",
			wantErr:   false,
		},
		{
			name:        "not-found image",
			imageName:   "not-found:latest",
			wantErr:     true,
			wantErrType: errdefs.ErrNotFound,
		},
		{
			name:      "not-found image with ignore",
			imageName: "not-found-ignore-not-found:latest",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.ImageInspect(context.Background(), tt.imageName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ImageInspect() error = nil, wantErr %v", tt.wantErr)
				}

				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("ImageInspect() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("ImageInspect() error = %v, wantErr %v", err, tt.wantErr)
				}

				if response.ID != "" {
					t.Errorf("ImageInspect() response.ID = %v, want empty", response.ID)
				}
			}
		})
	}
}

func TestImageService_ImageInspectWithRaw(t *testing.T) {
	service := &ImageService{}

	tests := []struct {
		name        string
		imageName   string
		wantErr     bool
		wantErrType error
		wantRaw     bool
	}{
		{
			name:      "valid image",
			imageName: "alpine:latest",
			wantErr:   false,
			wantRaw:   true,
		},
		{
			name:      "empty image",
			imageName: "",
			wantErr:   true,
		},
		{
			name:        "notfound image",
			imageName:   "notfound:latest",
			wantErr:     true,
			wantErrType: errdefs.ErrNotFound,
		},
		{
			name:      "notfound image with ignore",
			imageName: "notfound-ignorenotfound:latest",
			wantErr:   false,
			wantRaw:   true,
		},
		{
			name:        "not-found image",
			imageName:   "not-found:latest",
			wantErr:     true,
			wantErrType: errdefs.ErrNotFound,
		},
		{
			name:      "not-found image with ignore",
			imageName: "not-found-ignore-not-found:latest",
			wantErr:   false,
			wantRaw:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, raw, err := service.ImageInspectWithRaw(context.Background(), tt.imageName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ImageInspectWithRaw() error = nil, wantErr %v", tt.wantErr)
				}

				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("ImageInspectWithRaw() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("ImageInspectWithRaw() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantRaw {
					if len(raw) == 0 {
						t.Errorf("ImageInspectWithRaw() raw = empty, want data")
					}

					var unmarshaled image.InspectResponse
					if err := json.Unmarshal(raw, &unmarshaled); err != nil {
						t.Errorf("ImageInspectWithRaw() raw data invalid JSON: %v", err)
					}

					if response.ID == "" {
						t.Errorf("ImageInspectWithRaw() response.ID = empty, want generated ID")
					}

					if len(response.RepoTags) == 0 {
						t.Errorf("ImageInspectWithRaw() response.RepoTags = empty, want tags")
					}

					if response.Architecture != "amd64" {
						t.Errorf("ImageInspectWithRaw() response.Architecture = %v, want amd64", response.Architecture)
					}

					if response.Os != "linux" {
						t.Errorf("ImageInspectWithRaw() response.Os = %v, want linux", response.Os)
					}
				}
			}
		})
	}
}

func TestImageService_ImageList(t *testing.T) {
	service := &ImageService{}
	opts := image.ListOptions{}

	images, err := service.ImageList(context.Background(), opts)
	if err != nil {
		t.Errorf("ImageList() error = %v, want nil", err)
	}

	if images != nil {
		t.Errorf("ImageList() = %v, want nil", images)
	}
}

func TestImageService_ImageLoad(t *testing.T) {
	service := &ImageService{}

	response, err := service.ImageLoad(context.Background(), nil)
	if err != nil {
		t.Errorf("ImageLoad() error = %v, want nil", err)
	}

	if response.Body != nil {
		t.Errorf("ImageLoad() response.Body = %v, want nil", response.Body)
	}
}

func TestImageService_ImagePull(t *testing.T) {
	service := &ImageService{}
	opts := image.PullOptions{}

	tests := []struct {
		name        string
		imageName   string
		wantErr     bool
		wantErrType error
		wantBody    bool
	}{
		{
			name:      "valid image",
			imageName: "alpine:latest",
			wantErr:   false,
			wantBody:  true,
		},
		{
			name:      "empty image",
			imageName: "",
			wantErr:   true,
		},
		{
			name:        "notfound image",
			imageName:   "notfound:latest",
			wantErr:     true,
			wantErrType: errdefs.ErrNotFound,
		},
		{
			name:      "notfound image with ignore",
			imageName: "notfound-ignorenotfound:latest",
			wantErr:   false,
			wantBody:  true,
		},
		{
			name:        "not-found image",
			imageName:   "not-found:latest",
			wantErr:     true,
			wantErrType: errdefs.ErrNotFound,
		},
		{
			name:      "not-found image with ignore",
			imageName: "not-found-ignore-not-found:latest",
			wantErr:   false,
			wantBody:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.ImagePull(context.Background(), tt.imageName, opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ImagePull() error = nil, wantErr %v", tt.wantErr)
				}

				if tt.wantErrType != nil && !errdefs.IsNotFound(err) {
					t.Errorf("ImagePull() error type = %v, want %v", err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("ImagePull() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantBody {
					if response == nil {
						t.Errorf("ImagePull() response = nil, want body")
					} else {
						defer response.Close()

						body, err := io.ReadAll(response)
						if err != nil {
							t.Errorf("ImagePull() failed to read response body: %v", err)
						}

						bodyStr := string(body)
						if !strings.Contains(bodyStr, "Pulling from") {
							t.Errorf("ImagePull() response body missing 'Pulling from': %s", bodyStr)
						}

						if !strings.Contains(bodyStr, "Digest:") {
							t.Errorf("ImagePull() response body missing 'Digest:': %s", bodyStr)
						}

						if !strings.Contains(bodyStr, "Status:") {
							t.Errorf("ImagePull() response body missing 'Status:': %s", bodyStr)
						}
					}
				}
			}
		})
	}
}

func TestImageService_ImagePush(t *testing.T) {
	service := &ImageService{}
	opts := image.PushOptions{}

	response, err := service.ImagePush(context.Background(), "alpine:latest", opts)
	if err != nil {
		t.Errorf("ImagePush() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("ImagePush() = %v, want nil", response)
	}
}

func TestImageService_ImageRemove(t *testing.T) {
	service := &ImageService{}
	opts := image.RemoveOptions{}

	response, err := service.ImageRemove(context.Background(), "alpine:latest", opts)
	if err != nil {
		t.Errorf("ImageRemove() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("ImageRemove() = %v, want nil", response)
	}
}

func TestImageService_ImageSave(t *testing.T) {
	service := &ImageService{}
	imageIDs := []string{"alpine:latest"}

	response, err := service.ImageSave(context.Background(), imageIDs)
	if err != nil {
		t.Errorf("ImageSave() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("ImageSave() = %v, want nil", response)
	}
}

func TestImageService_ImageSearch(t *testing.T) {
	service := &ImageService{}
	opts := registry.SearchOptions{}

	results, err := service.ImageSearch(context.Background(), "alpine", opts)
	if err != nil {
		t.Errorf("ImageSearch() error = %v, want nil", err)
	}

	if results != nil {
		t.Errorf("ImageSearch() = %v, want nil", results)
	}
}

func TestImageService_ImageTag(t *testing.T) {
	service := &ImageService{}

	err := service.ImageTag(context.Background(), "alpine:latest", "alpine:test")
	if err != nil {
		t.Errorf("ImageTag() error = %v, want nil", err)
	}
}

func TestImageService_ImagesPrune(t *testing.T) {
	service := &ImageService{}
	pruneFilters := filters.Args{}

	report, err := service.ImagesPrune(context.Background(), pruneFilters)
	if err != nil {
		t.Errorf("ImagesPrune() error = %v, want nil", err)
	}

	if report.ImagesDeleted != nil {
		t.Errorf("ImagesPrune() report.ImagesDeleted = %v, want nil", report.ImagesDeleted)
	}

	if report.SpaceReclaimed != 0 {
		t.Errorf("ImagesPrune() report.SpaceReclaimed = %v, want 0", report.SpaceReclaimed)
	}
}

func TestImageService_InterfaceCompliance(_ *testing.T) {
	var _ client.ImageAPIClient = (*ImageService)(nil)
}
