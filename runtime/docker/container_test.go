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
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.InspectContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("InspectContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectContainer returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
			failure:   true,
			container: new(pipeline.Container),
		},
		{
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
		err = _engine.RemoveContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("RemoveContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RemoveContainer returned err: %v", err)
		}
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
		failure   bool
		pipeline  *pipeline.Build
		container *pipeline.Container
		volumes   []string
	}{
		{
			failure:   false,
			pipeline:  _pipeline,
			container: _container,
		},
		{
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
			failure:   true,
			pipeline:  _pipeline,
			container: new(pipeline.Container),
		},
		{
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
	}

	// run tests
	for _, test := range tests {
		if len(test.volumes) > 0 {
			_engine.config.Volumes = test.volumes
		}

		runtimeChannel := make(chan struct{})

		err = _engine.RunContainer(context.Background(), test.container, test.pipeline, runtimeChannel)

		if test.failure {
			if err == nil {
				t.Errorf("RunContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RunContainer returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
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
		err = _engine.SetupContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("SetupContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("SetupContainer returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// pass a closed channel to let TailContainer start right away
	runtimeChannel := make(chan struct{})
	close(runtimeChannel)

	// run tests
	for _, test := range tests {
		_, err = _engine.TailContainer(context.Background(), test.container, runtimeChannel)

		if test.failure {
			if err == nil {
				t.Errorf("TailContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("TailContainer returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.WaitContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("WaitContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WaitContainer returned err: %v", err)
		}
	}
}
