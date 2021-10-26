// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/mock/server"
	"github.com/go-vela/pkg-runtime/runtime/docker"
	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"
)

func TestLinux_Driver(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	want := constants.DriverLinux

	_engine, err := New(
		WithBuild(testBuild()),
		WithHostname("localhost"),
		WithPipeline(testSteps()),
		WithRepo(testRepo()),
		WithRuntime(_runtime),
		WithUser(testUser()),
		WithVelaClient(_client),
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
