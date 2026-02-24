// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"io"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestDocker_PollFileNames(t *testing.T) {
	// setup Docker
	_engine, err := NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		step      *pipeline.Container
		wantFiles int
	}{
		{
			name:    "artifacts directory search",
			failure: false,
			container: &pipeline.Container{
				ID:    "artifacts-container",
				Image: "alpine:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{"artifacts/*/*.txt"},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
			wantFiles: 4, // alpha.txt and beta.txt in test_results and build_results
		},
		{
			name:    "directory not found",
			failure: true,
			container: &pipeline.Container{
				ID:    "artifacts-container",
				Image: "alpine:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{"artifacts/*.txt"},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/not-found",
				},
			},
		},
		{
			name:    "empty container image",
			failure: false,
			container: &pipeline.Container{
				ID:    "no-image",
				Image: "",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{"*.txt"},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := _engine.PollFileNames(context.Background(), test.container, test.step)

			if test.failure {
				if err == nil {
					t.Errorf("PollFileNames should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("PollFileNames returned err: %v", err)
			}

			if test.wantFiles > 0 && len(got) != test.wantFiles {
				t.Errorf("PollFileNames returned %d files, want %d", len(got), test.wantFiles)
			}
		})
	}
}

func TestDocker_PollFileContent(t *testing.T) {
	// setup Docker
	_engine, err := NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		path      string
		wantBytes []byte
	}{
		{
			name:    "file content from default path",
			failure: false,
			container: &pipeline.Container{
				ID:    "content-container",
				Image: "alpine:latest",
			},
			path:      "/vela/artifacts/test_results/alpha.txt",
			wantBytes: []byte("results"),
		},
		{
			name:    "path not found",
			failure: false,
			container: &pipeline.Container{
				ID:    "content-container",
				Image: "alpine:latest",
			},
			path: "not-found",
		},
		{
			name:      "empty container image",
			failure:   false,
			container: new(pipeline.Container),
			path:      "/some/path",
		},
		{
			name:    "empty path",
			failure: false,
			container: &pipeline.Container{
				ID:    "content-container",
				Image: "alpine:latest",
			},
			path: "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader, size, err := _engine.PollFileContent(context.Background(), test.container, test.path)

			if test.failure {
				if err == nil {
					t.Errorf("PollFileContent should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("PollFileContent returned err: %v", err)
			}

			if test.wantBytes != nil {
				got, err := io.ReadAll(reader)
				if err != nil {
					t.Errorf("failed to read content: %v", err)
				}

				if string(got) != string(test.wantBytes) {
					t.Errorf("PollFileContent is %s, want %s", string(got), string(test.wantBytes))
				}

				if size != int64(len(test.wantBytes)) {
					t.Errorf("PollFileContent size is %d, want %d", size, int64(len(test.wantBytes)))
				}
			}
		})
	}
}
