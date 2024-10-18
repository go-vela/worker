// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
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
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.InspectBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("InspectBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectBuild returned err: %v", err)
			}
		})
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
		t.Run(test.name, func(t *testing.T) {
			err = _engine.SetupBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("SetupBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("SetupBuild returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_StreamBuild(t *testing.T) {
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock()
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			err = _engine.StreamBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("StreamBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("StreamBuild returned err: %v", err)
			}
		})
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
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock()
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			err = _engine.AssembleBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("AssembleBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("AssembleBuild returned err: %v", err)
			}
		})
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
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock()
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			err = _engine.RemoveBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("RemoveBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RemoveBuild returned err: %v", err)
			}
		})
	}
}
