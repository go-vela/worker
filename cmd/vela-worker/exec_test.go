// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/server/queue/redis"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/mock/worker"
	"github.com/go-vela/worker/runtime"
)

func TestWorker_exec(t *testing.T) {
	// setup types for tests
	_repo := testRepo()
	_user := testUser()

	type testStruct struct {
		name     string
		config   *workerTestConfig
		pipeline *pipeline.Build
		wantErr  bool
	}

	// this gets expanded into the product of runtime-executor-baseTest
	baseTests := []testStruct{
		{
			name:     "steps",
			pipeline: _steps,
			config:   &workerTestConfig{},
		},
		{
			name:     "stages",
			pipeline: _stages,
			config:   &workerTestConfig{},
		},
	}
	executors := []workerTestConfig{
		{
			name:           constants.DriverLinux,
			executorDriver: constants.DriverLinux,
			//executorLogMethod:
		},
		{
			name:           constants.DriverLocal,
			executorDriver: constants.DriverLocal,
			//executorLogMethod:
		},
	}
	runtimes := []workerTestConfig{
		{
			name:          constants.DriverDocker,
			runtimeDriver: constants.DriverDocker,
		},
		// TODO: kubernetes tests are hanging. Fix in a follow-up.
		//{
		//	name:              constants.DriverKubernetes,
		//	runtimeDriver:     constants.DriverKubernetes,
		//	runtimeNamespace:  "test",
		//	runtimeConfigFile: "../../runtime/kubernetes/testdata/config",
		//},
	}

	// if tests are needed beyond the matrix of tests, they can be added explicitly here.
	var tests []testStruct

	for _, r := range runtimes {
		for _, e := range executors {
			for _, test := range baseTests {
				tests = append(tests, testStruct{
					// matrix name format: runtime-executor-subtest
					name:     fmt.Sprintf("%s-%s-%s", r.name, e.name, test.name),
					pipeline: test.pipeline,
					config: (&workerTestConfig{
						// executor
						executorDriver:     e.executorDriver,
						executorLogMethod:  e.executorLogMethod,
						executorMaxLogSize: e.executorMaxLogSize,
						// runtime
						runtimeDriver:           r.runtimeDriver,
						runtimeConfigFile:       r.runtimeConfigFile,
						runtimeNamespace:        r.runtimeNamespace,
						runtimePodsTemplateName: r.runtimePodsTemplateName,
						runtimePodsTemplateFile: r.runtimePodsTemplateFile,
						runtimePrivilegedImages: r.runtimePrivilegedImages,
						runtimeHostVolumes:      r.runtimeHostVolumes,
						// etc
						buildTimeout:    test.config.buildTimeout,
						logFormat:       test.config.logFormat,
						logLevel:        test.config.logLevel,
						route:           test.config.route,
						queueRoutes:     test.config.queueRoutes,
						queuePopTimeout: test.config.queuePopTimeout,
					}).applyDefaults(),
				})
			}
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				err      error
				w        *Worker
				byteItem []byte
			)

			// make the worker with mocked
			w, err = mockWorker(test.config)
			if err != nil {
				t.Errorf("mockWorker error = %v", err)
			}

			_build := testBuild(w.Config)

			byteItem, err = json.Marshal(types.Item{
				Build:    _build,
				Pipeline: test.pipeline,
				Repo:     _repo,
				User:     _user,
			})
			if err != nil {
				t.Errorf("queue Item marshall error = %v", err)
			}

			// add Item to the queue
			err = w.Queue.Push(context.Background(), test.config.route, byteItem)
			if err != nil {
				t.Errorf("queue push error = %v", err)
			}

			// actually run our test
			if err = w.exec(0); (err != nil) != test.wantErr {
				t.Errorf("exec() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

type workerTestConfig struct {
	name                    string
	buildTimeout            time.Duration
	executorDriver          string
	executorLogMethod       string
	executorMaxLogSize      uint
	logFormat               string
	logLevel                string
	runtimeDriver           string
	runtimeConfigFile       string // only for k8s
	runtimeNamespace        string // only for k8s
	runtimePodsTemplateName string // only for k8s
	runtimePodsTemplateFile string // only for k8s
	runtimePrivilegedImages []string
	runtimeHostVolumes      []string
	route                   string
	queueRoutes             []string
	queuePopTimeout         time.Duration
}

// applyDefaults applies defaults for tests
// nolint:wsl // wsl prevents visual grouping of defaults
func (c *workerTestConfig) applyDefaults() *workerTestConfig {
	// defaults from flags.go (might need to lower for tests)
	if c.buildTimeout == 0 {
		c.buildTimeout = 30 * time.Minute
	}
	if c.logFormat == "" {
		c.logFormat = "json"
	}
	if c.logLevel == "" {
		c.logLevel = "info"
	}

	// defaults from executor.Flags
	if c.executorLogMethod == "" {
		c.executorLogMethod = "byte-chunks"
	}

	// defaults from runtime.Flags
	if c.runtimeDriver == "" {
		c.runtimeDriver = constants.DriverDocker
	}
	if c.runtimePrivilegedImages == nil {
		c.runtimePrivilegedImages = []string{"target/vela-docker"}
	}

	// defaults from runtime.Flags
	if c.queueRoutes == nil {
		c.queueRoutes = []string{constants.DefaultRoute}
	}
	if c.queuePopTimeout == 0 {
		c.queuePopTimeout = 60 * time.Second
	}

	// convenient default (not based on anything)
	if c.route == "" {
		c.route = constants.DefaultRoute
	}

	return c
}

// mockWorker creates a Worker with mocks for the Vela server,
// the queue, and the runtime client(s).
func mockWorker(cfg *workerTestConfig) (*Worker, error) {
	var err error

	// Worker initialized in run()
	w := &Worker{
		// worker configuration (skipping fields unused by exec())
		Config: &Config{
			// api configuration
			API: &API{},
			// build configuration
			Build: &Build{
				Limit:   1,
				Timeout: cfg.buildTimeout,
			},
			// executor configuration
			Executor: &executor.Setup{
				Driver:     cfg.executorDriver,
				LogMethod:  cfg.executorLogMethod,
				MaxLogSize: cfg.executorMaxLogSize,
			},
			// logger configuration
			Logger: &Logger{
				Format: cfg.logFormat,
				Level:  cfg.logLevel,
			},
			// runtime configuration
			Runtime: &runtime.Setup{
				Mock:             true,
				Driver:           cfg.runtimeDriver,
				ConfigFile:       cfg.runtimeConfigFile,
				Namespace:        cfg.runtimeNamespace,
				PodsTemplateName: cfg.runtimePodsTemplateName,
				PodsTemplateFile: cfg.runtimePodsTemplateFile,
				HostVolumes:      cfg.runtimeHostVolumes,
				PrivilegedImages: cfg.runtimePrivilegedImages,
			},
			// server configuration
			Server: &Server{
				// address is mocked below
				Secret: "server.secret",
			},
		},
		// exec() creates the runtime (including the mocked runtime).
		// exec() creates the executor and adds it here.
		Executors: make(map[int]executor.Engine),
	}

	// setup mock vela server
	s := httptest.NewServer(server.FakeHandler())
	w.Config.Server.Address = s.URL

	// setup mock vela worker
	api := httptest.NewServer(worker.FakeHandler())

	w.Config.API.Address, err = url.Parse(api.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse mock address: %w", err)
	}

	// set up VelaClient (setup happens in operate())
	w.VelaClient, err = setupClient(w.Config.Server)
	if err != nil {
		return nil, fmt.Errorf("vela client setup error = %w", err)
	}

	// set up mock Redis client (setup happens in operate())
	w.Queue, err = redis.NewTest(cfg.queueRoutes...)
	if err != nil {
		return nil, fmt.Errorf("queue setup error = %w", err)
	}

	return w, nil
}

var (
	_stages = &pipeline.Build{
		Version: "1",
		ID:      "github-octocat-1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service-github-octocat-1-postgres",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      4,
				Ports:       []string{"5432:5432"},
			},
		},
		Stages: pipeline.StageSlice{
			{
				Name: "init",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step-github-octocat-1-init-init",
						Directory:   "/vela/src/github.com/octocat/helloworld",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
				},
			},
			{
				Name:  "clone",
				Needs: []string{"init"},
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step-github-octocat-1-clone-clone",
						Directory:   "/vela/src/github.com/octocat/helloworld",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-git:v0.4.0",
						Name:        "clone",
						Number:      2,
						Pull:        "always",
					},
				},
			},
			{
				Name:  "echo",
				Needs: []string{"clone"},
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step-github-octocat-1-echo-echo",
						Commands:    []string{"echo hello"},
						Detach:      true,
						Directory:   "/vela/src/github.com/octocat/helloworld",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      3,
						Pull:        "always",
					},
				},
			},
		},
	}

	_steps = &pipeline.Build{
		Version: "1",
		ID:      "github-octocat-1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service-github-octocat-1-postgres",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      4,
				Ports:       []string{"5432:5432"},
			},
		},
		Steps: pipeline.ContainerSlice{
			{
				ID:          "step-github-octocat-1-init",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step-github-octocat-1-clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step-github-octocat-1-echo",
				Commands:    []string{"echo hello"},
				Detach:      true,
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
			},
		},
	}
)

// testBuild is a test helper function to create a Build
// type with all fields set to a fake value.
func testBuild(cfg *Config) *library.Build {
	return &library.Build{
		ID:           vela.Int64(1),
		Number:       vela.Int(1),
		Parent:       vela.Int(1),
		Event:        vela.String("push"),
		Status:       vela.String("success"),
		Error:        vela.String(""),
		Enqueued:     vela.Int64(1563474077),
		Created:      vela.Int64(1563474076),
		Started:      vela.Int64(1563474077),
		Finished:     vela.Int64(0),
		Deploy:       vela.String(""),
		Clone:        vela.String("https://github.com/github/octocat.git"),
		Source:       vela.String("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        vela.String("push received from https://github.com/github/octocat"),
		Message:      vela.String("First commit..."),
		Commit:       vela.String("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       vela.String("OctoKitty"),
		Author:       vela.String("OctoKitty"),
		Branch:       vela.String("master"),
		Ref:          vela.String("refs/heads/master"),
		BaseRef:      vela.String(""),
		Host:         vela.String(cfg.API.Address.Host),
		Runtime:      vela.String(cfg.Runtime.Driver),
		Distribution: vela.String(cfg.Executor.Driver),
	}
}

// testRepo is a test helper function to create a Repo
// type with all fields set to a fake value.
func testRepo() *library.Repo {
	return &library.Repo{
		ID:          vela.Int64(1),
		Org:         vela.String("github"),
		Name:        vela.String("octocat"),
		FullName:    vela.String("github/octocat"),
		Link:        vela.String("https://github.com/github/octocat"),
		Clone:       vela.String("https://github.com/github/octocat.git"),
		Branch:      vela.String("master"),
		Timeout:     vela.Int64(60),
		Visibility:  vela.String("public"),
		Private:     vela.Bool(false),
		Trusted:     vela.Bool(false),
		Active:      vela.Bool(true),
		AllowPull:   vela.Bool(false),
		AllowPush:   vela.Bool(true),
		AllowDeploy: vela.Bool(false),
		AllowTag:    vela.Bool(false),
	}
}

// testUser is a test helper function to create a User
// type with all fields set to a fake value.
func testUser() *library.User {
	return &library.User{
		ID:        vela.Int64(1),
		Name:      vela.String("octocat"),
		Token:     vela.String("superSecretToken"),
		Hash:      vela.String("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy"),
		Favorites: vela.Strings([]string{"github/octocat"}),
		Active:    vela.Bool(true),
		Admin:     vela.Bool(false),
	}
}
