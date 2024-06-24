// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"
)

func TestLinux_CreateStep(t *testing.T) {
	// setup types
	_build := testBuild()

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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "docker-init step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-init step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-basic step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-basic step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-step container with image not found",
			failure: true,
			runtime: _docker,
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
		//{
		//	name:    "kubernetes-step container with image not found",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "step-github-octocat-1-echo",
		//		Directory:   "/vela/src/github.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "alpine:notfound",
		//		Name:        "echo",
		//		Number:      1,
		//		Pull:        "not_present",
		//	},
		//},
		{
			name:      "docker-empty step container",
			failure:   true,
			runtime:   _docker,
			container: new(pipeline.Container),
		},
		{
			name:      "kubernetes-empty step container",
			failure:   true,
			runtime:   _kubernetes,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.CreateStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s CreateStep should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s CreateStep returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_PlanStep(t *testing.T) {
	// setup types
	_build := testBuild()

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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "docker-basic step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-basic step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-step container with nil environment",
			failure: true,
			runtime: _docker,
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
			name:    "kubernetes-step container with nil environment",
			failure: true,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: nil,
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:      "docker-empty step container",
			failure:   true,
			runtime:   _docker,
			container: new(pipeline.Container),
		},
		{
			name:      "kubernetes-empty step container",
			failure:   true,
			runtime:   _kubernetes,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.PlanStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s PlanStep should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s PlanStep returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_ExecStep(t *testing.T) {
	// setup types
	_build := testBuild()

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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	_kubernetes.PodTracker.Start(context.Background())

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "docker-init step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-init step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-basic step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-basic step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-detached step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-detached step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
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
			name:    "docker-step container with image not found",
			failure: true,
			runtime: _docker,
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
		//{
		//	name:    "kubernetes-step container with image not found",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "step-github-octocat-1-echo",
		//		Directory:   "/vela/src/github.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "alpine:notfound",
		//		Name:        "echo",
		//		Number:      1,
		//		Pull:        "not_present",
		//	},
		//},
		{
			name:      "docker-empty step container",
			failure:   true,
			runtime:   _docker,
			container: new(pipeline.Container),
		},
		{
			name:      "kubernetes-empty step container",
			failure:   true,
			runtime:   _kubernetes,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			if !test.container.Empty() {
				_engine.steps.Store(test.container.ID, new(library.Step))
				_engine.stepLogs.Store(test.container.ID, new(library.Log))
			}

			err = _engine.ExecStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s ExecStep should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s ExecStep returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_StreamStep(t *testing.T) {
	// setup types
	_build := testBuild()
	_logs := new(library.Log)

	// fill log with bytes
	_logs.SetData(make([]byte, 1000))

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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		logs      *library.Log
		container *pipeline.Container
	}{
		{
			name:    "docker-init step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-init step container",
			failure: false,
			runtime: _kubernetes,
			logs:    _logs,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-basic step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-basic step container",
			failure: false,
			runtime: _kubernetes,
			logs:    _logs,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-step container with name not found",
			failure: true,
			runtime: _docker,
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
		//{
		//	name:    "kubernetes-step container with name not found",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	logs:    _logs,
		//	container: &pipeline.Container{
		//		ID:          "step-github-octocat-1-notfound",
		//		Directory:   "/vela/src/github.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "alpine:latest",
		//		Name:        "notfound",
		//		Number:      1,
		//		Pull:        "not_present",
		//	},
		//},
		{
			name:      "docker-empty step container",
			failure:   true,
			runtime:   _docker,
			logs:      _logs,
			container: new(pipeline.Container),
		},
		{
			name:      "kubernetes-empty step container",
			failure:   true,
			runtime:   _kubernetes,
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
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			if !test.container.Empty() {
				_engine.steps.Store(test.container.ID, new(library.Step))
				_engine.stepLogs.Store(test.container.ID, new(library.Log))
			}

			err = _engine.StreamStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s StreamStep should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s StreamStep returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_DestroyStep(t *testing.T) {
	// setup types
	_build := testBuild()

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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "docker-init step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-init step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-basic step container",
			failure: false,
			runtime: _docker,
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
			name:    "kubernetes-basic step container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-step container with ignoring name not found",
			failure: true,
			runtime: _docker,
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
		//{
		//	name:    "kubernetes-step container with ignoring name not found",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "step-github-octocat-1-ignorenotfound",
		//		Directory:   "/vela/src/github.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "alpine:latest",
		//		Name:        "ignorenotfound",
		//		Number:      1,
		//		Pull:        "not_present",
		//	},
		//},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.DestroyStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s DestroyStep should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s DestroyStep returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_getSecretValues(t *testing.T) {
	fileSecret, err := os.ReadFile("./testdata/step/secret_text.txt")
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
