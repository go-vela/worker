// SPDX-License-Identifier: Apache-2.0

package local

import (
	"context"
	"errors"
	"flag"
	"sync"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime/docker"
)

func TestLocal_CreateStage(t *testing.T) {
	// setup types
	_file := "testdata/build/stages/basic.yml"
	_build := testBuild()
	_repo := testRepo()

	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_pipeline, _, err := compiler.
		Duplicate().
		WithBuild(_build).
		WithRepo(_repo).
		WithLocal(true).
		WithUser(_repo.GetOwner()).
		Compile(_file)
	if err != nil {
		t.Errorf("unable to compile pipeline %s: %v", _file, err)
	}

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		stage   *pipeline.Stage
	}{
		{
			name:    "basic stage",
			failure: false,
			stage: &pipeline.Stage{
				Name: "echo",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
		},
		{
			name:    "stage with step container with image not found",
			failure: true,
			stage: &pipeline.Stage{
				Name: "echo",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:notfound",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRepo(_repo),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			err = _engine.CreateStage(context.Background(), test.stage)

			if test.failure {
				if err == nil {
					t.Errorf("CreateStage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("CreateStage returned err: %v", err)
			}
		})
	}
}

func TestLocal_PlanStage(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	testMap := new(sync.Map)
	testMap.Store("foo", make(chan error, 1))

	tm, _ := testMap.Load("foo")
	tm.(chan error) <- nil
	close(tm.(chan error))

	errMap := new(sync.Map)
	errMap.Store("foo", make(chan error, 1))

	em, _ := errMap.Load("foo")
	em.(chan error) <- errors.New("bar")
	close(em.(chan error))

	// setup tests
	tests := []struct {
		name     string
		failure  bool
		stage    *pipeline.Stage
		stageMap *sync.Map
	}{
		{
			name:    "basic stage",
			failure: false,
			stage: &pipeline.Stage{
				Name: "echo",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
			stageMap: new(sync.Map),
		},
		{
			name:    "basic stage with nil stage map",
			failure: false,
			stage: &pipeline.Stage{
				Name:  "echo",
				Needs: []string{"foo"},
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
			stageMap: testMap,
		},
		{
			name:    "basic stage with error stage map",
			failure: true,
			stage: &pipeline.Stage{
				Name:  "echo",
				Needs: []string{"foo"},
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
			stageMap: errMap,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRepo(_repo),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.PlanStage(context.Background(), test.stage, test.stageMap)

			if test.failure {
				if err == nil {
					t.Errorf("PlanStage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("PlanStage returned err: %v", err)
			}
		})
	}
}

func TestLocal_ExecStage(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	// setup tests
	tests := []struct {
		name    string
		failure bool
		stage   *pipeline.Stage
	}{
		{
			name:    "basic stage",
			failure: false,
			stage: &pipeline.Stage{
				Name: "echo",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
		},
		{
			name:    "stage with step container with image not found",
			failure: true,
			stage: &pipeline.Stage{
				Name: "echo",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:notfound",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stageMap := new(sync.Map)
			stageMap.Store("echo", make(chan error))

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRepo(_repo),
				WithRuntime(_runtime),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.ExecStage(context.Background(), test.stage, stageMap)

			if test.failure {
				if err == nil {
					t.Errorf("ExecStage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ExecStage returned err: %v", err)
			}
		})
	}
}

func TestLocal_DestroyStage(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		stage   *pipeline.Stage
	}{
		{
			name:    "basic stage",
			failure: false,
			stage: &pipeline.Stage{
				Name: "echo",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRepo(_repo),
				WithRuntime(_runtime),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.DestroyStage(context.Background(), test.stage)

			if test.failure {
				if err == nil {
					t.Errorf("DestroyStage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("DestroyStage returned err: %v", err)
			}
		})
	}
}
