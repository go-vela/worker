// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"reflect"
	"sync"
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
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
	goodMap.Store(c.ID, new(library.Service))

	badMap := new(sync.Map)
	badMap.Store(c.ID, c)

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
		_map      *sync.Map
		want      *library.Service
	}{
		{
			failure:   false,
			container: c,
			want:      new(library.Service),
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
