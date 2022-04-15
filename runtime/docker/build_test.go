// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_InspectBuild(t *testing.T) {
	// setup Docker
	_engine, err := NewMock()
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
			name:     "steps",
			failure:  false,
			pipeline: _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		_, err = _engine.InspectBuild(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("InspectBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectBuild returned err: %v", err)
		}
	}
}

func TestDocker_SetupBuild(t *testing.T) {
	// setup Docker
	_engine, err := NewMock()
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
			name:     "steps",
			failure:  false,
			pipeline: _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.SetupBuild(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("SetupBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("SetupBuild returned err: %v", err)
		}
	}
}

func TestDocker_AssembleBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			name:     "steps",
			failure:  false,
			pipeline: _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := NewMock()
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		err = _engine.AssembleBuild(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("AssembleBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("AssembleBuild returned err: %v", err)
		}
	}
}

func TestDocker_RemoveBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			name:     "steps",
			failure:  false,
			pipeline: _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := NewMock()
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		err = _engine.RemoveBuild(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("RemoveBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RemoveBuild returned err: %v", err)
		}
	}
}
