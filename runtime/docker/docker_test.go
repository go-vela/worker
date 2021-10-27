// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"testing"

	"github.com/go-vela/types/pipeline"

	"gotest.tools/v3/env"
)

func TestDocker_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		envs    map[string]string
	}{
		{
			failure: false,
			envs:    map[string]string{},
		},
		{
			failure: true,
			envs: map[string]string{
				"DOCKER_CERT_PATH": "invalid/path",
			},
		},
	}

	// defer env cleanup
	defer env.PatchAll(t, nil)()

	// run tests
	for _, test := range tests {
		// patch environment for tests
		env.PatchAll(t, test.envs)

		_, err := New(
			WithPrivilegedImages([]string{"alpine"}),
			WithHostVolumes([]string{"/foo/bar.txt:/foo/bar.txt"}),
		)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}

// setup global variables used for testing.
var (
	_container = &pipeline.Container{
		ID:          "step_github_octocat_1_clone",
		Directory:   "/vela/src/github.com/octocat/helloworld",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "target/vela-git:v0.4.0",
		Name:        "clone",
		Number:      2,
		Pull:        "always",
	}

	_pipeline = &pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service_github_octocat_1_postgres",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
			},
		},
		Steps: pipeline.ContainerSlice{
			{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
			},
		},
	}
)
