// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestKubernetes_CreateNetwork(t *testing.T) {
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
		err := _engine.CreateNetwork(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("CreateNetwork should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateNetwork returned err: %v", err)
		}
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
		_, err = _engine.InspectNetwork(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("InspectNetwork should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectNetwork returned err: %v", err)
		}
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
		err = _engine.RemoveNetwork(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("RemoveNetwork should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RemoveNetwork returned err: %v", err)
		}
	}
}
