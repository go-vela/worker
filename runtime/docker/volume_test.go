// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_CreateVolume(t *testing.T) {
	// setup types
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
		{
			name:     "empty",
			failure:  true,
			pipeline: new(pipeline.Build),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.CreateVolume(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("CreateVolume should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("CreateVolume returned err: %v", err)
			}
		})
	}
}

func TestDocker_InspectVolume(t *testing.T) {
	// setup types
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
		{
			name:     "empty",
			failure:  true,
			pipeline: new(pipeline.Build),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.InspectVolume(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("InspectVolume should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectVolume returned err: %v", err)
			}
		})
	}
}

func TestDocker_RemoveVolume(t *testing.T) {
	// setup types
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
		{
			name:     "empty",
			failure:  true,
			pipeline: new(pipeline.Build),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.RemoveVolume(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("RemoveVolume should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RemoveVolume returned err: %v", err)
			}
		})
	}
}
