// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestKubernetes_CreateNetwork(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
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
			name:     "stages",
			failure:  false,
			pipeline: _stages,
		},
		{
			name:     "steps",
			failure:  false,
			pipeline: _steps,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := _engine.CreateNetwork(context.Background(), test.pipeline)

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

func TestKubernetes_InspectNetwork(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
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
			name:     "stages",
			failure:  false,
			pipeline: _stages,
		},
		{
			name:     "steps",
			failure:  false,
			pipeline: _steps,
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

func TestKubernetes_RemoveNetwork(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
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
			name:     "stages",
			failure:  false,
			pipeline: _stages,
		},
		{
			name:     "steps",
			failure:  false,
			pipeline: _steps,
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
