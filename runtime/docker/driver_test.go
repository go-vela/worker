// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestDocker_Driver(t *testing.T) {
	// setup types
	want := constants.DriverDocker

	_engine, err := NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// run tes
	got := _engine.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
