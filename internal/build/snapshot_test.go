// SPDX-License-Identifier: Apache-2.0

package build

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/mock/server"
)

func TestBuild_Snapshot(t *testing.T) {
	// setup types
	r := &api.Repo{
		ID:         new(int64(1)),
		Org:        new("github"),
		Name:       new("octocat"),
		FullName:   new("github/octocat"),
		Link:       new("https://github.com/github/octocat"),
		Clone:      new("https://github.com/github/octocat.git"),
		Branch:     new("main"),
		Timeout:    new(int32(60)),
		Visibility: new("public"),
		Private:    new(false),
		Trusted:    new(false),
		Active:     new(true),
		AllowEvents: &api.Events{
			Push: &actions.Push{
				Branch: new(true),
			},
		},
	}

	b := &api.Build{
		ID:           new(int64(1)),
		Repo:         r,
		Number:       new(int64(1)),
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

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		name   string
		build  *api.Build
		client *vela.Client
		err    error
	}{
		{
			name:   "build with error",
			build:  b,
			client: _client,
			err:    errors.New("unable to create network"),
		},
		{
			name:   "nil build with error",
			build:  nil,
			client: _client,
			err:    errors.New("unable to create network"),
		},
		{
			name:   "nil everything",
			build:  nil,
			client: nil,
			err:    nil,
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			Snapshot(t.Context(), test.build, test.client, test.err, nil)
		})
	}
}
