// SPDX-License-Identifier: Apache-2.0

package step

import (
	"reflect"
	"sync"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/types/library"
)

func TestStep_Load(t *testing.T) {
	// setup types
	c := &pipeline.Container{
		ID:          "step_github_octocat_1_init",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "#init",
		Name:        "init",
		Number:      1,
		Pull:        "always",
	}

	goodMap := new(sync.Map)
	goodMap.Store(c.ID, new(api.Step))

	badMap := new(sync.Map)
	badMap.Store(c.ID, c)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *api.Step
	}{
		{
			name:      "good map",
			failure:   false,
			container: c,
			want:      new(api.Step),
			_map:      goodMap,
		},
		{
			name:      "bad map",
			failure:   true,
			container: c,
			want:      nil,
			_map:      badMap,
		},
		{
			name:      "empty map",
			failure:   true,
			container: new(pipeline.Container),
			want:      nil,
			_map:      new(sync.Map),
		},
		{
			name:      "nil map",
			failure:   true,
			container: nil,
			want:      nil,
			_map:      nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Load(test.container, test._map)

			if test.failure {
				if err == nil {
					t.Errorf("Load should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("Load returned err: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Load is %v, want %v", got, test.want)
			}
		})
	}
}

func TestStep_LoadInit(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
		want     *pipeline.Container
	}{
		{
			name:    "stages",
			failure: false,
			pipeline: &pipeline.Build{
				Version: "1",
				ID:      "github_octocat_1",
				Stages: pipeline.StageSlice{
					{
						Name: "init",
						Steps: pipeline.ContainerSlice{
							{
								ID:          "github_octocat_1_init_init",
								Directory:   "/vela/src/github.com/github/octocat",
								Environment: map[string]string{"FOO": "bar"},
								Image:       "#init",
								Name:        "init",
								Number:      1,
								Pull:        "always",
							},
						},
					},
				},
			},
			want: &pipeline.Container{
				ID:          "github_octocat_1_init_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
		},
		{
			name:    "steps",
			failure: false,
			pipeline: &pipeline.Build{
				Version: "1",
				ID:      "github_octocat_1",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step_github_octocat_1_init",
						Directory:   "/vela/src/github.com/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
				},
			},
			want: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
		},
		{
			name:     "nil failure",
			failure:  true,
			pipeline: nil,
			want:     nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := LoadInit(test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("LoadInit should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("LoadInit returned err: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LoadInit is %v, want %v", got, test.want)
			}
		})
	}
}

func TestStep_LoadLogs(t *testing.T) {
	// setup types
	c := &pipeline.Container{
		ID:          "step_github_octocat_1_init",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "#init",
		Name:        "init",
		Number:      1,
		Pull:        "always",
	}

	goodMap := new(sync.Map)
	goodMap.Store(c.ID, new(library.Log))

	badMap := new(sync.Map)
	badMap.Store(c.ID, c)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *library.Log
	}{
		{
			name:      "good map",
			failure:   false,
			container: c,
			want:      new(library.Log),
			_map:      goodMap,
		},
		{
			name:      "bad map",
			failure:   true,
			container: c,
			want:      nil,
			_map:      badMap,
		},
		{
			name:      "empty map",
			failure:   true,
			container: new(pipeline.Container),
			want:      nil,
			_map:      new(sync.Map),
		},
		{
			name:      "nil map",
			failure:   true,
			container: nil,
			want:      nil,
			_map:      nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := LoadLogs(test.container, test._map)

			if test.failure {
				if err == nil {
					t.Errorf("LoadLogs should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("LoadLogs returned err: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LoadLogs is %v, want %v", got, test.want)
			}
		})
	}
}
