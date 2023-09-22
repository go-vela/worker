// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestKubernetes_InspectImage(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
	}{
		{
			name:      "valid image",
			failure:   false,
			container: _container,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.InspectImage(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("InspectImage should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectImage returned err: %v", err)
			}
		})
	}
}
