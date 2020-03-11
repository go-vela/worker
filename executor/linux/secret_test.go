// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/mock/server"

	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/google/go-cmp/cmp"
)

func TestExecutor_PullSecret_Success(t *testing.T) {
	// setup
	r, _ := docker.NewMock()

	// setup context
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())
	vela, _ := vela.NewClient(s.URL, nil)

	e, _ := New(vela, r)
	e.WithPipeline(&pipeline.Build{
		Version: "1",
		ID:      "__0",
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "__0_clone",
				Environment: map[string]string{},
				Image:       "target/vela-plugins/git:1",
				Name:        "clone",
				Number:      1,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "__0_exit",
				Environment: map[string]string{},
				Image:       "alpine:latest",
				Name:        "exit",
				Number:      2,
				Pull:        true,
				Ruleset: pipeline.Ruleset{
					Continue: true,
				},
				Commands: []string{"exit 1"},
			},
			&pipeline.Container{
				ID:          "__0_echo",
				Environment: map[string]string{},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      2,
				Pull:        true,
				Commands:    []string{"echo ${FOOBAR}"},
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "foobar",
						Target: "foobar",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "foo",
				Key:    "github/octocat/foo",
				Engine: "native",
				Type:   "repo",
			},
			&pipeline.Secret{
				Name:   "foo",
				Key:    "github/foo",
				Engine: "native",
				Type:   "org",
			},
			&pipeline.Secret{
				Name:   "foo",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "shared",
			},
		},
	})

	// run test
	got := e.PullSecret(context.Background())

	if got != nil {
		t.Errorf("PullSecret is %v, want nil", got)
	}
}

func TestLinux_Secret_injectSecret(t *testing.T) {
	// name and value of secret
	v := "foo"

	// setup types
	tests := []struct {
		step *pipeline.Container
		msec map[string]*library.Secret
		want *pipeline.Container
	}{
		// Tests for secrets with image ACLs
		{
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
		_ = injectSecrets(test.step, test.msec)
		got := test.step

		// Preferred use of reflect.DeepEqual(x, y interface) is giving false positives.
		// Switching to a Google library for increased clarity.
		// https://github.com/google/go-cmp
		if diff := cmp.Diff(test.want.Environment, got.Environment); diff != "" {
			t.Errorf("injectSecrets mismatch (-want +got):\n%s", diff)
		}
	}
}
