// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"testing"

	"github.com/go-vela/server/constants"
)

func TestRuntime_Setup_Docker(t *testing.T) {
	tests := []struct {
		name string
		mock bool
	}{
		{name: "standard", mock: false},
		{name: "mocked", mock: true},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			_setup := &Setup{
				Mock:   test.mock,
				Driver: constants.DriverDocker,
			}

			_, err := _setup.Docker()
			if err != nil {
				t.Errorf("Docker returned err: %v", err)
			}
		})
	}
}

func TestRuntime_Setup_Kubernetes(t *testing.T) {
	tests := []struct {
		name string
		mock bool
	}{
		{name: "standard", mock: false},
		{name: "mocked", mock: true},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			_setup := &Setup{
				Mock:       test.mock,
				Driver:     constants.DriverKubernetes,
				ConfigFile: "testdata/config",
				Namespace:  "docker",
			}

			_, err := _setup.Kubernetes()
			if err != nil {
				t.Errorf("Kubernetes returned err: %v", err)
			}
		})
	}
}

func TestRuntime_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		name    string
		failure bool
		setup   *Setup
		want    error
	}{
		{
			name:    "docker driver",
			failure: false,
			setup: &Setup{
				Driver:           constants.DriverDocker,
				DropCapabilities: []string{"CAP_DAC_OVERRIDE"},
			},
		},
		{
			name:    "docker driver bad cap",
			failure: true,
			setup: &Setup{
				Driver:           constants.DriverDocker,
				DropCapabilities: []string{"BAD"},
			},
		},
		{
			name:    "kubernetes driver",
			failure: false,
			setup: &Setup{
				Driver:    constants.DriverKubernetes,
				Namespace: "docker",
			},
		},
		{
			name:    "empty driver",
			failure: true,
			setup: &Setup{
				Driver: "",
			},
		},
		{
			name:    "kubernetes driver-missing namespace",
			failure: true,
			setup: &Setup{
				Driver: constants.DriverKubernetes,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.setup.Validate()

			if test.failure {
				if err == nil {
					t.Errorf("Validate should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("Validate returned err: %v", err)
			}
		})
	}
}
