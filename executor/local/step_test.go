// SPDX-License-Identifier: Apache-2.0

package local

import (
	"context"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime/docker"
)

func TestLocal_CreateStep(t *testing.T) {
	// setup types
	_build := testBuild()

	_runtime, err := docker.NewMock()
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
			name:    "init step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        constants.InitName,
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "basic step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "step container with image not found",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:notfound",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:      "empty step container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.CreateStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("CreateStep should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("CreateStep returned err: %v", err)
			}
		})
	}
}

func TestLocal_PlanStep(t *testing.T) {
	// setup types
	_build := testBuild()

	_runtime, err := docker.NewMock()
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
			name:    "basic step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:      "empty step container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.PlanStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("PlanStep should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("PlanStep returned err: %v", err)
			}
		})
	}
}

func TestLocal_ExecStep(t *testing.T) {
	// setup types
	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
	}{
		{
			name:    "init step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        constants.InitName,
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "basic step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "detached step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "step container with image not found",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:notfound",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:      "empty step container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(_runtime),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			if !test.container.Empty() {
				_engine.steps.Store(test.container.ID, new(api.Step))
			}

			err = _engine.ExecStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("ExecStep should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ExecStep returned err: %v", err)
			}
		})
	}
}

func TestLocal_StreamStep(t *testing.T) {
	// setup types
	_build := testBuild()

	_runtime, err := docker.NewMock()
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
			name:    "init step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        constants.InitName,
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "basic step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "basic stage container",
			failure: false,
			container: &pipeline.Container{
				ID:          "github_octocat_1_echo_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"VELA_STEP_STAGE": "foo"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:      "empty step container",
			failure:   true,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.StreamStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("StreamStep should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("StreamStep returned err: %v", err)
			}
		})
	}
}

func TestLocal_DestroyStep(t *testing.T) {
	// setup types
	_build := testBuild()

	_runtime, err := docker.NewMock()
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
			name:    "init step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        constants.InitName,
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "basic step container",
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.DestroyStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("DestroyStep should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("DestroyStep returned err: %v", err)
			}
		})
	}
}
