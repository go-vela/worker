// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"net/http/httptest"
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

func TestLinux_CreateService(t *testing.T) {
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
			name:    "docker-basic service container",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "kubernetes-basic service container",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-service container with image not found",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_echo",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:notfound",
				Name:        "echo",
				Number:      1,
				Pull:        "not_present",
			},
		},
		//{
		//	name:    "kubernetes-service container with image not found",
		//	failure: true, // FIXME: notfound image w/ k8s mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "service_github_octocat_1_echo",
		//		Directory:   "/vela/src/github.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "alpine:notfound",
		//		Name:        "echo",
		//		Number:      1,
		//		Pull:        "not_present",
		//	},
		//},
		{
			name:      "docker-empty service container",
			failure:   true,
			runtime:   _docker,
			container: new(pipeline.Container),
		},
		{
			name:      "kubernetes-empty service container",
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
				WithRepo(_repo),
				WithRuntime(test.runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.CreateService(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s CreateService should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s CreateService returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_PlanService(t *testing.T) {
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
			name:    "basic service container",
			failure: false,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:    "service container with nil environment",
			failure: true,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: nil,
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:      "empty service container",
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

			err = _engine.PlanService(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("PlanService should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("PlanService returned err: %v", err)
			}
		})
	}
}

func TestLinux_ExecService(t *testing.T) {
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
			name:    "basic service container",
			failure: false,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:    "service container with image not found",
			failure: true,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:notfound",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:      "empty service container",
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
				_engine.services.Store(test.container.ID, new(library.Service))
				_engine.serviceLogs.Store(test.container.ID, new(library.Log))
			}

			err = _engine.ExecService(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("ExecService should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ExecService returned err: %v", err)
			}
		})
	}
}

func TestLinux_StreamService(t *testing.T) {
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
			name:    "basic service container",
			failure: false,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:    "service container with name not found",
			failure: true,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_notfound",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "notfound",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:      "empty service container",
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

			if !test.container.Empty() {
				_engine.services.Store(test.container.ID, new(library.Service))
				_engine.serviceLogs.Store(test.container.ID, new(library.Log))
			}

			err = _engine.StreamService(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("StreamService should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("StreamService returned err: %v", err)
			}
		})
	}
}

func TestLinux_DestroyService(t *testing.T) {
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
			name:    "basic service container",
			failure: false,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:    "service container with ignoring name not found",
			failure: true,
			container: &pipeline.Container{
				ID:          "service_github_octocat_1_ignorenotfound",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "ignorenotfound",
				Number:      1,
				Ports:       []string{"5432:5432"},
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

			err = _engine.DestroyService(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("DestroyService should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("DestroyService returned err: %v", err)
			}
		})
	}
}
