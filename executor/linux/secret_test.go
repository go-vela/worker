// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"
)

func TestLinux_Secret_create(t *testing.T) {
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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
	}{
		{
			name:    "docker-good image tag",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
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
			name:    "kubernetes-good image tag",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "secret-github-octocat-1-vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:    "docker-notfound image tag",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:notfound",
				Name:        "vault",
				Number:      1,
				Pull:        "not_present",
			},
		},
		//{
		//	name:    "kubernetes-notfound image tag",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "secret-github-octocat-1-vault",
		//		Directory:   "/vela/src/vcs.company.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "target/secret-vault:notfound",
		//		Name:        "vault",
		//		Number:      1,
		//		Pull:        "not_present",
		//	},
		//},
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

			err = _engine.secret.create(context.Background(), test.container, "")

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

func TestLinux_Secret_delete(t *testing.T) {
	// setup types
	_build := testBuild()
	_dockerSteps := testSteps(constants.DriverDocker)
	_kubernetesSteps := testSteps(constants.DriverKubernetes)

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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	_step := new(api.Step)
	_step.SetName("clone")
	_step.SetNumber(2)
	_step.SetStatus(constants.StatusPending)

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		container *pipeline.Container
		step      *api.Step
		steps     *pipeline.Build
	}{
		{
			name:    "docker-running container-empty step",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "always",
			},
			step:  new(api.Step),
			steps: _dockerSteps,
		},
		{
			name:    "kubernetes-running container-empty step",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "secret-github-octocat-1-vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "always",
			},
			step:  new(api.Step),
			steps: _kubernetesSteps,
		},
		{
			name:    "docker-running container-pending step",
			failure: false,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      2,
				Pull:        "always",
			},
			step:  _step,
			steps: _dockerSteps,
		},
		{
			name:    "kubernetes-running container-pending step",
			failure: false,
			runtime: _kubernetes,
			container: &pipeline.Container{
				ID:          "secret-github-octocat-1-vault",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "vault",
				Number:      2,
				Pull:        "always",
			},
			step:  _step,
			steps: _kubernetesSteps,
		},
		{
			name:    "docker-inspecting container failure due to invalid container id",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_notfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "notfound",
				Number:      2,
				Pull:        "always",
			},
			step:  new(api.Step),
			steps: _dockerSteps,
		},
		//{
		//	name:    "kubernetes-inspecting container failure due to invalid container id",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "secret-github-octocat-1-notfound",
		//		Directory:   "/vela/src/vcs.company.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "target/secret-vault:latest",
		//		Name:        "notfound",
		//		Number:      2,
		//		Pull:        "always",
		//	},
		//	step:  new(api.Step),
		//	steps: _kubernetesSteps,
		//},
		{
			name:    "docker-removing container failure",
			failure: true,
			runtime: _docker,
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_ignorenotfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "ignorenotfound",
				Number:      2,
				Pull:        "always",
			},
			step:  new(api.Step),
			steps: _dockerSteps,
		},
		//{
		//	name:    "kubernetes-removing container failure",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	container: &pipeline.Container{
		//		ID:          "secret-github-octocat-1-ignorenotfound",
		//		Directory:   "/vela/src/vcs.company.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "target/secret-vault:latest",
		//		Name:        "ignorenotfound",
		//		Number:      2,
		//		Pull:        "always",
		//	},
		//	step:  new(api.Step),
		//	steps: _kubernetesSteps,
		//},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(test.steps),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime.Driver() == constants.DriverKubernetes {
				err = _engine.Runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
				if err != nil {
					t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
				}
			}

			_engine.steps.Store(test.container.ID, test.step)

			err = _engine.secret.destroy(context.Background(), test.container)

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

func TestLinux_Secret_exec(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)

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
			name:     "docker-pipeline with secret name not found",
			failure:  true,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-pipeline with secret name not found",
		//	failure:  true, // FIXME:  make Kubernetes mock simulate failure similar to Docker mock
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/name_notfound.yml",
		//},
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
				Compile(context.Background(), file)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			// Docker uses _ while Kubernetes uses -
			p = p.Sanitize(test.runtime)

			var _runtime runtime.Engine

			switch test.runtime {
			case constants.DriverKubernetes:
				_pod := testPodFor(p)
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
				WithPipeline(p),
				WithRuntime(_runtime),
				WithVelaClient(_client),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			_engine.build.SetStatus(constants.StatusSuccess)

			// add init container info to client
			_ = _engine.CreateBuild(context.Background())

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime == constants.DriverKubernetes {
				err = _runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
				if err != nil {
					t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
				}

				go _runtime.(kubernetes.MockKubernetesRuntime).SimulateResync(nil)
			}

			err = _engine.secret.exec(context.Background(), &p.Secrets)

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

func TestLinux_Secret_pull(t *testing.T) {
	// setup types
	_build := testBuild()

	gin.SetMode(gin.TestMode)

	server := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(server.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_docker, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create docker runtime engine: %v", err)
	}

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		runtime runtime.Engine
		secret  *pipeline.Secret
	}{
		{
			name:    "docker-success with org secret",
			failure: false,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-success with org secret",
			failure: false,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-failure with invalid org secret",
			failure: true,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo/foo/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-failure with invalid org secret",
			failure: true,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo/foo/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-failure with org secret key not found",
			failure: true,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-failure with org secret key not found",
			failure: true,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-success with repo secret",
			failure: false,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octocat/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-success with repo secret",
			failure: false,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octocat/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-failure with invalid repo secret",
			failure: true,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo/foo/foo/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-failure with invalid repo secret",
			failure: true,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo/foo/foo/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-failure with repo secret key not found",
			failure: true,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-failure with repo secret key not found",
			failure: true,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-success with shared secret",
			failure: false,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-success with shared secret",
			failure: false,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-failure with shared secret key not found",
			failure: true,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-failure with shared secret key not found",
			failure: true,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "not-found",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "docker-failure with invalid type",
			failure: true,
			runtime: _docker,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "invalid",
				Origin: &pipeline.Container{},
			},
		},
		{
			name:    "kubernetes-failure with invalid type",
			failure: true,
			runtime: _kubernetes,
			secret: &pipeline.Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "invalid",
				Origin: &pipeline.Container{},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(_build),
				WithPipeline(testSteps(constants.DriverDocker)),
				WithRuntime(test.runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			_, err = _engine.secret.pull(test.secret)

			if test.failure {
				if err == nil {
					t.Errorf("%s pull should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s pull returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_Secret_stream(t *testing.T) {
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

	_kubernetes, err := kubernetes.NewMock(testPod(false))
	if err != nil {
		t.Errorf("unable to create kubernetes runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		runtime   runtime.Engine
		logs      *api.Log
		container *pipeline.Container
	}{
		{
			name:    "docker-container step succeeds",
			failure: false,
			runtime: _docker,
			logs:    new(api.Log),
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
		},
		{
			name:    "kubernetes-container step succeeds",
			failure: false,
			runtime: _kubernetes,
			logs:    new(api.Log),
			container: &pipeline.Container{
				ID:          "step-github-octocat-1-init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
		},
		{
			name:    "docker-container step fails because of invalid container id",
			failure: true,
			runtime: _docker,
			logs:    new(api.Log),
			container: &pipeline.Container{
				ID:          "secret_github_octocat_1_notfound",
				Directory:   "/vela/src/vcs.company.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/secret-vault:latest",
				Name:        "notfound",
				Number:      2,
				Pull:        "always",
			},
		},
		//{
		//	name:    "kubernetes-container step fails because of invalid container id",
		//	failure: true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	runtime: _kubernetes,
		//	logs:    new(api.Log),
		//	container: &pipeline.Container{
		//		ID:          "secret-github-octocat-1-notfound",
		//		Directory:   "/vela/src/vcs.company.com/github/octocat",
		//		Environment: map[string]string{"FOO": "bar"},
		//		Image:       "target/secret-vault:latest",
		//		Name:        "notfound",
		//		Number:      2,
		//		Pull:        "always",
		//	},
		//},
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

			// Kubernetes runtime needs to set up the Mock after CreateBuild is called
			if test.runtime.Driver() == constants.DriverKubernetes {
				err = _engine.Runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
				if err != nil {
					t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
				}
			}

			err = _engine.secret.stream(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("%s stream should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s stream returned err: %v", test.name, err)
			}
		})
	}
}

func TestLinux_Secret_injectSecret(t *testing.T) {
	// name and value of secret
	v := "foo"
	tBool := true

	// setup types
	tests := []struct {
		name string
		step *pipeline.Container
		msec map[string]*api.Secret
		want *pipeline.Container
	}{
		// Tests for secrets with image ACLs
		{
			name: "secret with empty image ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{"FOO": {Name: &v, Value: &v, Images: &[]string{""}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
			},
		},
		{
			name: "secret with matching image ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine"},
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Branch: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with matching image:tag ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine:latest"},
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Branch: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with non-matching image ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{"FOO": {Name: &v, Value: &v, Images: &[]string{"centos"}}},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: make(map[string]string),
			},
		},

		// Tests for secrets with event ACLs
		{ // push event checks
			name: "secret with matching push event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine:latest"},
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Branch: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with non-matching push event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:  &v,
					Value: &v,
					AllowEvents: &api.Events{
						Deployment: &actions.Deploy{
							Created: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
		},
		{ // pull_request event checks
			name: "secret with matching pull_request event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "pull_request", "VELA_BUILD_EVENT_ACTION": "opened"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine:latest"},
					AllowEvents: &api.Events{
						PullRequest: &actions.Pull{
							Opened: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "pull_request", "VELA_BUILD_EVENT_ACTION": "opened"},
			},
		},
		{
			name: "secret with non-matching pull_request event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "pull_request"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:  &v,
					Value: &v,
					AllowEvents: &api.Events{
						Deployment: &actions.Deploy{
							Created: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "pull_request"},
			},
		},
		{ // tag event checks
			name: "secret with matching tag event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "tag"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine:latest"},
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Tag: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "tag"},
			},
		},
		{
			name: "secret with non-matching tag event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "tag"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:  &v,
					Value: &v,
					AllowEvents: &api.Events{
						Deployment: &actions.Deploy{
							Created: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "tag"},
			},
		},
		{ // deployment event checks
			name: "secret with matching deployment event ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "deployment", "VELA_BUILD_EVENT_ACTION": "created"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine:latest"},
					AllowEvents: &api.Events{
						Deployment: &actions.Deploy{
							Created: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "deployment", "VELA_BUILD_EVENT_ACTION": "created"},
			},
		},
		{
			name: "secret with non-matching deployment event ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "deployment"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:  &v,
					Value: &v,
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Tag: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "deployment"},
			},
		},

		// Tests for secrets with event and image ACLs
		{
			name: "secret with matching event ACL and non-matching image ACL not injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:  &v,
					Value: &v,
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Branch: &tBool,
						},
					},
					Images: &[]string{"centos"},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with non-matching event ACL and matching image ACL not injected",
			step: &pipeline.Container{
				Image:       "centos:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:  &v,
					Value: &v,
					AllowEvents: &api.Events{
						PullRequest: &actions.Pull{
							Opened:      &tBool,
							Synchronize: &tBool,
							Reopened:    &tBool,
						},
					},
					Images: &[]string{"centos"},
				},
			},
			want: &pipeline.Container{
				Image:       "centos:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
		},
		{
			name: "secret with matching event ACL and matching image ACL injected",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Secrets:     pipeline.StepSecretSlice{{Source: "FOO", Target: "FOO"}},
			},
			msec: map[string]*api.Secret{
				"FOO": {
					Name:   &v,
					Value:  &v,
					Images: &[]string{"alpine:latest"},
					AllowEvents: &api.Events{
						Push: &actions.Push{
							Branch: &tBool,
						},
					},
				},
			},
			want: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"FOO": "foo", "VELA_BUILD_EVENT": "push"},
			},
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = injectSecrets(test.step, test.msec)
			got := test.step

			// Preferred use of reflect.DeepEqual(x, y interface) is giving false positives.
			// Switching to a Google library for increased clarity.
			// https://github.com/google/go-cmp
			if diff := cmp.Diff(test.want.Environment, got.Environment); diff != "" {
				t.Errorf("injectSecrets mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLinux_Secret_escapeNewlineSecrets(t *testing.T) {
	// name and value of secret
	n := "foo"
	v := "bar\\nbaz"
	vEscaped := "bar\\\nbaz"

	// desired secret value
	w := "bar\\\nbaz"

	// setup types
	tests := []struct {
		name      string
		secretMap map[string]*api.Secret
		want      map[string]*api.Secret
	}{
		{
			name:      "not escaped",
			secretMap: map[string]*api.Secret{"FOO": {Name: &n, Value: &v}},
			want:      map[string]*api.Secret{"FOO": {Name: &n, Value: &w}},
		},
		{
			name:      "already escaped",
			secretMap: map[string]*api.Secret{"FOO": {Name: &n, Value: &vEscaped}},
			want:      map[string]*api.Secret{"FOO": {Name: &n, Value: &w}},
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			escapeNewlineSecrets(test.secretMap)
			got := test.secretMap

			// Preferred use of reflect.DeepEqual(x, y interface) is giving false positives.
			// Switching to a Google library for increased clarity.
			// https://github.com/google/go-cmp
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("escapeNewlineSecrets mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
