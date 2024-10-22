// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"testing"

	"github.com/go-vela/server/constants"
)

func TestRuntime_New(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		setup   *Setup
	}{
		{
			name:    "docker driver",
			failure: false,
			setup: &Setup{
				Driver: constants.DriverDocker,
			},
		},
		{
			name:    "kubernetes driver",
			failure: false,
			setup: &Setup{
				Driver:     constants.DriverKubernetes,
				Namespace:  "docker",
				ConfigFile: "testdata/config",
			},
		},
		{
			name:    "invalid driver fails",
			failure: true,
			setup: &Setup{
				Driver: "invalid",
			},
		},
		{
			name:    "empty driver fails",
			failure: true,
			setup: &Setup{
				Driver: "",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.setup)

			if test.failure {
				if err == nil {
					t.Errorf("New should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("New returned err: %v", err)
			}
		})
	}
}
