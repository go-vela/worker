// SPDX-License-Identifier: Apache-2.0

package local

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
)

func TestLocal_Opt_WithBuild(t *testing.T) {
	// setup types
	_build := testBuild()

	// setup tests
	tests := []struct {
		name  string
		build *api.Build
	}{
		{
			name:  "build",
			build: _build,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(test.build),
			)
			if err != nil {
				t.Errorf("WithBuild returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.build, _build) {
				t.Errorf("WithBuild is %v, want %v", _engine.build, _build)
			}
		})
	}
}

func TestLocal_Opt_WithHostname(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		hostname string
		want     string
	}{
		{
			name:     "dns hostname",
			hostname: "vela.worker.localhost",
			want:     "vela.worker.localhost",
		},
		{
			name:     "empty hostname is localhost",
			hostname: "",
			want:     "localhost",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithHostname(test.hostname),
			)
			if err != nil {
				t.Errorf("unable to create local engine: %v", err)
			}

			if !reflect.DeepEqual(_engine.Hostname, test.want) {
				t.Errorf("WithHostname is %v, want %v", _engine.Hostname, test.want)
			}
		})
	}
}

func TestLocal_Opt_WithPipeline(t *testing.T) {
	// setup types
	_steps := testSteps()

	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			name:     "steps pipeline",
			failure:  false,
			pipeline: _steps,
		},
		{
			name:     "nil pipeline",
			failure:  true,
			pipeline: nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithPipeline(test.pipeline),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithPipeline should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithPipeline returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.pipeline, _steps) {
				t.Errorf("WithPipeline is %v, want %v", _engine.pipeline, _steps)
			}
		})
	}
}

func TestLocal_Opt_WithRuntime(t *testing.T) {
	// setup types
	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		runtime runtime.Engine
	}{
		{
			name:    "docker runtime",
			failure: false,
			runtime: _runtime,
		},
		{
			name:    "nil runtime",
			failure: true,
			runtime: nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithRuntime(test.runtime),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithRuntime should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithRuntime returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.Runtime, _runtime) {
				t.Errorf("WithRuntime is %v, want %v", _engine.Runtime, _runtime)
			}
		})
	}
}

func TestLocal_Opt_WithVelaClient(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	// setup tests
	tests := []struct {
		name   string
		client *vela.Client
	}{
		{
			name:   "vela client",
			client: _client,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithVelaClient(test.client),
			)
			if err != nil {
				t.Errorf("WithVelaClient returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.Vela, _client) {
				t.Errorf("WithVelaClient is %v, want %v", _engine.Vela, _client)
			}
		})
	}
}

func TestLocal_Opt_WithVersion(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version",
			version: "v1.0.0",
			want:    "v1.0.0",
		},
		{
			name:    "empty version",
			version: "",
			want:    "v0.0.0",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithVersion(test.version),
			)
			if err != nil {
				t.Errorf("unable to create local engine: %v", err)
			}

			if !reflect.DeepEqual(_engine.Version, test.want) {
				t.Errorf("WithVersion is %v, want %v", _engine.Version, test.want)
			}
		})
	}
}

func TestLocal_Opt_WithMockStdout(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		mock    bool
		wantNil bool
	}{
		{
			name:    "standard",
			mock:    false,
			wantNil: true,
		},
		{
			name:    "mocked",
			mock:    true,
			wantNil: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithMockStdout(test.mock),
			)
			if err != nil {
				t.Errorf("unable to create local engine: %v", err)
			}

			if !reflect.DeepEqual(_engine.MockStdout() == nil, test.wantNil) {
				t.Errorf("WithMockStdout is %v, wantNil = %v", _engine.MockStdout() == nil, test.wantNil)
			}
		})
	}
}
