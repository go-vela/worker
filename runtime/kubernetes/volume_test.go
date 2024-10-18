// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestKubernetes_CreateVolume(t *testing.T) {
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

func TestKubernetes_InspectVolume(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
		pod      *v1.Pod
	}{
		{
			name:     "stages",
			failure:  false,
			pipeline: _stages,
			pod:      _pod,
		},
		{
			name:     "steps",
			failure:  false,
			pipeline: _steps,
			pod:      _pod,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock(test.pod)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

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

func TestKubernetes_RemoveVolume(t *testing.T) {
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
