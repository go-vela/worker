// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// https://github.com/go-vela/worker/blob/main/runtime/kubernetes/kubernetes_test.go#L83
var _pod = &v1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "github-octocat-1",
		Namespace: "test",
		Labels: map[string]string{
			"pipeline": "github-octocat-1",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
	Status: v1.PodStatus{
		Phase: v1.PodRunning,
		ContainerStatuses: []v1.ContainerStatus{
			{
				Name: "step-github-octocat-1-clone",
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						Reason:   "Completed",
						ExitCode: 0,
					},
				},
				Image: "target/vela-git:v0.6.0",
			},
			{
				Name: "step-github-octocat-1-echo",
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						Reason:   "Completed",
						ExitCode: 0,
					},
				},
				Image: "alpine:latest",
			},
		},
	},
	Spec: v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:            "step-github-octocat-1-clone",
				Image:           "target/vela-git:v0.6.0",
				WorkingDir:      "/vela/src/github.com/octocat/helloworld",
				ImagePullPolicy: v1.PullAlways,
			},
			{
				Name:            "step-github-octocat-1-echo",
				Image:           "alpine:latest",
				WorkingDir:      "/vela/src/github.com/octocat/helloworld",
				ImagePullPolicy: v1.PullAlways,
			},
			{
				Name:            "service-github-octocat-1-postgres",
				Image:           "postgres:12-alpine",
				WorkingDir:      "/vela/src/github.com/octocat/helloworld",
				ImagePullPolicy: v1.PullAlways,
			},
		},
		HostAliases: []v1.HostAlias{
			{
				IP: "127.0.0.1",
				Hostnames: []string{
					"postgres.local",
					"echo.local",
				},
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "github-octocat-1",
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			},
		},
	},
}

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

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	_kubernetes, err := kubernetes.NewMock(_pod)
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
			name:    "docker init step container",
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
			name:    "kubernetes init step container",
			failure: false,
			runtime: _kubernetes,
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
			name:    "docker basic step container",
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
			name:    "kubernetes basic step container",
			failure: false,
			runtime: _kubernetes,
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
			name:    "docker step container with image not found",
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
		{
			name:    "kubernetes step container with image not found",
			failure: false,
			runtime: _kubernetes,
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
			name:      "docker empty step container",
			failure:   true,
			runtime:   _docker,
			container: new(pipeline.Container),
		},
		{
			name:      "kubernetes empty step container",
			failure:   true,
			runtime:   _kubernetes,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		name := filepath.Join(test.runtime.Driver(), test.name)

		t.Run(name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(new(pipeline.Build)),
				WithRepo(_repo),
				WithRuntime(test.runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", name, err)
			}

			err = _engine.CreateStep(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s CreateStep should have returned err", name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s CreateStep returned err: %v", name, err)
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
