// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/mock/server"

	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func TestLinux_CreateStep(t *testing.T) {
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
				Name:        "init",
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
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
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

func TestLinux_PlanStep(t *testing.T) {
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
			name:    "step container with nil environment",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: nil,
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
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
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

func TestLinux_ExecStep(t *testing.T) {
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
				Name:        "init",
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
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			if !test.container.Empty() {
				_engine.steps.Store(test.container.ID, new(library.Step))
				_engine.stepLogs.Store(test.container.ID, new(library.Log))
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

func TestLinux_StreamStep(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_logs := new(library.Log)

	// fill log with bytes
	_logs.SetData(make([]byte, 1000))

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
		name      string
		failure   bool
		logs      *library.Log
		container *pipeline.Container
	}{
		{
			name:    "init step container",
			failure: false,
			logs:    _logs,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "basic step container",
			failure: false,
			logs:    _logs,
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
			name:    "step container with name not found",
			failure: true,
			logs:    _logs,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_notfound",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "notfound",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:      "empty step container",
			failure:   true,
			logs:      _logs,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithMaxLogSize(10),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			if !test.container.Empty() {
				_engine.steps.Store(test.container.ID, new(library.Step))
				_engine.stepLogs.Store(test.container.ID, new(library.Log))
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

func TestLinux_DestroyStep(t *testing.T) {
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
				Name:        "init",
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
			name:    "step container with ignoring name not found",
			failure: true,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "ignorenotfound",
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
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
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

func TestLinux_getSecretValues(t *testing.T) {
	fileSecret, err := ioutil.ReadFile("./testdata/step/secret_text.txt")
	if err != nil {
		t.Errorf("unable to read from test data file secret. Err: %v", err)
	}

	tests := []struct {
		name      string
		want      []string
		container *pipeline.Container
	}{
		{
			name: "no secrets container",
			want: []string{},
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name: "secrets container",
			want: []string{"secretUser", "secretPass"},
			container: &pipeline.Container{
				ID:        "step_github_octocat_1_echo",
				Directory: "/vela/src/github.com/github/octocat",
				Environment: map[string]string{
					"FOO":             "bar",
					"SECRET_USERNAME": "secretUser",
					"SECRET_PASSWORD": "secretPass",
				},
				Image:  "alpine:latest",
				Name:   "echo",
				Number: 1,
				Pull:   "not_present",
				Secrets: pipeline.StepSecretSlice{
					{
						Source: "someSource",
						Target: "secret_username",
					},
					{
						Source: "someOtherSource",
						Target: "secret_password",
					},
					{
						Source: "disallowedSecret",
						Target: "cannot_find",
					},
				},
			},
		},
		{
			name: "secrets container with file as value",
			want: []string{"secretUser", "this is a secret"},
			container: &pipeline.Container{
				ID:        "step_github_octocat_1_ignorenotfound",
				Directory: "/vela/src/github.com/github/octocat",
				Environment: map[string]string{
					"FOO":             "bar",
					"SECRET_USERNAME": "secretUser",
					"SECRET_PASSWORD": string(fileSecret),
				},
				Image:  "alpine:latest",
				Name:   "ignorenotfound",
				Number: 1,
				Pull:   "not_present",
				Secrets: pipeline.StepSecretSlice{
					{
						Source: "someSource",
						Target: "secret_username",
					},
					{
						Source: "someOtherSource",
						Target: "secret_password",
					},
				},
			},
		},
	}
	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getSecretValues(test.container)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getSecretValues is %v, want %v", got, test.want)
			}
		})
	}
}
