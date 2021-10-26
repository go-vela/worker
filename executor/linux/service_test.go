// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/mock/server"

	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
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

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
	}{
		{ // basic service container
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
		{ // service container with image not found
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
		{ // empty service container
			failure:   true,
			container: new(pipeline.Container),
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

		err = _engine.CreateService(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("CreateService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateService returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{ // basic service container
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
		{ // service container with nil environment
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
		{ // empty service container
			failure:   true,
			container: new(pipeline.Container),
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

		err = _engine.PlanService(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("PlanService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("PlanService returned err: %v", err)
		}
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

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
	}{
		{ // basic service container
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
		{ // service container with image not found
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
		{ // empty service container
			failure:   true,
			container: new(pipeline.Container),
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

		if !test.container.Empty() {
			_engine.services.Store(test.container.ID, new(library.Service))
			_engine.serviceLogs.Store(test.container.ID, new(library.Log))
		}

		err = _engine.ExecService(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("ExecService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("ExecService returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{ // basic service container
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
		{ // service container with name not found
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
		{ // empty service container
			failure:   true,
			container: new(pipeline.Container),
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

		if !test.container.Empty() {
			_engine.services.Store(test.container.ID, new(library.Service))
			_engine.serviceLogs.Store(test.container.ID, new(library.Log))
		}

		err = _engine.StreamService(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("StreamService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("StreamService returned err: %v", err)
		}
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
		failure   bool
		container *pipeline.Container
	}{
		{ // basic service container
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
		{ // service container with ignoring name not found
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

			continue
		}

		if err != nil {
			t.Errorf("DestroyService returned err: %v", err)
		}
	}
}
