// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_InspectImage(t *testing.T) {
	// setup types
	_engine, err := NewMock()
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
			name:      "tag exists",
			failure:   false,
			container: _container,
		},
		{
			name:      "empty build container",
			failure:   true,
			container: new(pipeline.Container),
		},
		{
			name:    "tag notfound",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:notfound",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.InspectImage(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("InspectImage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectImage returned err: %v", err)
			}
		})
	}
}
