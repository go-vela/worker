// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/mock/server"
	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/library"
)

func TestBuild_Upload(t *testing.T) {
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

	_canceled := *_build
	_canceled.SetStatus("canceled")

	_error := *_build
	_error.SetStatus("error")

	_pending := *_build
	_pending.SetStatus("pending")

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

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		build  *library.Build
		client *vela.Client
		err    error
		repo   *library.Repo
	}{
		{
			build:  _build,
			client: _client,
			err:    errors.New("unable to create network"),
			repo:   _repo,
		},
		{
			build:  &_canceled,
			client: _client,
			err:    errors.New("unable to create network"),
			repo:   _repo,
		},
		{
			build:  &_error,
			client: _client,
			err:    errors.New("unable to create network"),
			repo:   _repo,
		},
		{
			build:  &_pending,
			client: _client,
			err:    errors.New("unable to create network"),
			repo:   _repo,
		},
		{
			build:  nil,
			client: _client,
			err:    errors.New("unable to create network"),
			repo:   _repo,
		},
		{
			build:  nil,
			client: nil,
			err:    nil,
			repo:   nil,
		},
	}

	// run test
	for _, test := range tests {
		Upload(test.build, test.client, test.err, nil, test.repo)
	}
}
