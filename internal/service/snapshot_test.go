// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/mock/server"
	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func TestService_Snapshot(t *testing.T) {
	// setup types
	_build := &library.Build{
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

	_repo := &library.Repo{
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

	_service := &library.Service{
		ID:           vela.Int64(1),
		BuildID:      vela.Int64(1),
		RepoID:       vela.Int64(1),
		Number:       vela.Int(1),
		Name:         vela.String("postgres"),
		Image:        vela.String("postgres:12-alpine"),
		Status:       vela.String("running"),
		ExitCode:     vela.Int(0),
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
		build     *library.Build
		client    *vela.Client
		container *pipeline.Container
		repo      *library.Repo
		service   *library.Service
	}{
		{
			build:     _build,
			client:    _client,
			container: _container,
			repo:      _repo,
			service:   _service,
		},
		{
			build:     _build,
			client:    _client,
			container: _exitCode,
			repo:      _repo,
			service:   nil,
		},
	}

	// run test
	for _, test := range tests {
		Snapshot(test.container, test.build, test.client, nil, test.repo, test.service)
	}
}
