// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/mock/server"

	"github.com/go-vela/worker/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/google/go-cmp/cmp"
)

func TestLinux_Secret_create(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_steps := testSteps()

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
		name      string
		failure   bool
		container *pipeline.Container
	}{
		{
			name:    "good image tag",
			failure: false,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "notfound image tag",
			failure: true,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:notfound",
				Name:        "vault",
				Number:      1,
				Pull:        "not_present",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_steps),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.secret.create(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("create should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("create returned err: %v", err)
			}
		})
	}
}

func TestLinux_Secret_delete(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_steps := testSteps()

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

	_step := new(library.Step)
	_step.SetName("clone")
	_step.SetNumber(2)
	_step.SetStatus(constants.StatusPending)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		step      *library.Step
	}{
		{
			name:    "running container-empty step",
			failure: false,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "always",
			},
			step: new(library.Step),
		},
		{
			name:    "running container-pending step",
			failure: false,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      2,
				Pull:        "always",
			},
			step: _step,
		},
		{
			name:    "inspecting container failure due to invalid container id",
			failure: true,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_notfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "notfound",
				Number:      2,
				Pull:        "always",
			},
			step: new(library.Step),
		},
		{
			name:    "removing container failure",
			failure: true,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "ignorenotfound",
				Number:      2,
				Pull:        "always",
			},
			step: new(library.Step),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_steps),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			_ = _engine.CreateBuild(context.Background())

			_engine.steps.Store(test.container.ID, test.step)

			err = _engine.secret.destroy(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("destroy should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("destroy returned err: %v", err)
			}
		})
	}
}

func TestLinux_Secret_exec(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

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
		pipeline string
	}{
		{
			name:     "basic secrets pipeline",
			failure:  false,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "pipeline with secret name not found",
			failure:  true,
			pipeline: "testdata/build/secrets/name_notfound.yml",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, _ := ioutil.ReadFile(test.pipeline)

			p, _ := compiler.
				WithBuild(_build).
				WithRepo(_repo).
				WithUser(_user).
				WithMetadata(_metadata).
				Compile(file)

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(p),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			_engine.build.SetStatus(constants.StatusSuccess)

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			err = _engine.secret.exec(context.Background(), &p.Secrets)

			if test.failure {
				if err == nil {
					t.Errorf("exec should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("exec returned err: %v", err)
			}
		})
	}
}

func TestLinux_Secret_pull(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

	gin.SetMode(gin.TestMode)

	server := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(server.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		secret  *pipeline.Secret
	}{
		{
			name:    "success with org secret",
			failure: false,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "failure with invalid org secret",
			failure: true,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo/foo/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "failure with org secret key not found",
			failure: true,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "success with repo secret",
			failure: false,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octocat/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "failure with invalid repo secret",
			failure: true,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo/foo/foo/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "failure with repo secret key not found",
			failure: true,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "success with shared secret",
			failure: false,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "failure with shared secret key not found",
			failure: true,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "failure with invalid type",
			failure: true,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "invalid",
				Origin: &pipeline.Container{},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(testSteps()),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			_, err = _engine.secret.pull(test.secret)

			if test.failure {
				if err == nil {
					t.Errorf("pull should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("pull returned err: %v", err)
			}
		})
	}
}

func TestLinux_Secret_stream(t *testing.T) {
	// setup types
	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_steps := testSteps()

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
		name      string
		failure   bool
		logs      *library.Log
		container *pipeline.Container
	}{
		{
			name:    "container step succeeds",
			failure: false,
			logs:    new(library.Log),
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
		},
		{
			name:    "container step fails because of invalid container id",
			failure: true,
			logs:    new(library.Log),
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_notfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "notfound",
				Number:      2,
				Pull:        "always",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_steps),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			err = _engine.secret.stream(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("stream should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("stream returned err: %v", err)
			}
		})
	}
}

func TestLinux_Secret_injectSecret(t *testing.T) {
	// name and value of secret
	v := "foo"

	// setup types
	tests := []struct {
		name string
		step *pipeline.Container
		msec map[string]*library.Secret
		want *pipeline.Container
	}{
		// Tests for secrets with image ACLs
		{
			name: "secret with empty image ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Images: &[]string{""}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
			},
		},
		{
			name: "secret with matching image ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Images: &[]string{"alpine"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo"},
			},
		},
		{
			name: "secret with matching image:tag ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Images: &[]string{"alpine:latest"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo"},
			},
		},
		{
			name: "secret with non-matching image ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Images: &[]string{"centos"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
			},
		},

		// Tests for secrets with event ACLs
		{ // push event checks
			name: "secret with matching push event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"push"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with non-matching push event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"deployment"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
			},
		},
		{ // pull_request event checks
			name: "secret with matching pull_request event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "pull_request"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"pull_request"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "BUILD_EVENT": "pull_request"},
			},
		},
		{
			name: "secret with non-matching pull_request event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "pull_request"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"deployment"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "pull_request"},
			},
		},
		{ // tag event checks
			name: "secret with matching tag event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "tag"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"tag"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "BUILD_EVENT": "tag"},
			},
		},
		{
			name: "secret with non-matching tag event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "tag"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"deployment"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "tag"},
			},
		},
		{ // deployment event checks
			name: "secret with matching deployment event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "deployment"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"deployment"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "BUILD_EVENT": "deployment"},
			},
		},
		{
			name: "secret with non-matching deployment event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "deployment"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"tag"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "deployment"},
			},
		},

		// Tests for secrets with event and image ACLs
		{
			name: "secret with matching event ACL and non-matching image ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"push"}, Images: &[]string{"centos"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with non-matching event ACL and matching image ACL not injected",
			step: &pipeline.Container{
				Image:       "centos:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"pull_request"}, Images: &[]string{"centos"}}},
			want: &pipeline.Container{
				Image:       "centos:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with matching event ACL and matching image ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*library.Secret{"FOO": {Name: &v, Value: &v, Events: &[]string{"push"}, Images: &[]string{"alpine"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "BUILD_EVENT": "push"},
			},
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = injectSecrets(test.step, test.msec)
			got := test.step

			// Preferred use of reflect.DeepEqual(x, y interface) is giving false positives.
			// Switching to a Google library for increased clarity.
			// https://github.com/google/go-cmp
			if diff := cmp.Diff(test.want.Environment, got.Environment); diff != "" {
				t.Errorf("injectSecrets mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLinux_Secret_escapeNewlineSecrets(t *testing.T) {
	// name and value of secret
	n := "foo"
	v := "bar\\nbaz"
	vEscaped := "bar\\\nbaz"

	// desired secret value
	w := "bar\\\nbaz"

	// setup types
	tests := []struct {
		name      string
		secretMap map[string]*library.Secret
		want      map[string]*library.Secret
	}{
		{
			name:      "not escaped",
			secretMap: map[string]*library.Secret{"FOO": {Name: &n, Value: &v}},
			want:      map[string]*library.Secret{"FOO": {Name: &n, Value: &w}},
		},
		{
			name:      "already escaped",
			secretMap: map[string]*library.Secret{"FOO": {Name: &n, Value: &vEscaped}},
			want:      map[string]*library.Secret{"FOO": {Name: &n, Value: &w}},
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			escapeNewlineSecrets(test.secretMap)
			got := test.secretMap

			// Preferred use of reflect.DeepEqual(x, y interface) is giving false positives.
			// Switching to a Google library for increased clarity.
			// https://github.com/google/go-cmp
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("escapeNewlineSecrets mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
