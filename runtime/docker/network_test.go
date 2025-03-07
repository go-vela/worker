// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestDocker_CreateNetwork(t *testing.T) {
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
			err = _engine.CreateNetwork(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("CreateNetwork should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("CreateNetwork returned err: %v", err)
			}
		})
	}
}

func TestDocker_InspectNetwork(t *testing.T) {
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
			_, err = _engine.InspectNetwork(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("InspectNetwork should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectNetwork returned err: %v", err)
			}
		})
	}
}

func TestDocker_RemoveNetwork(t *testing.T) {
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
			err = _engine.RemoveNetwork(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("RemoveNetwork should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RemoveNetwork returned err: %v", err)
			}
		})
	}
}
