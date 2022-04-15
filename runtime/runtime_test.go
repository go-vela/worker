// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"testing"

	"github.com/go-vela/types/constants"
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
		_, err := New(test.setup)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}
