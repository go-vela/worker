// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	"github.com/urfave/cli/v2"

	"github.com/go-vela/compiler/compiler/native"
	"github.com/go-vela/mock/server"

	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/pipeline"
)

func TestLinux_CreateStage(t *testing.T) {
	// setup types
	_file := "testdata/build/stages/basic.yml"
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_pipeline, err := compiler.
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

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		stage   *pipeline.Stage
	}{
		{ // basic stage
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
		{ // stage with step container with image not found
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
		{ // empty stage
			failure: true,
			stage:   new(pipeline.Stage),
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithBuild(_build),
			WithPipeline(_pipeline),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		if len(test.stage.Name) > 0 {
			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}
		}

		err = _engine.CreateStage(context.Background(), test.stage)

		if test.failure {
			if err == nil {
				t.Errorf("CreateStage should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateStage returned err: %v", err)
		}
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
		failure  bool
		stage    *pipeline.Stage
		stageMap *sync.Map
	}{
		{ // basic stage
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
		{ // basic stage with nil stage map
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
		{ // basic stage with error stage map
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
		_engine, err := New(
			WithBuild(_build),
			WithPipeline(new(pipeline.Build)),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		err = _engine.PlanStage(context.Background(), test.stage, test.stageMap)

		if test.failure {
			if err == nil {
				t.Errorf("PlanStage should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("PlanStage returned err: %v", err)
		}
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

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		stage   *pipeline.Stage
	}{
		{ // basic stage
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
		{ // stage with step container with image not found
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
		{ // stage with step container with bad number
			failure: true,
			stage: &pipeline.Stage{
				Name: "echo",
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
		stageMap := new(sync.Map)
		stageMap.Store("echo", make(chan error, 1))

		_engine, err := New(
			WithBuild(_build),
			WithPipeline(new(pipeline.Build)),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		err = _engine.ExecStage(context.Background(), test.stage, stageMap)

		if test.failure {
			if err == nil {
				t.Errorf("ExecStage should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("ExecStage returned err: %v", err)
		}
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

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		stage   *pipeline.Stage
	}{
		{ // basic stage
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
		_engine, err := New(
			WithBuild(_build),
			WithPipeline(new(pipeline.Build)),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		err = _engine.DestroyStage(context.Background(), test.stage)

		if test.failure {
			if err == nil {
				t.Errorf("DestroyStage should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DestroyStage returned err: %v", err)
		}
	}
}
