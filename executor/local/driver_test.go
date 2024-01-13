// SPDX-License-Identifier: Apache-2.0

package local

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/worker/runtime/docker"
)

func TestLocal_Driver(t *testing.T) {
	// setup types
	want := constants.DriverLocal

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	_engine, err := New(
		WithBuild(testBuild()),
		WithHostname("localhost"),
		WithPipeline(testSteps()),
		WithRuntime(_runtime),
	)
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	// run tes
	got := _engine.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
