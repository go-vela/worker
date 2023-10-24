// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestKubernetes_Driver(t *testing.T) {
	// setup types
	want := constants.DriverKubernetes

	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// run tes
	got := _engine.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
