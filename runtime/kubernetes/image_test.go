// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
		_, err = _engine.InspectImage(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("InspectImage should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectImage returned err: %v", err)
		}
	}
}
