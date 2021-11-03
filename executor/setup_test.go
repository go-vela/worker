// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

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
		linux.WithLogMethod("byte-chunks"),
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
		Build:     _build,
		Client:    _client,
		Driver:    constants.DriverLinux,
		LogMethod: "byte-chunks",
		Hostname:  "localhost",
		Pipeline:  _pipeline,
		Repo:      _repo,
		Runtime:   _runtime,
		User:      _user,
		Version:   "v1.0.0",
	}

	// run test
	got, err := _setup.Linux()
	if err != nil {
		t.Errorf("Linux returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Linux is %v, want %v", got, want)
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

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Local is %v, want %v", got, want)
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
		setup   *Setup
		failure bool
	}{
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: false,
		},
		{
			setup: &Setup{
				Build:     nil,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    nil,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    "",
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  nil,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      nil,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   nil,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "byte-chunks",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      nil,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
		{
			setup: &Setup{
				Build:     _build,
				Client:    _client,
				Driver:    constants.DriverLinux,
				LogMethod: "foobar",
				Pipeline:  _pipeline,
				Repo:      _repo,
				Runtime:   _runtime,
				User:      _user,
			},
			failure: true,
		},
	}

	// run tests
	for _, test := range tests {
		err = test.setup.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
