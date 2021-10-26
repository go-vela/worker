// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"reflect"
	"sync"
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
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
	goodMap.Store(c.ID, new(library.Step))

	badMap := new(sync.Map)
	badMap.Store(c.ID, c)

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *library.Step
	}{
		{
			failure:   false,
			container: c,
			want:      new(library.Step),
			_map:      goodMap,
		},
		{
			failure:   true,
			container: c,
			want:      nil,
			_map:      badMap,
		},
		{
			failure:   true,
			container: new(pipeline.Container),
			want:      nil,
			_map:      new(sync.Map),
		},
		{
			failure:   true,
			container: nil,
			want:      nil,
			_map:      nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := Load(test.container, test._map)

		if test.failure {
			if err == nil {
				t.Errorf("Load should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Load returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Load is %v, want %v", got, test.want)
		}
	}
}

func TestStep_LoadInit(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
		want     *pipeline.Container
	}{
		{
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
			failure:  true,
			pipeline: nil,
			want:     nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := LoadInit(test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("LoadInit should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("LoadInit returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("LoadInit is %v, want %v", got, test.want)
		}
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
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *library.Log
	}{
		{
			failure:   false,
			container: c,
			want:      new(library.Log),
			_map:      goodMap,
		},
		{
			failure:   true,
			container: c,
			want:      nil,
			_map:      badMap,
		},
		{
			failure:   true,
			container: new(pipeline.Container),
			want:      nil,
			_map:      new(sync.Map),
		},
		{
			failure:   true,
			container: nil,
			want:      nil,
			_map:      nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := LoadLogs(test.container, test._map)

		if test.failure {
			if err == nil {
				t.Errorf("LoadLogs should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("LoadLogs returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("LoadLogs is %v, want %v", got, test.want)
		}
	}
}
