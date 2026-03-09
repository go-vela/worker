// SPDX-License-Identifier: Apache-2.0

package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/mock/server"
)

func TestService_Upload(t *testing.T) {
	// setup types
	_repo := &api.Repo{
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
		AllowEvents: api.NewEventsFromMask(1),
	}

	_build := &api.Build{
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

	_container := &pipeline.Container{
		ID:          "step_github_octocat_1_init",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "#init",
		Name:        constants.InitName,
		Number:      1,
		Pull:        "always",
	}

	_exitCode := &pipeline.Container{
		ID:          "step_github_octocat_1_init",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		ExitCode:    137,
		Image:       "#init",
		Name:        constants.InitName,
		Number:      1,
		Pull:        "always",
	}

	_service := &api.Service{
		ID:           new(int64(1)),
		BuildID:      new(int64(1)),
		RepoID:       new(int64(1)),
		Number:       new(int32(1)),
		Name:         new("postgres"),
		Image:        new("postgres:12-alpine"),
		Status:       new("running"),
		ExitCode:     new(int32(0)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(0)),
		Finished:     new(int64(1563474079)),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_canceled := *_service
	_canceled.SetStatus("canceled")

	_error := *_service
	_error.SetStatus("error")

	_pending := *_service
	_pending.SetStatus("pending")

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		name      string
		build     *api.Build
		client    *vela.Client
		container *pipeline.Container
		service   *api.Service
	}{
		{
			name:      "running service",
			build:     _build,
			client:    _client,
			container: _container,
			service:   _service,
		},
		{
			name:      "canceled service",
			build:     _build,
			client:    _client,
			container: _container,
			service:   &_canceled,
		},
		{
			name:      "errored service",
			build:     _build,
			client:    _client,
			container: _container,
			service:   &_error,
		},
		{
			name:      "pending service",
			build:     _build,
			client:    _client,
			container: _container,
			service:   &_pending,
		},
		{
			name:      "exited service",
			build:     _build,
			client:    _client,
			container: _exitCode,
			service:   nil,
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			Upload(t.Context(), test.container, test.build, test.client, nil, test.service)
		})
	}
}
