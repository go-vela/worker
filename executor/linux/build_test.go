// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"flag"
	"net/http/httptest"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"
	"github.com/urfave/cli/v2"
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

	tests := []struct {
		name     string
		failure  bool
		runtime  string
		build    *library.Build
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-steps pipeline with empty build",
			failure:  true,
			runtime:  constants.DriverDocker,
			build:    new(library.Build),
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-steps pipeline with empty build",
			failure:  true,
			runtime:  constants.DriverKubernetes,
			build:    new(library.Build),
			pipeline: "testdata/build/steps/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_repo).
				WithMetadata(_metadata).
				WithUser(_user).
				Compile(test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			_pipeline = _pipeline.Sanitize(test.runtime)

			var _runtime runtime.Engine

			switch test.runtime {
			case constants.DriverKubernetes:
				_pod := testPodFor(_pipeline)
				_runtime, err = kubernetes.NewMock(_pod)
				if err != nil {
					t.Errorf("unable to create kubernetes runtime engine: %v", err)
				}
			case constants.DriverDocker:
				_runtime, err = docker.NewMock()
				if err != nil {
					t.Errorf("unable to create docker runtime engine: %v", err)
				}
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
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.CreateBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("%s CreateBuild should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s CreateBuild returned err: %v", test.name, err)
			}
		})
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

	tests := []struct {
		name     string
		failure  bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_repo).
				WithMetadata(_metadata).
				WithUser(_user).
				Compile(test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			_pipeline = _pipeline.Sanitize(test.runtime)

			var _runtime runtime.Engine

			switch test.runtime {
			case constants.DriverKubernetes:
				_pod := testPodFor(_pipeline)
				_runtime, err = kubernetes.NewMock(_pod)
				if err != nil {
					t.Errorf("unable to create kubernetes runtime engine: %v", err)
				}
			case constants.DriverDocker:
				_runtime, err = docker.NewMock()
				if err != nil {
					t.Errorf("unable to create docker runtime engine: %v", err)
				}
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
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			err = _engine.PlanBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("%s PlanBuild should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s PlanBuild returned err: %v", test.name, err)
			}
		})
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

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	tests := []struct {
		name     string
		failure  bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-secrets pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-secrets pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/img_notfound.yml",
		//},
		{
			name:     "docker-secrets pipeline with ignoring image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-secrets pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/img_ignorenotfound.yml",
		//},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-services pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/img_notfound.yml",
		//},
		{
			name:     "docker-services pipeline with ignoring image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/img_ignorenotfound.yml",
		//},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-steps pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/img_notfound.yml",
		//},
		{
			name:     "docker-steps pipeline with ignoring image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/img_ignorenotfound.yml",
		//},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-stages pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/img_notfound.yml",
		//},
		{
			name:     "docker-stages pipeline with ignoring image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/img_ignorenotfound.yml",
		//},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_repo).
				WithMetadata(_metadata).
				WithUser(_user).
				Compile(test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			_pipeline = _pipeline.Sanitize(test.runtime)

			var _runtime runtime.Engine

			switch test.runtime {
			case constants.DriverKubernetes:
				_runtime, err = kubernetes.NewMock(&v1.Pod{}) // do not use _pod here! AssembleBuild will load it.
				if err != nil {
					t.Errorf("unable to create kubernetes runtime engine: %v", err)
				}
			case constants.DriverDocker:
				_runtime, err = docker.NewMock()
				if err != nil {
					t.Errorf("unable to create docker runtime engine: %v", err)
				}
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime == constants.DriverKubernetes {
				go func() {
					_mockRuntime := _runtime.(kubernetes.MockKubernetesRuntime)
					// This handles waiting until runtime.AssembleBuild has prepared the PodTracker.
					_mockRuntime.WaitForPodTrackerReady()
					// Normally, runtime.StreamBuild (which runs in a goroutine) calls PodTracker.Start.
					_mockRuntime.StartPodTracker(context.Background())

					_pod := testPodFor(_pipeline)

					// Now wait until the pod is created at the end of runtime.AssembleBuild.
					_mockRuntime.WaitForPodCreate(_pod.GetNamespace(), _pod.GetName())

					var stepsRunningCount int

					percents := []int{0, 0, 50, 100}
					lastIndex := len(percents) - 1
					for index, stepsCompletedPercent := range percents {
						if index == 0 || index == lastIndex {
							stepsRunningCount = 0
						} else {
							stepsRunningCount = 1
						}

						err := _mockRuntime.SimulateStatusUpdate(_pod,
							testContainerStatuses(
								_pipeline, true, stepsRunningCount, stepsCompletedPercent,
							),
						)
						if err != nil {
							t.Errorf("%s - failed to simulate pod update: %s", test.name, err)
						}

						// simulate exec build duration
						time.Sleep(100 * time.Microsecond)
					}
				}()
			}

			err = _engine.AssembleBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("%s AssembleBuild should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s AssembleBuild returned err: %v", test.name, err)
			}
		})
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

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	tests := []struct {
		name     string
		failure  bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false, // fixed
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-services pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/img_notfound.yml",
		//},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-steps pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/img_notfound.yml",
		//},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-stages pipeline with image not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/img_notfound.yml",
		//},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_repo).
				WithMetadata(_metadata).
				WithUser(_user).
				Compile(test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			_pipeline = _pipeline.Sanitize(test.runtime)

			var (
				_runtime runtime.Engine
				_pod     *v1.Pod
			)

			switch test.runtime {
			case constants.DriverKubernetes:
				_pod = testPodFor(_pipeline)
				_runtime, err = kubernetes.NewMock(_pod)
				if err != nil {
					t.Errorf("unable to create kubernetes runtime engine: %v", err)
				}
			case constants.DriverDocker:
				_runtime, err = docker.NewMock()
				if err != nil {
					t.Errorf("unable to create docker runtime engine: %v", err)
				}
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("%s unable to create build: %v", test.name, err)
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
				t.Errorf("unable to create docker runtime volume: %v", err)
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

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime == constants.DriverKubernetes {
				err = _runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
				if err != nil {
					t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
				}

				_runtime.(kubernetes.MockKubernetesRuntime).StartPodTracker(context.Background())

				go func() {
					_runtime.(kubernetes.MockKubernetesRuntime).SimulateResync()

					var stepsRunningCount int

					percents := []int{0, 0, 50, 100}
					lastIndex := len(percents) - 1
					for index, stepsCompletedPercent := range percents {
						if index == 0 || index == lastIndex {
							stepsRunningCount = 0
						} else {
							stepsRunningCount = 1
						}

						err := _runtime.(kubernetes.MockKubernetesRuntime).SimulateStatusUpdate(_pod,
							testContainerStatuses(
								_pipeline, true, stepsRunningCount, stepsCompletedPercent,
							),
						)
						if err != nil {
							t.Errorf("%s - failed to simulate pod update: %s", test.name, err)
						}

						// simulate exec build duration
						time.Sleep(100 * time.Microsecond)
					}
				}()
			}

			err = _engine.ExecBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("%s ExecBuild for %s should have returned err", test.name, test.pipeline)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s ExecBuild for %s returned err: %v", test.name, test.pipeline, err)
			}
		})
	}
}

func TestLinux_StreamBuild(t *testing.T) {
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

	type planFuncType = func(context.Context, *pipeline.Container) error

	// planNothing is a planFuncType that does nothing
	planNothing := func(ctx context.Context, container *pipeline.Container) error {
		return nil
	}

	tests := []struct {
		name       string
		failure    bool
		runtime    string
		pipeline   string
		messageKey string
		ctn        *pipeline.Container
		streamFunc func(*client) message.StreamFunc
		planFunc   func(*client) planFuncType
	}{
		{
			name:       "docker-basic services pipeline",
			failure:    false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanService
			},
			ctn: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:       "kubernetes-basic services pipeline",
			failure:    false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanService
			},
			ctn: &pipeline.Container{
				ID:          "service-github-octocat-1-postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:       "docker-basic services pipeline with StreamService failure",
			failure:    false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(c *client) planFuncType {
				// simulate failure to call PlanService
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:       "kubernetes-basic services pipeline with StreamService failure",
			failure:    false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(c *client) planFuncType {
				// simulate failure to call PlanService
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "service-github-octocat-1-postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:       "docker-basic steps pipeline",
			failure:    false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanStep
			},
			ctn: &pipeline.Container{
				ID:          "step_github_octocat_1_test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "kubernetes-basic steps pipeline",
			failure:    false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanStep
			},
			ctn: &pipeline.Container{
				ID:          "step-github-octocat-1-test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "docker-basic steps pipeline with StreamStep failure",
			failure:    false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				// simulate failure to call PlanStep
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "step_github_octocat_1_test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "kubernetes-basic steps pipeline with StreamStep failure",
			failure:    false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				// simulate failure to call PlanStep
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "step-github-octocat-1-test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "docker-basic stages pipeline",
			failure:    false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/stages/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanStep
			},
			ctn: &pipeline.Container{
				ID:          "step_github_octocat_1_test_test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "kubernetes-basic stages pipeline",
			failure:    false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/stages/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanStep
			},
			ctn: &pipeline.Container{
				ID:          "step-github-octocat-1-test-test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "docker-basic secrets pipeline",
			failure:    false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/secrets/basic.yml",
			messageKey: "secret",
			streamFunc: func(c *client) message.StreamFunc {
				return c.secret.stream
			},
			planFunc: func(c *client) planFuncType {
				// no plan function equivalent for secret containers
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "kubernetes-basic secrets pipeline",
			failure:    false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/secrets/basic.yml",
			messageKey: "secret",
			streamFunc: func(c *client) message.StreamFunc {
				return c.secret.stream
			},
			planFunc: func(c *client) planFuncType {
				// no plan function equivalent for secret containers
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "secret-github-octocat-1-vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "not_present",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buildCtx, done := context.WithCancel(context.Background())
			defer done()

			streamRequests := make(chan message.StreamRequest)

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_repo).
				WithMetadata(_metadata).
				WithUser(_user).
				Compile(test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			_pipeline = _pipeline.Sanitize(test.runtime)

			var _runtime runtime.Engine

			switch test.runtime {
			case constants.DriverKubernetes:
				_pod := testPodFor(_pipeline)
				_runtime, err = kubernetes.NewMock(_pod)
				if err != nil {
					t.Errorf("unable to create kubernetes runtime engine: %v", err)
				}
			case constants.DriverDocker:
				_runtime, err = docker.NewMock()
				if err != nil {
					t.Errorf("unable to create docker runtime engine: %v", err)
				}
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRepo(_repo),
				WithRuntime(_runtime),
				WithUser(_user),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(buildCtx)
			if err != nil {
				t.Errorf("%s unable to create build: %v", test.name, err)
			}

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime == constants.DriverKubernetes {
				err = _runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
				if err != nil {
					t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
				}

				// Runtime.StreamBuild calls PodTracker.Start after the PodTracker is marked Ready
				_runtime.(kubernetes.MockKubernetesRuntime).MarkPodTrackerReady()
				// FIXME:
				//		msg="error while requesting pod/logs stream for container service-github-octocat-1-postgres: context canceled"
				//		msg="exponential backoff error while tailing container service-github-octocat-1-postgres: context canceled"
				// 		msg="exponential backoff error while tailing container service-github-octocat-1-postgres: context canceled"
				// 		msg="unable to tail container output for upload: context canceled" service=postgres
			}

			// simulate ExecBuild() which runs concurrently with StreamBuild()
			go func() {
				// ExecBuild calls PlanService()/PlanStep() before ExecService()/ExecStep()
				// (ExecStage() calls PlanStep() before ExecStep()).
				_engine.err = test.planFunc(_engine)(buildCtx, test.ctn)

				// ExecService()/ExecStep()/secret.exec() send this message
				streamRequests <- message.StreamRequest{
					Key:       test.messageKey,
					Stream:    test.streamFunc(_engine),
					Container: test.ctn,
				}

				// simulate exec build duration
				time.Sleep(100 * time.Microsecond)

				// signal the end of the build so StreamBuild can terminate
				done()
			}()

			err = _engine.StreamBuild(buildCtx)

			if test.failure {
				if err == nil {
					t.Errorf("%s StreamBuild for %s should have returned err", test.name, test.pipeline)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s StreamBuild for %s returned err: %v", test.name, test.pipeline, err)
			}
		})
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

	tests := []struct {
		name     string
		failure  bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-secrets pipeline with name not found",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-secrets pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/name_notfound.yml",
		//},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-services pipeline with name not found",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/name_notfound.yml",
		//},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-steps pipeline with name not found",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/name_notfound.yml",
		//},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-stages pipeline with name not found",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/name_notfound.yml",
		//},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_repo).
				WithMetadata(_metadata).
				WithUser(_user).
				Compile(test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			_pipeline = _pipeline.Sanitize(test.runtime)

			var _runtime runtime.Engine

			switch test.runtime {
			case constants.DriverKubernetes:
				_pod := testPodFor(_pipeline)
				_runtime, err = kubernetes.NewMock(_pod)
				if err != nil {
					t.Errorf("unable to create kubernetes runtime engine: %v", err)
				}
			case constants.DriverDocker:
				_runtime, err = docker.NewMock()
				if err != nil {
					t.Errorf("unable to create docker runtime engine: %v", err)
				}
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
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime == constants.DriverKubernetes {
				err = _runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
				if err != nil {
					t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
				}
			}

			err = _engine.DestroyBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("%s DestroyBuild should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s DestroyBuild returned err: %v", test.name, err)
			}
		})
	}
}
