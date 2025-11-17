// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"testing"

	"gotest.tools/v3/env"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

func TestDocker_New(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		envs    map[string]string
	}{
		{
			name:    "default",
			failure: false,
			envs:    map[string]string{},
		},
		{
			name:    "with invalid DOCKER_CERT_PATH",
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
		t.Run(test.name, func(t *testing.T) {
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

				return // continue to next test
			}

			if err != nil {
				t.Errorf("New returned err: %v", err)
			}
		})
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
				Name:        constants.InitName,
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
