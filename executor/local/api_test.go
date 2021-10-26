// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"reflect"
	"testing"
)

func TestLocal_GetBuild(t *testing.T) {
	// setup types
	_build := testBuild()

	_engine, err := New(
		WithBuild(_build),
	)
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		engine  *client
	}{
		{
			failure: false,
			engine:  _engine,
		},
		{
			failure: true,
			engine:  new(client),
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.engine.GetBuild()

		if test.failure {
			if err == nil {
				t.Errorf("GetBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(got, _build) {
			t.Errorf("GetBuild is %v, want %v", got, _build)
		}
	}
}

func TestLocal_GetPipeline(t *testing.T) {
	// setup types
	_steps := testSteps()

	_engine, err := New(
		WithPipeline(_steps),
	)
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		engine  *client
	}{
		{
			failure: false,
			engine:  _engine,
		},
		{
			failure: true,
			engine:  new(client),
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.engine.GetPipeline()

		if test.failure {
			if err == nil {
				t.Errorf("GetPipeline should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetPipeline returned err: %v", err)
		}

		if !reflect.DeepEqual(got, _steps) {
			t.Errorf("GetPipeline is %v, want %v", got, _steps)
		}
	}
}

func TestLocal_GetRepo(t *testing.T) {
	// setup types
	_repo := testRepo()

	_engine, err := New(
		WithRepo(_repo),
	)
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		engine  *client
	}{
		{
			failure: false,
			engine:  _engine,
		},
		{
			failure: true,
			engine:  new(client),
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.engine.GetRepo()

		if test.failure {
			if err == nil {
				t.Errorf("GetRepo should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepo returned err: %v", err)
		}

		if !reflect.DeepEqual(got, _repo) {
			t.Errorf("GetRepo is %v, want %v", got, _repo)
		}
	}
}
