// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestLinux_GetBuild(t *testing.T) {
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
		name    string
		failure bool
		engine  *client
	}{
		{
			name:    "with build",
			failure: false,
			engine:  _engine,
		},
		{
			name:    "missing build",
			failure: true,
			engine:  new(client),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.engine.GetBuild()

			if test.failure {
				if err == nil {
					t.Errorf("GetBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("GetBuild returned err: %v", err)
			}

			if !reflect.DeepEqual(got, _build) {
				t.Errorf("GetBuild is %v, want %v", got, _build)
			}
		})
	}
}

func TestLinux_GetPipeline(t *testing.T) {
	// setup types
	_steps := testSteps(constants.DriverDocker)

	_engine, err := New(
		WithPipeline(_steps),
	)
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		engine  *client
	}{
		{
			name:    "with pipeline",
			failure: false,
			engine:  _engine,
		},
		{
			name:    "missing pipeline",
			failure: true,
			engine:  new(client),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.engine.GetPipeline()

			if test.failure {
				if err == nil {
					t.Errorf("GetPipeline should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("GetPipeline returned err: %v", err)
			}

			if !reflect.DeepEqual(got, _steps) {
				t.Errorf("GetPipeline is %v, want %v", got, _steps)
			}
		})
	}
}
