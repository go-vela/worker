// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"flag"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/compiler/compiler/native"
	"github.com/go-vela/mock/server"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestLinux_CreateBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

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

	tests := []struct {
		failure  bool
		build    *library.Build
		pipeline string
	}{
		{ // basic secrets pipeline
			failure:  false,
			build:    _build,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{ // basic services pipeline
			failure:  false,
			build:    _build,
			pipeline: "testdata/build/services/basic.yml",
		},
		{ // basic steps pipeline
			failure:  false,
			build:    _build,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{ // basic stages pipeline
			failure:  false,
			build:    _build,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{ // steps pipeline with empty build
			failure:  true,
			build:    new(library.Build),
			pipeline: "testdata/build/steps/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		_pipeline, err := compiler.
			WithBuild(_build).
			WithRepo(_repo).
			WithMetadata(_metadata).
			WithUser(_user).
			Compile(test.pipeline)
		if err != nil {
			t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
		}

		_engine, err := New(
			WithBuild(test.build),
			WithPipeline(_pipeline),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		err = _engine.CreateBuild(context.Background())

		if test.failure {
			if err == nil {
				t.Errorf("CreateBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateBuild returned err: %v", err)
		}
	}
}

func TestLinux_PlanBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

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

	tests := []struct {
		failure  bool
		pipeline string
	}{
		{ // basic secrets pipeline
			failure:  false,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{ // basic services pipeline
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{ // basic steps pipeline
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{ // basic stages pipeline
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		_pipeline, err := compiler.
			WithBuild(_build).
			WithRepo(_repo).
			WithMetadata(_metadata).
			WithUser(_user).
			Compile(test.pipeline)
		if err != nil {
			t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
		}

		_engine, err := New(
			WithBuild(_build),
			WithPipeline(_pipeline),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		// run create to init steps to be created properly
		err = _engine.CreateBuild(context.Background())
		if err != nil {
			t.Errorf("unable to create build: %v", err)
		}

		err = _engine.PlanBuild(context.Background())

		if test.failure {
			if err == nil {
				t.Errorf("PlanBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("PlanBuild returned err: %v", err)
		}
	}
}

func TestLinux_AssembleBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

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

	tests := []struct {
		failure  bool
		pipeline string
	}{
		{ // basic secrets pipeline
			failure:  false,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{ // secrets pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/secrets/img_notfound.yml",
		},
		{ // secrets pipeline with ignoring image not found
			failure:  true,
			pipeline: "testdata/build/secrets/img_ignorenotfound.yml",
		},
		{ // basic services pipeline
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{ // services pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		{ // services pipeline with ignoring image not found
			failure:  true,
			pipeline: "testdata/build/services/img_ignorenotfound.yml",
		},
		{ // basic steps pipeline
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{ // steps pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		{ // steps pipeline with ignoring image not found
			failure:  true,
			pipeline: "testdata/build/steps/img_ignorenotfound.yml",
		},
		{ // basic stages pipeline
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{ // stages pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
		{ // stages pipeline with ignoring image not found
			failure:  true,
			pipeline: "testdata/build/stages/img_ignorenotfound.yml",
		},
	}

	// run test
	for _, test := range tests {
		_pipeline, err := compiler.
			WithBuild(_build).
			WithRepo(_repo).
			WithMetadata(_metadata).
			WithUser(_user).
			Compile(test.pipeline)
		if err != nil {
			t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
		}

		_engine, err := New(
			WithBuild(_build),
			WithPipeline(_pipeline),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		// run create to init steps to be created properly
		err = _engine.CreateBuild(context.Background())
		if err != nil {
			t.Errorf("unable to create build: %v", err)
		}

		err = _engine.AssembleBuild(context.Background())

		if test.failure {
			if err == nil {
				t.Errorf("AssembleBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("AssembleBuild returned err: %v", err)
		}
	}
}

func TestLinux_ExecBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

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

	tests := []struct {
		failure  bool
		pipeline string
	}{
		{ // basic services pipeline
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{ // services pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		{ // basic steps pipeline
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{ // steps pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		{ // basic stages pipeline
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{ // stages pipeline with image not found
			failure:  true,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
	}

	// run test
	for _, test := range tests {
		_pipeline, err := compiler.
			WithBuild(_build).
			WithRepo(_repo).
			WithMetadata(_metadata).
			WithUser(_user).
			Compile(test.pipeline)
		if err != nil {
			t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
		}

		_engine, err := New(
			WithBuild(_build),
			WithPipeline(_pipeline),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		// run create to init steps to be created properly
		err = _engine.CreateBuild(context.Background())
		if err != nil {
			t.Errorf("unable to create build: %v", err)
		}

		// TODO: hack - remove this
		//
		// When using our Docker mock we default the list of
		// Docker images that have privileged access. One of
		// these images is target/vela-git which is injected
		// by the compiler into every build.
		//
		// The problem this causes is that we haven't called
		// all the necessary functions we do in a real world
		// scenario to ensure we can run privileged images.
		//
		// This is the necessary function to create the
		// runtime host config so we can run images
		// in a privileged fashion.
		err = _runtime.CreateVolume(context.Background(), _pipeline)
		if err != nil {
			t.Errorf("unable to create runtime volume: %w", err)
		}

		// TODO: hack - remove this
		//
		// When calling CreateBuild(), it will automatically set the
		// test build object to a status of `created`. This happens
		// because we use a mock for the go-vela/server in our tests
		// which only returns dummy based responses.
		//
		// The problem this causes is that our container.Execute()
		// function isn't setup to handle builds in a `created` state.
		//
		// In a real world scenario, we never would have a build
		// in this state when we call ExecBuild() because the
		// go-vela/server has logic to set it to an expected state.
		_engine.build.SetStatus("running")

		err = _engine.ExecBuild(context.Background())

		if test.failure {
			if err == nil {
				t.Errorf("ExecBuild for %s should have returned err", test.pipeline)
			}

			continue
		}

		if err != nil {
			t.Errorf("ExecBuild for %s returned err: %v", test.pipeline, err)
		}
	}
}

func TestLinux_DestroyBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()
	_metadata := testMetadata()

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

	tests := []struct {
		failure  bool
		pipeline string
	}{
		{ // basic secrets pipeline
			failure:  false,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{ // secrets pipeline with name not found
			failure:  false,
			pipeline: "testdata/build/secrets/name_notfound.yml",
		},
		{ // basic services pipeline
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{ // services pipeline with name not found
			failure:  false,
			pipeline: "testdata/build/services/name_notfound.yml",
		},
		{ // basic steps pipeline
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{ // steps pipeline with name not found
			failure:  false,
			pipeline: "testdata/build/steps/name_notfound.yml",
		},
		{ // basic stages pipeline
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{ // stages pipeline with name not found
			failure:  false,
			pipeline: "testdata/build/stages/name_notfound.yml",
		},
	}

	// run test
	for _, test := range tests {
		_pipeline, err := compiler.
			WithBuild(_build).
			WithRepo(_repo).
			WithMetadata(_metadata).
			WithUser(_user).
			Compile(test.pipeline)
		if err != nil {
			t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
		}

		_engine, err := New(
			WithBuild(_build),
			WithPipeline(_pipeline),
			WithRepo(_repo),
			WithRuntime(_runtime),
			WithUser(_user),
			WithVelaClient(_client),
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		// run create to init steps to be created properly
		err = _engine.CreateBuild(context.Background())
		if err != nil {
			t.Errorf("unable to create build: %v", err)
		}

		err = _engine.DestroyBuild(context.Background())

		if test.failure {
			if err == nil {
				t.Errorf("DestroyBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DestroyBuild returned err: %v", err)
		}
	}
}
