// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"errors"
	"flag"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/urfave/cli/v2"
)

func TestLinux_CreateStage(t *testing.T) {
	// setup types
	_file := "testdata/build/stages/basic.yml"
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", "target/vela-git:latest", "doc")
	compiler, _ := native.New(cli.NewContext(nil, set, nil))

	_pipeline, _, err := compiler.
		Duplicate().
		WithBuild(_build).
		WithRepo(_repo).
		WithMetadata(_metadata).
		WithUser(_user).
		Compile(_file)
	if err != nil {
		t.Errorf("unable to compile pipeline %s: %v", _file, err)
	}

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		runtime runtime.Engine
		stage   *pipeline.Stage
	}{
		{
			name:    "docker-basic stage",
			failure: false,
			runtime: _docker,
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
			name:    "docker-stage with step container with image not found",
			failure: true,
			runtime: _docker,
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
		{
			name:    "docker-empty stage",
			failure: true,
			runtime: _docker,
			stage:   new(pipeline.Stage),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRepo(_repo),
				WithRuntime(test.runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			if len(test.stage.Name) > 0 {
				// run create to init steps to be created properly
				err = _engine.CreateBuild(context.Background())
				if err != nil {
					t.Errorf("unable to create %s build: %v", test.name, err)
				}
			}

			err = _engine.CreateStage(context.Background(), test.stage)

			if test.failure {
				if err == nil {
					t.Errorf("%s CreateStage should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s CreateStage returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_PlanStage(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	dockerTestMap := new(sync.Map)
	dockerTestMap.Store("foo", make(chan error, 1))

	dtm, _ := dockerTestMap.Load("foo")
	dtm.(chan error) <- nil
	close(dtm.(chan error))

	dockerErrMap := new(sync.Map)
	dockerErrMap.Store("foo", make(chan error, 1))

	dem, _ := dockerErrMap.Load("foo")
	dem.(chan error) <- errors.New("bar")
	close(dem.(chan error))

	// setup tests
	tests := []struct {
		name     string
		failure  bool
		runtime  runtime.Engine
		stage    *pipeline.Stage
		stageMap *sync.Map
	}{
		{
			name:    "docker-basic stage",
			failure: false,
			runtime: _docker,
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
			name:    "docker-basic stage with nil stage map",
			failure: false,
			runtime: _docker,
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
			stageMap: dockerTestMap,
		},
		{
			name:    "docker-basic stage with error stage map",
			failure: true,
			runtime: _docker,
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
			stageMap: dockerErrMap,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRepo(_repo),
				WithRuntime(test.runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.PlanStage(context.Background(), test.stage, test.stageMap)

			if test.failure {
				if err == nil {
					t.Errorf("%s PlanStage should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s PlanStage returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_ExecStage(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	// setup tests
	tests := []struct {
		name    string
		failure bool
		runtime runtime.Engine
		stage   *pipeline.Stage
	}{
		{
			name:    "docker-basic stage",
			failure: false,
			runtime: _docker,
			stage: &pipeline.Stage{
				Independent: true,
				Name:        "echo",
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
			name:    "docker-stage with step container with image not found",
			failure: true,
			runtime: _docker,
			stage: &pipeline.Stage{
				Name:        "echo",
				Independent: true,
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
		{
			name:    "docker-stage with step container with bad number",
			failure: true,
			runtime: _docker,
			stage: &pipeline.Stage{
				Name:        "echo",
				Independent: true,
				Steps: pipeline.ContainerSlice{
					{
						ID:          "github_octocat_1_echo_echo",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      0,
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
			stageMap.Store("echo", make(chan error, 1))

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRepo(_repo),
				WithRuntime(test.runtime),
				WithUser(_user),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.ExecStage(context.Background(), test.stage, stageMap)

			if test.failure {
				if err == nil {
					t.Errorf("%s ExecStage should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s ExecStage returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_DestroyStage(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		runtime runtime.Engine
		stage   *pipeline.Stage
	}{
		{
			name:    "docker-basic stage",
			failure: false,
			runtime: _docker,
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
				WithRuntime(test.runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.DestroyStage(context.Background(), test.stage)

			if test.failure {
				if err == nil {
					t.Errorf("%s DestroyStage should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s DestroyStage returned err: %v", test.name, err)
			}
		})
	}
}
