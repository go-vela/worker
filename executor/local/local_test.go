// SPDX-License-Identifier: Apache-2.0

package local

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/worker/runtime/docker"
)

func TestEqual(t *testing.T) {
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

	_local, err := New(
		WithBuild(testBuild()),
		WithHostname("localhost"),
		WithPipeline(testSteps()),
		WithRuntime(_runtime),
		WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create local executor: %v", err)
	}

	_alternate, err := New(
		WithBuild(testBuild()),
		WithHostname("a.different.host"),
		WithPipeline(testSteps()),
		WithRuntime(_runtime),
		WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create alternate local executor: %v", err)
	}

	tests := []struct {
		name string
		a    *client
		b    *client
		want bool
	}{
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "left nil",
			a:    nil,
			b:    _local,
			want: false,
		},
		{
			name: "right nil",
			a:    _local,
			b:    nil,
			want: false,
		},
		{
			name: "equal",
			a:    _local,
			b:    _local,
			want: true,
		},
		{
			name: "not equal",
			a:    _local,
			b:    _alternate,
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Equal(test.a, test.b); got != test.want {
				t.Errorf("Equal() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLocal_New(t *testing.T) {
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
		name     string
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			name:     "steps pipeline",
			failure:  false,
			pipeline: testSteps(),
		},
		{
			name:     "nil pipeline",
			failure:  true,
			pipeline: nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(
				WithBuild(testBuild()),
				WithHostname("localhost"),
				WithPipeline(test.pipeline),
				WithRuntime(_runtime),
				WithVelaClient(_client),
			)

			if test.failure {
				if err == nil {
					t.Errorf("New should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("New returned err: %v", err)
			}
		})
	}
}

// testBuild is a test helper function to create a Build
// type with all fields set to a fake value.
func testBuild() *api.Build {
	return &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         testRepo(),
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
		Distribution: new("Local"),
	}
}

// testRepo is a test helper function to create a Repo
// type with all fields set to a fake value.
func testRepo() *api.Repo {
	return &api.Repo{
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
		Owner:      testUser(),
	}
}

// testUser is a test helper function to create a User
// type with all fields set to a fake value.
func testUser() *api.User {
	return &api.User{
		ID:        new(int64(1)),
		Name:      new("octocat"),
		Token:     new("superSecretToken"),
		Favorites: new([]string{"github/octocat"}),
		Active:    new(true),
		Admin:     new(false),
	}
}

// testSteps is a test helper function to create a steps
// pipeline with fake steps.
func testSteps() *pipeline.Build {
	return &pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service_github_octocat_1_postgres",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		Steps: pipeline.ContainerSlice{
			{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        constants.InitName,
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo hello"},
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
			},
		},
		Secrets: pipeline.SecretSlice{
			{
				Name:   "foo",
				Key:    "github/octocat/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			{
				Name:   "foo",
				Key:    "github/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
			{
				Name:   "foo",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
	}
}
