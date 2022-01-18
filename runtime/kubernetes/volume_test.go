// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"

	v1 "k8s.io/api/core/v1"
)

func TestKubernetes_CreateVolume(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			failure:  false,
			pipeline: _stages,
		},
		{
			failure:  false,
			pipeline: _steps,
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.CreateVolume(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("CreateVolume should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateVolume returned err: %v", err)
		}
	}
}

func TestKubernetes_InspectVolume(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
		pod      *v1.Pod
	}{
		{
			failure:  false,
			pipeline: _stages,
			pod:      _pod,
		},
		{
			failure:  false,
			pipeline: _steps,
			pod:      _pod,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := NewMock(test.pod)
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		_, err = _engine.InspectVolume(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("InspectVolume should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectVolume returned err: %v", err)
		}
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
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			failure:  false,
			pipeline: _stages,
		},
		{
			failure:  false,
			pipeline: _steps,
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.RemoveVolume(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("RemoveVolume should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RemoveVolume returned err: %v", err)
		}
	}
}
