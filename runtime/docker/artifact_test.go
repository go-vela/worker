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
			name:    "test-results XML files",
			failure: false,
			container: &pipeline.Container{
				ID:    "test-results-container",
				Image: "alpine:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{"test-results/*.xml"},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
			wantFiles: 2, // junit.xml and report.xml
		},
		{
			name:    "cypress screenshots PNG files",
			failure: false,
			container: &pipeline.Container{
				ID:    "cypress-screenshots-container",
				Image: "cypress/browsers:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{"cypress/screenshots/**/*.png"},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
			wantFiles: 3, // screenshot1.png, screenshot2.png, error.png
		},
		{
			name:    "cypress videos MP4 files",
			failure: false,
			container: &pipeline.Container{
				ID:    "cypress-videos-container",
				Image: "cypress/browsers:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{"cypress/videos/**/*.mp4"},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
			wantFiles: 2, // test1.mp4, test2.mp4
		},
		{
			name:    "multiple cypress patterns",
			failure: false,
			container: &pipeline.Container{
				ID:    "cypress-all-container",
				Image: "cypress/browsers:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{
						"cypress/screenshots/**/*.png",
						"cypress/videos/**/*.mp4",
					},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
			wantFiles: 5, // 3 PNG screenshots + 2 MP4 videos
		},
		{
			name:    "combined test-results and cypress artifacts",
			failure: false,
			container: &pipeline.Container{
				ID:    "combined-artifacts-container",
				Image: "cypress/browsers:latest",
			},
			step: &pipeline.Container{
				Artifacts: pipeline.Artifacts{
					Paths: []string{
						"test-results/*.xml",
						"cypress/screenshots/**/*.png",
						"cypress/videos/**/*.mp4",
					},
				},
				Environment: map[string]string{
					"VELA_WORKSPACE": "/vela/workspace",
				},
			},
			wantFiles: 7, // 2 XML + 3 PNG + 2 MP4
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
			name:    "test-results XML file content",
			failure: false,
			container: &pipeline.Container{
				ID:    "content-container",
				Image: "alpine:latest",
			},
			path:      "/vela/workspace/test-results/junit.xml",
			wantBytes: []byte("<?xml version=\"1.0\"?><testsuites></testsuites>"),
		},
		{
			name:    "cypress screenshot PNG file",
			failure: false,
			container: &pipeline.Container{
				ID:    "content-container",
				Image: "cypress/browsers:latest",
			},
			path:      "/vela/workspace/cypress/screenshots/test1/screenshot1.png",
			wantBytes: []byte("PNG_BINARY_DATA"),
		},
		{
			name:    "cypress video MP4 file",
			failure: false,
			container: &pipeline.Container{
				ID:    "content-container",
				Image: "cypress/browsers:latest",
			},
			path:      "/vela/workspace/cypress/videos/test1.mp4",
			wantBytes: []byte("MP4_BINARY_DATA"),
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
