// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/mock/server"

	"github.com/go-vela/pkg-runtime/runtime"
	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func TestLinux_Opt_WithBuild(t *testing.T) {
	// setup types
	_build := testBuild()

	// setup tests
	tests := []struct {
		failure bool
		build   *library.Build
	}{
		{
			failure: false,
			build:   _build,
		},
		{
			failure: true,
			build:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithBuild(test.build),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.build, _build) {
			t.Errorf("WithBuild is %v, want %v", _engine.build, _build)
		}
	}
}

func TestLinux_Opt_WithHostname(t *testing.T) {
	// setup tests
	tests := []struct {
		hostname string
		want     string
	}{
		{
			hostname: "vela.worker.localhost",
			want:     "vela.worker.localhost",
		},
		{
			hostname: "",
			want:     "localhost",
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithHostname(test.hostname),
		)
		if err != nil {
			t.Errorf("unable to create linux engine: %v", err)
		}

		if !reflect.DeepEqual(_engine.Hostname, test.want) {
			t.Errorf("WithHostname is %v, want %v", _engine.Hostname, test.want)
		}
	}
}

func TestLinux_Opt_WithPipeline(t *testing.T) {
	// setup types
	_steps := testSteps()

	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			failure:  false,
			pipeline: _steps,
		},
		{
			failure:  true,
			pipeline: nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithPipeline(test.pipeline),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithPipeline should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithPipeline returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.pipeline, _steps) {
			t.Errorf("WithPipeline is %v, want %v", _engine.pipeline, _steps)
		}
	}
}

func TestLinux_Opt_WithRepo(t *testing.T) {
	// setup types
	_repo := testRepo()

	// setup tests
	tests := []struct {
		failure bool
		repo    *library.Repo
	}{
		{
			failure: false,
			repo:    _repo,
		},
		{
			failure: true,
			repo:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithRepo(test.repo),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithRepo should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithRepo returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.repo, _repo) {
			t.Errorf("WithRepo is %v, want %v", _engine.repo, _repo)
		}
	}
}

func TestLinux_Opt_WithRuntime(t *testing.T) {
	// setup types
	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		runtime runtime.Engine
	}{
		{
			failure: false,
			runtime: _runtime,
		},
		{
			failure: true,
			runtime: nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithRuntime(test.runtime),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithRuntime should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithRuntime returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.Runtime, _runtime) {
			t.Errorf("WithRuntime is %v, want %v", _engine.Runtime, _runtime)
		}
	}
}

func TestLinux_Opt_WithUser(t *testing.T) {
	// setup types
	_user := testUser()

	// setup tests
	tests := []struct {
		failure bool
		user    *library.User
	}{
		{
			failure: false,
			user:    _user,
		},
		{
			failure: true,
			user:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithUser(test.user),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithUser should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithUser returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.user, _user) {
			t.Errorf("WithUser is %v, want %v", _engine.user, _user)
		}
	}
}

func TestLinux_Opt_WithVelaClient(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		client  *vela.Client
	}{
		{
			failure: false,
			client:  _client,
		},
		{
			failure: true,
			client:  nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithVelaClient(test.client),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithVelaClient should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithVelaClient returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.Vela, _client) {
			t.Errorf("WithVelaClient is %v, want %v", _engine.Vela, _client)
		}
	}
}

func TestLinux_Opt_WithVersion(t *testing.T) {
	// setup tests
	tests := []struct {
		version string
		want    string
	}{
		{
			version: "v1.0.0",
			want:    "v1.0.0",
		},
		{
			version: "",
			want:    "v0.0.0",
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithVersion(test.version),
		)
		if err != nil {
			t.Errorf("unable to create linux engine: %v", err)
		}

		if !reflect.DeepEqual(_engine.Version, test.want) {
			t.Errorf("WithVersion is %v, want %v", _engine.Version, test.want)
		}
	}
}
