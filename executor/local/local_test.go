// SPDX-License-Identifier: Apache-2.0

package local

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/mock/server"

	"github.com/go-vela/worker/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
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
		WithRepo(testRepo()),
		WithRuntime(_runtime),
		WithUser(testUser()),
		WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create local executor: %v", err)
	}

	_alternate, err := New(
		WithBuild(testBuild()),
		WithHostname("a.different.host"),
		WithPipeline(testSteps()),
		WithRepo(testRepo()),
		WithRuntime(_runtime),
		WithUser(testUser()),
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
				WithRepo(testRepo()),
				WithRuntime(_runtime),
				WithUser(testUser()),
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
func testBuild() *library.Build {
	return &library.Build{
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
		Branch:       vela.String("main"),
		Ref:          vela.String("refs/heads/main"),
		BaseRef:      vela.String(""),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("Local"),
	}
}

// testRepo is a test helper function to create a Repo
// type with all fields set to a fake value.
func testRepo() *api.Repo {
	return &api.Repo{
		ID:          vela.Int64(1),
		Org:         vela.String("github"),
		Name:        vela.String("octocat"),
		FullName:    vela.String("github/octocat"),
		Link:        vela.String("https://github.com/github/octocat"),
		Clone:       vela.String("https://github.com/github/octocat.git"),
		Branch:      vela.String("main"),
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
}

// testUser is a test helper function to create a User
// type with all fields set to a fake value.
func testUser() *library.User {
	return &library.User{
		ID:        vela.Int64(1),
		Name:      vela.String("octocat"),
		Token:     vela.String("superSecretToken"),
		Hash:      vela.String("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy"),
		Favorites: vela.Strings([]string{"github/octocat"}),
		Active:    vela.Bool(true),
		Admin:     vela.Bool(false),
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
				Name:        "init",
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
