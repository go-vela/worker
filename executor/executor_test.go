// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/worker/executor/linux"
	"github.com/go-vela/worker/executor/local"
	"github.com/go-vela/worker/runtime/docker"
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
		linux.WithRuntime(_runtime),
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
		local.WithRuntime(_runtime),
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
		equal   any
	}{
		{
			name:    "driver-darwin",
			failure: true,
			setup: &Setup{
				Build:    _build,
				Client:   _client,
				Driver:   constants.DriverDarwin,
				Pipeline: _pipeline,
				Runtime:  _runtime,
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
				Runtime:    _runtime,
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
				Runtime:  _runtime,
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
				Runtime:  _runtime,
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
				Runtime:  _runtime,
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
				Runtime:  _runtime,
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
	_pipeline = &pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Steps: pipeline.ContainerSlice{
			{
				ID:        "step_github_octocat_1_init",
				Directory: "/home/github/octocat",
				Image:     "#init",
				Name:      constants.InitName,
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

	_user = &api.User{
		ID:        new(int64(1)),
		Name:      new("octocat"),
		Token:     new("superSecretToken"),
		Favorites: new([]string{"github/octocat"}),
		Active:    new(true),
		Admin:     new(false),
	}

	_allowEvents = &api.Events{
		Push: &actions.Push{
			Branch: new(true),
			Tag:    new(true),
		},
		PullRequest: &actions.Pull{
			Opened:      new(true),
			Synchronize: new(true),
			Edited:      new(true),
			Reopened:    new(true),
			Labeled:     new(true),
			Unlabeled:   new(true),
		},
		Comment: &actions.Comment{
			Created: new(true),
			Edited:  new(true),
		},
		Deployment: &actions.Deploy{
			Created: new(true),
		},
	}

	_repo = &api.Repo{
		ID:          new(int64(1)),
		Org:         new("github"),
		Name:        new("octocat"),
		FullName:    new("github/octocat"),
		Link:        new("https://github.com/github/octocat"),
		Clone:       new("https://github.com/github/octocat.git"),
		Branch:      new("main"),
		Timeout:     new(int32(60)),
		Visibility:  new("public"),
		Private:     new(false),
		Trusted:     new(false),
		Active:      new(true),
		AllowEvents: _allowEvents,
		Owner:       _user,
	}

	_build = &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("push"),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}
)
