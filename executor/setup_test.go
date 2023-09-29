// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/mock/server"

	"github.com/go-vela/worker/executor/linux"
	"github.com/go-vela/worker/executor/local"

	"github.com/go-vela/worker/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/constants"
)

func TestExecutor_Setup_Darwin(t *testing.T) {
	// setup types
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

	_setup := &Setup{
		Build:    _build,
		Client:   _client,
		Driver:   constants.DriverDarwin,
		Pipeline: _pipeline,
		Repo:     _repo,
		Runtime:  _runtime,
		User:     _user,
	}

	got, err := _setup.Darwin()
	if err == nil {
		t.Errorf("Darwin should have returned err")
	}

	if got != nil {
		t.Errorf("Darwin is %v, want nil", got)
	}
}

func TestExecutor_Setup_Linux(t *testing.T) {
	// setup types
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

	want, err := linux.New(
		linux.WithBuild(_build),
		linux.WithMaxLogSize(2097152),
		linux.WithLogStreamingTimeout(1*time.Second),
		linux.WithHostname("localhost"),
		linux.WithPipeline(_pipeline),
		linux.WithRepo(_repo),
		linux.WithRuntime(_runtime),
		linux.WithUser(_user),
		linux.WithVelaClient(_client),
		linux.WithVersion("v1.0.0"),
	)
	if err != nil {
		t.Errorf("unable to create linux engine: %v", err)
	}

	_setup := &Setup{
		Build:      _build,
		Client:     _client,
		Driver:     constants.DriverLinux,
		MaxLogSize: 2097152,
		Hostname:   "localhost",
		Pipeline:   _pipeline,
		Repo:       _repo,
		Runtime:    _runtime,
		User:       _user,
		Version:    "v1.0.0",
	}

	// run test
	got, err := _setup.Linux()
	if err != nil {
		t.Errorf("Linux returned err: %v", err)
	}

	// Comparing with reflect.DeepEqual(x, y interface) panics due to the
	// unexported streamRequests channel.
	if diff := cmp.Diff(want, got, cmp.Comparer(linux.Equal)); diff != "" {
		t.Errorf("linux Engine mismatch (-want +got):\n%v", diff)
	}
}

func TestExecutor_Setup_Local(t *testing.T) {
	// setup types
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

	want, err := local.New(
		local.WithBuild(_build),
		local.WithHostname("localhost"),
		local.WithPipeline(_pipeline),
		local.WithRepo(_repo),
		local.WithRuntime(_runtime),
		local.WithUser(_user),
		local.WithVelaClient(_client),
		local.WithVersion("v1.0.0"),
	)
	if err != nil {
		t.Errorf("unable to create local engine: %v", err)
	}

	_setup := &Setup{
		Build:    _build,
		Client:   _client,
		Driver:   "local",
		Hostname: "localhost",
		Pipeline: _pipeline,
		Repo:     _repo,
		Runtime:  _runtime,
		User:     _user,
		Version:  "v1.0.0",
	}

	// run test
	got, err := _setup.Local()
	if err != nil {
		t.Errorf("Local returned err: %v", err)
	}

	// Comparing with reflect.DeepEqual(x, y interface) panics due to the
	// unexported streamRequests channel.
	if diff := cmp.Diff(want, got, cmp.Comparer(local.Equal)); diff != "" {
		t.Errorf("local Engine mismatch (-want +got):\n%v", diff)
	}
}

func TestExecutor_Setup_Windows(t *testing.T) {
	// setup types
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

	_setup := &Setup{
		Build:    _build,
		Client:   _client,
		Driver:   constants.DriverWindows,
		Pipeline: _pipeline,
		Repo:     _repo,
		Runtime:  _runtime,
		User:     _user,
	}

	got, err := _setup.Windows()
	if err == nil {
		t.Errorf("Windows should have returned err")
	}

	if got != nil {
		t.Errorf("Windows is %v, want nil", got)
	}
}

func TestExecutor_Setup_Validate(t *testing.T) {
	// setup types
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
		name    string
		setup   *Setup
		failure bool
	}{
		{
			name: "complete",
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       _user,
			},
			failure: false,
		},
		{
			name: "nil build",
			setup: &Setup{
				Build:      nil,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       _user,
			},
			failure: true,
		},
		{
			name: "nil client",
			setup: &Setup{
				Build:      _build,
				Client:     nil,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       _user,
			},
			failure: true,
		},
		{
			name: "empty driver",
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     "",
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       _user,
			},
			failure: true,
		},
		{
			name: "nil pipeline",
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   nil,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       _user,
			},
			failure: true,
		},
		{
			name: "nil repo",
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       nil,
				Runtime:    _runtime,
				User:       _user,
			},
			failure: true,
		},
		{
			name: "nil runtime",
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    nil,
				User:       _user,
			},
			failure: true,
		},
		{
			name: "nil user",
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       nil,
			},
			failure: true,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = test.setup.Validate()

			if test.failure {
				if err == nil {
					t.Errorf("Validate should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("Validate returned err: %v", err)
			}
		})
	}
}
