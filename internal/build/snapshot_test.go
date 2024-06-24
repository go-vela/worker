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
		ID:         vela.Int64(1),
		Org:        vela.String("github"),
		Name:       vela.String("octocat"),
		FullName:   vela.String("github/octocat"),
		Link:       vela.String("https://github.com/github/octocat"),
		Clone:      vela.String("https://github.com/github/octocat.git"),
		Branch:     vela.String("main"),
		Timeout:    vela.Int64(60),
		Visibility: vela.String("public"),
		Private:    vela.Bool(false),
		Trusted:    vela.Bool(false),
		Active:     vela.Bool(true),
		AllowEvents: &api.Events{
			Push: &actions.Push{
				Branch: vela.Bool(true),
			},
		},
	}

	b := &api.Build{
		ID:           vela.Int64(1),
		Repo:         r,
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
		Branch:       vela.String("main"),
		Ref:          vela.String("refs/heads/main"),
		BaseRef:      vela.String(""),
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
			Snapshot(test.build, test.client, test.err, nil)
		})
	}
}
