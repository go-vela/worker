// SPDX-License-Identifier: Apache-2.0

package service

import (
	"reflect"
	"sync"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestService_Load(t *testing.T) {
	// setup types
	c := &pipeline.Container{
		ID:          "service_github_octocat_1_postgres",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "postgres:12-alpine",
		Name:        "postgres",
		Number:      1,
		Ports:       []string{"5432:5432"},
		Pull:        "not_present",
	}

	goodMap := new(sync.Map)
	goodMap.Store(c.ID, new(api.Service))

	badMap := new(sync.Map)
	badMap.Store(c.ID, c)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *api.Service
	}{
		{
			name:      "good map",
			failure:   false,
			container: c,
			want:      new(api.Service),
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

func TestStep_LoadLogs(t *testing.T) {
	// setup types
	c := &pipeline.Container{
		ID:          "service_github_octocat_1_postgres",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "postgres:12-alpine",
		Name:        "postgres",
		Number:      1,
		Ports:       []string{"5432:5432"},
		Pull:        "not_present",
	}

	goodMap := new(sync.Map)
	goodMap.Store(c.ID, new(api.Log))

	badMap := new(sync.Map)
	badMap.Store(c.ID, c)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *api.Log
	}{
		{
			name:      "good map",
			failure:   false,
			container: c,
			want:      new(api.Log),
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
