// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"context"
	"flag"
	"testing"

	"github.com/go-vela/compiler/compiler/native"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/pkg-runtime/runtime/docker"
)

func TestLocal_CreateBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

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
			WithLocal(true).
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

func TestLocal_PlanBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

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
			WithLocal(true).
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

func TestLocal_AssembleBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

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
			WithLocal(true).
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

func TestLocal_ExecBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

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
			WithLocal(true).
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
		)
		if err != nil {
			t.Errorf("unable to create executor engine: %v", err)
		}

		// run create to init steps to be created properly
		err = _engine.CreateBuild(context.Background())
		if err != nil {
			t.Errorf("unable to create build: %v", err)
		}

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

func TestLocal_DestroyBuild(t *testing.T) {
	// setup types
	compiler, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	_build := testBuild()
	_repo := testRepo()
	_user := testUser()

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
			WithLocal(true).
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
