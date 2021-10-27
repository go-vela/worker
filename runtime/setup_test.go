// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"testing"

	"github.com/go-vela/types/constants"
)

func TestRuntime_Setup_Docker(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver: constants.DriverDocker,
	}

	// run test
	_, err := _setup.Docker()
	if err != nil {
		t.Errorf("Docker returned err: %v", err)
	}
}

func TestRuntime_Setup_Kubernetes(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:     constants.DriverKubernetes,
		ConfigFile: "testdata/config",
		Namespace:  "docker",
	}

	// run test
	_, err := _setup.Kubernetes()
	if err != nil {
		t.Errorf("Kubernetes returned err: %v", err)
	}
}

func TestRuntime_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		failure bool
		setup   *Setup
		want    error
	}{
		{
			failure: false,
			setup: &Setup{
				Driver: constants.DriverDocker,
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:    constants.DriverKubernetes,
				Namespace: "docker",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver: "",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver: constants.DriverKubernetes,
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.setup.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
