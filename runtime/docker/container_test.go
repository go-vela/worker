// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_InspectContainer(t *testing.T) {
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
	}{
		{
			name:      "build container",
			failure:   false,
			container: _container,
		},
		{
			name:      "empty build container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.InspectContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("InspectContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectContainer returned err: %v", err)
			}
		})
	}
}

func TestDocker_RemoveContainer(t *testing.T) {
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
	}{
		{
			name:      "build container",
			failure:   false,
			container: _container,
		},
		{
			name:      "empty build container",
			failure:   true,
			container: new(pipeline.Container),
		},
		{
			name:    "absent build container",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "ignorenotfound",
				Number:      2,
				Pull:        "always",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.RemoveContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("RemoveContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RemoveContainer returned err: %v", err)
			}
		})
	}
}

func TestDocker_RunContainer(t *testing.T) {
	// setup Docker
	_engine, err := NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		pipeline  *pipeline.Build
		container *pipeline.Container
		volumes   []string
	}{
		{
			name:      "steps-clone step",
			failure:   false,
			pipeline:  _pipeline,
			container: _container,
		},
		{
			name:     "steps-echo step",
			failure:  false,
			pipeline: _pipeline,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo", "hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      2,
				Pull:        "always",
			},
		},
		{
			name:     "steps-privileged",
			failure:  false,
			pipeline: _pipeline,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo", "hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Image:       "target/vela-docker:latest",
				Name:        "echo",
				Number:      2,
				Pull:        "always",
			},
		},
		{
			name:     "steps-kaniko-volumes",
			failure:  false,
			pipeline: _pipeline,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo", "hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Image:       "target/vela-kaniko:latest",
				Name:        "echo",
				Number:      2,
				Pull:        "always",
			},
			volumes: []string{"/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt:rw"},
		},
		{
			name:      "steps-empty build container",
			failure:   true,
			pipeline:  _pipeline,
			container: new(pipeline.Container),
		},
		{
			name:     "steps-absent build container",
			failure:  true,
			pipeline: _pipeline,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "ignorenotfound",
				Number:      2,
				Pull:        "always",
			},
		},
		{
			name:     "steps-user-absent build container",
			failure:  true,
			pipeline: _pipeline,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "ignorenotfound",
				Number:      2,
				Pull:        "always",
				User:        "foo",
			},
		},
		{
			name:     "steps-user-echo step",
			failure:  false,
			pipeline: _pipeline,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo", "hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      2,
				Pull:        "always",
				User:        "foo",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.volumes) > 0 {
				_engine.config.Volumes = test.volumes
			}

			err = _engine.RunContainer(context.Background(), test.container, test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("RunContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RunContainer returned err: %v", err)
			}
		})
	}
}

func TestDocker_SetupContainer(t *testing.T) {
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
	}{
		{
			name:      "pull-always-tag_exists",
			failure:   false,
			container: _container,
		},
		{
			name:    "pull-not_present-tag_exists",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
		},
		{
			name:    "pull-not_present-mock tag ignorenotfound", // mock returns as if this exists
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:ignorenotfound",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
		},
		{
			name:    "pull-always-tag notfound fails",
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
		{
			name:    "pull-not_present-tag notfound fails",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:notfound",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.SetupContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("SetupContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("SetupContainer returned err: %v", err)
			}
		})
	}
}

func TestDocker_TailContainer(t *testing.T) {
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
	}{
		{
			name:      "build container",
			failure:   false,
			container: _container,
		},
		{
			name:      "empty build container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.TailContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("TailContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("TailContainer returned err: %v", err)
			}
		})
	}
}

func TestDocker_WaitContainer(t *testing.T) {
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
	}{
		{
			name:      "build container",
			failure:   false,
			container: _container,
		},
		{
			name:      "empty build container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.WaitContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("WaitContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WaitContainer returned err: %v", err)
			}
		})
	}
}
