// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"flag"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
)

func TestLinux_Outputs_create(t *testing.T) {
	// setup types
	_build := testBuild()
	_steps := testSteps(constants.DriverDocker)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "good image tag",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "outputs",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "notfound image tag",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:notfound",
				Name:        "outputs",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "not supplied image tag",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "",
				Name:        "outputs",
				Number:      1,
				Pull:        "not_present",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_steps),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
				WithOutputCtn(test.container),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			err = _engine.outputs.create(context.Background(), test.container, 30)

			if test.failure {
				if err == nil {
					t.Errorf("%s create should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s create returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_Outputs_delete(t *testing.T) {
	// setup types
	_build := testBuild()
	_dockerSteps := testSteps(constants.DriverDocker)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	_step := new(library.Step)
	_step.SetName("clone")
	_step.SetNumber(2)
	_step.SetStatus(constants.StatusPending)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
		step      *library.Step
		steps     *pipeline.Build
	}{
		{
			name:    "docker-running container-empty step",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "outputs",
				Number:      1,
				Pull:        "always",
			},
			step:  new(library.Step),
			steps: _dockerSteps,
		},
		{
			name:    "docker-running container-pending step",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "outputs",
				Number:      2,
				Pull:        "always",
			},
			step:  _step,
			steps: _dockerSteps,
		},
		{
			name:    "docker-inspecting container failure due to invalid container id",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1_notfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "notfound",
				Number:      2,
				Pull:        "always",
			},
			step:  new(library.Step),
			steps: _dockerSteps,
		},
		{
			name:    "docker-removing container failure",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "ignorenotfound",
				Number:      2,
				Pull:        "always",
			},
			step:  new(library.Step),
			steps: _dockerSteps,
		},
		{
			name:    "no outputs image provided",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "",
				Name:        "outputs",
				Number:      2,
				Pull:        "always",
			},
			step:  _step,
			steps: _dockerSteps,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(test.steps),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
				WithOutputCtn(test.container),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			_engine.steps.Store(test.container.ID, test.step)

			err = _engine.outputs.destroy(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s destroy should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s destroy returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_Outputs_exec(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", "target/vela-git:latest", "doc")
	compiler, _ := native.FromCLIContext(cli.NewContext(nil, set, nil))

	_build := testBuild()

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	// setup tests
	tests := []struct {
		name     string
		failure  bool
		runtime  string
		pipeline string
	}{
		{
			name:     "basic pipeline",
			failure:  false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, _ := os.ReadFile(test.pipeline)

			p, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(file)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			p = p.Sanitize(test.runtime)

			var _runtime runtime.Engine

			_runtime, err = docker.NewMock()
			if err != nil {
				t.Errorf("unable to create docker runtime engine: %v", err)
			}

			outputsCtn := &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "outputs",
				Number:      2,
				Pull:        "always",
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(p),
				WithRuntime(_runtime),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
				WithOutputCtn(outputsCtn),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			_engine.build.SetStatus(constants.StatusSuccess)

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			err = _engine.outputs.exec(context.Background(), outputsCtn)

			if test.failure {
				if err == nil {
					t.Errorf("%s exec should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s exec returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_Outputs_poll(t *testing.T) {
	// setup types
	_build := testBuild()
	_steps := testSteps(constants.DriverDocker)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "succeeds",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "outputs",
				Number:      1,
				Pull:        "always",
			},
		},
		{
			name:    "no outputs image provided",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "outputs_github_octocat_1_notfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "",
				Name:        "notfound",
				Number:      2,
				Pull:        "always",
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_steps),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			_, _, err = _engine.outputs.poll(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s poll should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s poll returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_Outputs_toMap(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		runtime string
		input   []byte
		wantMap map[string]string
	}{
		{
			name:    "basic",
			runtime: constants.DriverDocker,
			input:   []byte("FOO=bar\r\nONE=1\r\nTEST=hello world\r\n"),
			wantMap: map[string]string{
				"FOO":  "bar",
				"ONE":  "1",
				"TEST": "hello world",
			},
		},
		{
			name:    "bad format",
			runtime: constants.DriverDocker,
			input:   []byte("FOO;bar//ONE^TEST,,,hello world"),
			wantMap: make(map[string]string),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := toMap(test.input)

			if !reflect.DeepEqual(got, test.wantMap) {
				t.Errorf("toMap is %v, want %v", got, test.wantMap)
			}
		})
	}
}
