// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/mock/server"

	"github.com/go-vela/worker/executor/linux"
	"github.com/go-vela/worker/executor/local"

	"github.com/go-vela/worker/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func TestExecutor_New(t *testing.T) {
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

	_linux, err := linux.New(
		linux.WithBuild(_build),
		linux.WithHostname("localhost"),
		linux.WithMaxLogSize(2097152),
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

	_local, err := local.New(
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

	// setup tests
	tests := []struct {
		name    string
		failure bool
		setup   *Setup
		want    Engine
		equal   interface{}
	}{
		{
			name:    "driver-darwin",
			failure: true,
			setup: &Setup{
				Build:    _build,
				Client:   _client,
				Driver:   constants.DriverDarwin,
				Pipeline: _pipeline,
				Repo:     _repo,
				Runtime:  _runtime,
				User:     _user,
				Version:  "v1.0.0",
			},
			want:  nil,
			equal: reflect.DeepEqual,
		},
		{
			name:    "driver-linux",
			failure: false,
			setup: &Setup{
				Build:      _build,
				Client:     _client,
				Driver:     constants.DriverLinux,
				MaxLogSize: 2097152,
				Pipeline:   _pipeline,
				Repo:       _repo,
				Runtime:    _runtime,
				User:       _user,
				Version:    "v1.0.0",
			},
			want:  _linux,
			equal: linux.Equal,
		},
		{
			name:    "driver-local",
			failure: false,
			setup: &Setup{
				Build:    _build,
				Client:   _client,
				Driver:   "local",
				Pipeline: _pipeline,
				Repo:     _repo,
				Runtime:  _runtime,
				User:     _user,
				Version:  "v1.0.0",
			},
			want:  _local,
			equal: local.Equal,
		},
		{
			name:    "driver-windows",
			failure: true,
			setup: &Setup{
				Build:    _build,
				Client:   _client,
				Driver:   constants.DriverWindows,
				Pipeline: _pipeline,
				Repo:     _repo,
				Runtime:  _runtime,
				User:     _user,
				Version:  "v1.0.0",
			},
			want:  nil,
			equal: reflect.DeepEqual,
		},
		{
			name:    "driver-invalid",
			failure: true,
			setup: &Setup{
				Build:    _build,
				Client:   _client,
				Driver:   "invalid",
				Pipeline: _pipeline,
				Repo:     _repo,
				Runtime:  _runtime,
				User:     _user,
				Version:  "v1.0.0",
			},
			want:  nil,
			equal: reflect.DeepEqual,
		},
		{
			name:    "driver-empty",
			failure: true,
			setup: &Setup{
				Build:    _build,
				Client:   _client,
				Driver:   "",
				Pipeline: _pipeline,
				Repo:     _repo,
				Runtime:  _runtime,
				User:     _user,
				Version:  "v1.0.0",
			},
			want:  nil,
			equal: reflect.DeepEqual,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(test.setup)

			if test.failure {
				if err == nil {
					t.Errorf("New should have returned err")
				}

				if !reflect.DeepEqual(got, test.want) {
					t.Errorf("New is %v, want %v", got, test.want)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			// Comparing with reflect.DeepEqual(x, y interface) panics due to the
			// unexported streamRequests channel.
			if diff := cmp.Diff(test.want, got, cmp.Comparer(test.equal)); diff != "" {
				t.Errorf("engine mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

// setup global variables used for testing.
var (
	_build = &library.Build{
		ID:           vela.Int64(1),
		Number:       vela.Int(1),
		Parent:       vela.Int(1),
		Event:        vela.String("push"),
		Status:       vela.String("success"),
		Error:        vela.String(""),
		Enqueued:     vela.Int64(1563474077),
		Created:      vela.Int64(1563474076),
		Started:      vela.Int64(1563474077),
		Finished:     vela.Int64(0),
		Deploy:       vela.String(""),
		Clone:        vela.String("https://github.com/github/octocat.git"),
		Source:       vela.String("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        vela.String("push received from https://github.com/github/octocat"),
		Message:      vela.String("First commit..."),
		Commit:       vela.String("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       vela.String("OctoKitty"),
		Author:       vela.String("OctoKitty"),
		Branch:       vela.String("master"),
		Ref:          vela.String("refs/heads/master"),
		BaseRef:      vela.String(""),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("linux"),
	}

	_pipeline = &pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Steps: pipeline.ContainerSlice{
			{
				ID:        "step_github_octocat_1_init",
				Directory: "/home/github/octocat",
				Image:     "#init",
				Name:      "init",
				Number:    1,
				Pull:      "always",
			},
			{
				ID:        "step_github_octocat_1_clone",
				Directory: "/home/github/octocat",
				Image:     "target/vela-git:v0.3.0",
				Name:      "clone",
				Number:    2,
				Pull:      "always",
			},
			{
				ID:        "step_github_octocat_1_echo",
				Commands:  []string{"echo hello"},
				Directory: "/home/github/octocat",
				Image:     "alpine:latest",
				Name:      "echo",
				Number:    3,
				Pull:      "always",
			},
		},
	}

	_repo = &library.Repo{
		ID:          vela.Int64(1),
		Org:         vela.String("github"),
		Name:        vela.String("octocat"),
		FullName:    vela.String("github/octocat"),
		Link:        vela.String("https://github.com/github/octocat"),
		Clone:       vela.String("https://github.com/github/octocat.git"),
		Branch:      vela.String("master"),
		Timeout:     vela.Int64(60),
		Visibility:  vela.String("public"),
		Private:     vela.Bool(false),
		Trusted:     vela.Bool(false),
		Active:      vela.Bool(true),
		AllowPull:   vela.Bool(false),
		AllowPush:   vela.Bool(true),
		AllowDeploy: vela.Bool(false),
		AllowTag:    vela.Bool(false),
	}

	_user = &library.User{
		ID:        vela.Int64(1),
		Name:      vela.String("octocat"),
		Token:     vela.String("superSecretToken"),
		Hash:      vela.String("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy"),
		Favorites: vela.Strings([]string{"github/octocat"}),
		Active:    vela.Bool(true),
		Admin:     vela.Bool(false),
	}
)
