// SPDX-License-Identifier: Apache-2.0

package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/mock/server"
)

func TestService_Snapshot(t *testing.T) {
	// setup types
	_repo := &api.Repo{
		ID:          vela.Int64(1),
		Org:         vela.String("github"),
		Name:        vela.String("octocat"),
		FullName:    vela.String("github/octocat"),
		Link:        vela.String("https://github.com/github/octocat"),
		Clone:       vela.String("https://github.com/github/octocat.git"),
		Branch:      vela.String("main"),
		Timeout:     vela.Int32(60),
		Visibility:  vela.String("public"),
		Private:     vela.Bool(false),
		Trusted:     vela.Bool(false),
		Active:      vela.Bool(true),
		AllowEvents: api.NewEventsFromMask(1),
	}

	_build := &api.Build{
		ID:           vela.Int64(1),
		Number:       vela.Int64(1),
		Repo:         _repo,
		Parent:       vela.Int64(1),
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
		Branch:       vela.String("main"),
		Ref:          vela.String("refs/heads/main"),
		BaseRef:      vela.String(""),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("linux"),
	}

	_container := &pipeline.Container{
		ID:          "service_github_octocat_1_postgres",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "postgres:12-alpine",
		Name:        "postgres",
		Number:      1,
		Ports:       []string{"5432:5432"},
		Pull:        "not_present",
	}

	_exitCode := &pipeline.Container{
		ID:          "service_github_octocat_1_postgres",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		ExitCode:    137,
		Image:       "postgres:12-alpine",
		Name:        "postgres",
		Number:      1,
		Ports:       []string{"5432:5432"},
		Pull:        "not_present",
	}

	_service := &api.Service{
		ID:           vela.Int64(1),
		BuildID:      vela.Int64(1),
		RepoID:       vela.Int64(1),
		Number:       vela.Int32(1),
		Name:         vela.String("postgres"),
		Image:        vela.String("postgres:12-alpine"),
		Status:       vela.String("running"),
		ExitCode:     vela.Int32(0),
		Created:      vela.Int64(1563474076),
		Started:      vela.Int64(0),
		Finished:     vela.Int64(1563474079),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("linux"),
	}

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
			Snapshot(test.container, test.build, test.client, nil, test.service)
		})
	}
}
