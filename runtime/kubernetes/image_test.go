// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"strings"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

func TestKubernetes_CreateImage(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		container *pipeline.Container
	}{
		{
			name:      "valid container",
			container: _container,
		},
		{
			name: "different container",
			container: &pipeline.Container{
				ID:     "different",
				Image:  "alpine:latest",
				Number: 2,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := _engine.CreateImage(context.Background(), test.container)
			if err != nil {
				t.Errorf("CreateImage returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_InspectImage(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
	}{
		{
			name:      "valid image",
			failure:   false,
			container: _container,
		},
		{
			name:    "pull on start policy",
			failure: false,
			container: &pipeline.Container{
				ID:     "test_container",
				Image:  "alpine:latest",
				Number: 1,
				Pull:   constants.PullOnStart,
			},
		},
		{
			name:    "pull always policy",
			failure: false,
			container: &pipeline.Container{
				ID:     "test_container",
				Image:  "alpine:latest",
				Number: 1,
				Pull:   constants.PullAlways,
			},
		},
		{
			name:    "pull never policy",
			failure: false,
			container: &pipeline.Container{
				ID:     "test_container",
				Image:  "alpine:latest",
				Number: 1,
				Pull:   constants.PullNever,
			},
		},
		{
			name:    "pull not present policy",
			failure: false,
			container: &pipeline.Container{
				ID:     "test_container",
				Image:  "alpine:latest",
				Number: 1,
				Pull:   constants.PullNotPresent,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := _engine.InspectImage(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("InspectImage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectImage returned err: %v", err)
			}

			if output == nil {
				t.Error("InspectImage returned nil output")
			}

			outputStr := string(output)

			// Check for pull on start special case
			if strings.EqualFold(test.container.Pull, constants.PullOnStart) {
				if !strings.Contains(outputStr, "skipped for container") {
					t.Errorf("Expected skip message for pull on start, got: %s", outputStr)
				}

				if !strings.Contains(outputStr, test.container.ID) {
					t.Errorf("Expected container ID %s in output: %s", test.container.ID, outputStr)
				}
			} else {
				// Should contain kubectl command
				if !strings.Contains(outputStr, "kubectl get pod") {
					t.Errorf("Expected kubectl command in output: %s", outputStr)
				}

				if !strings.Contains(outputStr, test.container.ID) {
					t.Errorf("Expected container ID %s in output: %s", test.container.ID, outputStr)
				}
			}
		})
	}
}

func TestKubernetes_InspectImage_EdgeCases(t *testing.T) {
	// Test edge cases for InspectImage
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// Test with empty container ID
	container := &pipeline.Container{
		ID:     "",
		Image:  "alpine:latest",
		Number: 1,
		Pull:   constants.PullAlways,
	}

	_, err = _engine.InspectImage(context.Background(), container)
	// This may or may not error, just test that it doesn't panic
	if err != nil {
		t.Logf("InspectImage with empty ID returned expected error: %v", err)
	}
}

func TestImageConstants(t *testing.T) {
	// Test that the constants are defined correctly
	if pauseImage != "kubernetes/pause:latest" {
		t.Errorf("pauseImage constant = %s, want kubernetes/pause:latest", pauseImage)
	}

	// Test that imagePatch contains expected format strings
	if !strings.Contains(imagePatch, "%s") {
		t.Errorf("imagePatch should contain %%s format specifiers")
	}

	if !strings.Contains(imagePatch, "spec") {
		t.Error("imagePatch should contain 'spec' field")
	}

	if !strings.Contains(imagePatch, "containers") {
		t.Error("imagePatch should contain 'containers' field")
	}

	if !strings.Contains(imagePatch, "name") {
		t.Error("imagePatch should contain 'name' field")
	}

	if !strings.Contains(imagePatch, "image") {
		t.Error("imagePatch should contain 'image' field")
	}
}
