// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/urfave/cli/v3"
	v1 "k8s.io/api/core/v1"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"
)

func TestLinux_UpdateSCMAuth_ConcurrentRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var refreshCalls int32

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/install_token") {
			atomic.AddInt32(&refreshCalls, 1)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"refreshed-token","expiration":4102444800}`))

			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	apiClient, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Fatalf("unable to create Vela API client: %v", err)
	}

	// Use a token expiring soon to force refresh for the first caller.
	apiClient.Authentication.SetBuildTokenAuth("build-token", "initial-token", time.Now().Unix()+60, "github/octocat", 1)

	engine := &client{
		Logger: logrus.NewEntry(logrus.New()),
		Vela:   apiClient,
		build:  testBuild(),
	}

	const goroutines = 10

	ctns := make([]*pipeline.Container, goroutines)
	for i := range ctns {
		ctns[i] = &pipeline.Container{Environment: map[string]string{}}
	}

	errCh := make(chan error, goroutines)

	var wg sync.WaitGroup
	for i := range ctns {
		wg.Add(1)

		ctn := ctns[i]

		go func() {
			defer wg.Done()

			err := engine.UpdateSCMAuth(context.Background(), ctn)
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("UpdateSCMAuth returned err: %v", err)
		}
	}

	if got, want := atomic.LoadInt32(&refreshCalls), int32(1); got != want {
		t.Fatalf("unexpected refresh call count: got %d, want %d", got, want)
	}

	for i, ctn := range ctns {
		if got, want := ctn.Environment["VELA_NETRC_PASSWORD"], "refreshed-token"; got != want {
			t.Fatalf("container %d has unexpected VELA_NETRC_PASSWORD: got %q, want %q", i, got, want)
		}

		if got, want := ctn.Environment["VELA_GIT_TOKEN"], "refreshed-token"; got != want {
			t.Fatalf("container %d has unexpected VELA_GIT_TOKEN: got %q, want %q", i, got, want)
		}

		if got, want := ctn.Environment["VELA_GIT_TOKEN_EXPIRATION"], "4102444800"; got != want {
			t.Fatalf("container %d has unexpected VELA_GIT_TOKEN_EXPIRATION: got %q, want %q", i, got, want)
		}
	}
}

func TestLinux_UpdateSCMAuthForAllContainers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		failure       bool
		refreshFails  bool
		expiresIn     time.Duration
		wantToken     string
		wantRefreshes int32
	}{
		{
			name:          "valid token is refreshed before assembly",
			expiresIn:     2 * time.Hour,
			wantToken:     "refreshed-token",
			wantRefreshes: 1,
		},
		{
			name:          "failed refresh keeps a valid token",
			refreshFails:  true,
			expiresIn:     2 * time.Hour,
			wantToken:     "initial-token",
			wantRefreshes: 1,
		},
		{
			name:         "failed refresh of an expiring token fails",
			failure:      true,
			refreshFails: true,
			expiresIn:    time.Minute,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var refreshCalls int32

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/install_token") {
					atomic.AddInt32(&refreshCalls, 1)

					if test.refreshFails {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = w.Write([]byte(`{"error":"unable to refresh install token"}`))

						return
					}

					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"token":"refreshed-token","expiration":4102444800}`))

					return
				}

				w.WriteHeader(http.StatusNotFound)
			}))
			defer s.Close()

			apiClient, err := vela.NewClient(s.URL, "", nil)
			if err != nil {
				t.Fatalf("unable to create Vela API client: %v", err)
			}

			apiClient.Authentication.SetBuildTokenAuth("build-token", "initial-token", time.Now().Add(test.expiresIn).Unix(), "github/octocat", 1)

			_init := &pipeline.Container{ID: "step_github_octocat_1_init", Name: constants.InitName, Environment: map[string]string{}}
			_test := &pipeline.Container{ID: "step_github_octocat_1_test", Name: "test", Environment: map[string]string{}}
			_release := &pipeline.Container{
				ID:          "step_github_octocat_1_release",
				Name:        "release",
				Environment: map[string]string{},
				Ruleset: pipeline.Ruleset{
					If: pipeline.Rules{
						Branch: []string{"release"},
						Status: []string{constants.StatusSuccess, constants.StatusFailure},
					},
				},
			}

			engine := &client{
				Logger:   logrus.NewEntry(logrus.New()),
				Vela:     apiClient,
				build:    testBuild(),
				pipeline: &pipeline.Build{Steps: pipeline.ContainerSlice{_init, _test, _release}},
			}

			err = engine.updateSCMAuthForAllContainers(context.Background())

			if test.failure {
				if err == nil {
					t.Fatal("updateSCMAuthForAllContainers should have returned err")
				}

				return
			}

			if err != nil {
				t.Fatalf("updateSCMAuthForAllContainers returned err: %v", err)
			}

			if got := atomic.LoadInt32(&refreshCalls); got != test.wantRefreshes {
				t.Errorf("unexpected refresh call count: got %d, want %d", got, test.wantRefreshes)
			}

			if got := _test.Environment["VELA_NETRC_PASSWORD"]; got != test.wantToken {
				t.Errorf("runnable step has VELA_NETRC_PASSWORD %q, want %q", got, test.wantToken)
			}

			for _, ctn := range []*pipeline.Container{_init, _release} {
				if _, ok := ctn.Environment["VELA_NETRC_PASSWORD"]; ok {
					t.Errorf("step %s should not receive SCM credentials", ctn.Name)
				}
			}
		})
	}
}

func TestLinux_UpdateSCMAuth_ConcurrentNoRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var refreshCalls int32

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/install_token") {
			atomic.AddInt32(&refreshCalls, 1)
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	apiClient, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Fatalf("unable to create Vela API client: %v", err)
	}

	// Use a token expiration far in the future so refresh is not needed.
	expiration := time.Now().Unix() + int64((2 * time.Hour).Seconds())
	apiClient.Authentication.SetBuildTokenAuth("build-token", "initial-token", expiration, "github/octocat", 1)

	engine := &client{
		Logger: logrus.NewEntry(logrus.New()),
		Vela:   apiClient,
		build:  testBuild(),
	}

	const goroutines = 10

	ctns := make([]*pipeline.Container, goroutines)
	for i := range ctns {
		ctns[i] = &pipeline.Container{Environment: map[string]string{}}
	}

	errCh := make(chan error, goroutines)

	var wg sync.WaitGroup
	for i := range ctns {
		wg.Add(1)

		ctn := ctns[i]

		go func() {
			defer wg.Done()

			err := engine.UpdateSCMAuth(context.Background(), ctn)
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("UpdateSCMAuth returned err: %v", err)
		}
	}

	if got, want := atomic.LoadInt32(&refreshCalls), int32(0); got != want {
		t.Fatalf("unexpected refresh call count: got %d, want %d", got, want)
	}

	wantExp := strconv.FormatInt(expiration, 10)

	for i, ctn := range ctns {
		if got, want := ctn.Environment["VELA_NETRC_PASSWORD"], "initial-token"; got != want {
			t.Fatalf("container %d has unexpected VELA_NETRC_PASSWORD: got %q, want %q", i, got, want)
		}

		if got, want := ctn.Environment["VELA_GIT_TOKEN"], "initial-token"; got != want {
			t.Fatalf("container %d has unexpected VELA_GIT_TOKEN: got %q, want %q", i, got, want)
		}

		if got := ctn.Environment["VELA_GIT_TOKEN_EXPIRATION"]; got != wantExp {
			t.Fatalf("container %d has unexpected VELA_GIT_TOKEN_EXPIRATION: got %q, want %q", i, got, wantExp)
		}
	}
}

func TestLinux_CreateBuild(t *testing.T) {
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
	if err != nil {
		t.Errorf("unable to create pipeline compiler: %v", err)
	}

	_build := testBuild()

	testLogger := logrus.New()
	loggerHook := logrusTest.NewLocal(testLogger)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		logError bool
		runtime  string
		build    *api.Build
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			build:    _build,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			build:    _build,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-steps pipeline with empty build",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			build:    new(api.Build),
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-steps pipeline with empty build",
			failure:  true,
			logError: false,
			runtime:  constants.DriverKubernetes,
			build:    new(api.Build),
			pipeline: "testdata/build/steps/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := testLogger.WithFields(logrus.Fields{"test": test.name})

			defer loggerHook.Reset()

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			_pipeline.Prepare(_build.GetRepo().GetOrg(), _build.GetRepo().GetName(), _build.GetNumber(), false)

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
				WithLogger(logger),
				WithBuild(test.build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
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

			loggedError := false

			for _, logEntry := range loggerHook.AllEntries() {
				// Many errors during StreamBuild get logged and ignored.
				// So, Make sure there are no errors logged during StreamBuild.
				if logEntry.Level == logrus.ErrorLevel {
					loggedError = true

					if !test.logError {
						t.Errorf("%s StreamBuild for %s logged an Error: %v", test.name, test.pipeline, logEntry.Message)
					}
				}
			}

			if test.logError && !loggedError {
				t.Errorf("%s StreamBuild for %s did not log an Error but should have", test.name, test.pipeline)
			}
		})
	}
}

func TestLinux_PlanBuild(t *testing.T) {
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
	if err != nil {
		t.Errorf("unable to create pipeline compiler: %v", err)
	}

	_build := testBuild()

	testLogger := logrus.New()
	loggerHook := logrusTest.NewLocal(testLogger)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		logError bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-basic steps with lazy and regular secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic_secrets.yml",
		},
		{
			name:     "kubernetes-basic steps with lazy and regular secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic_secrets.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := testLogger.WithFields(logrus.Fields{"test": test.name})

			defer loggerHook.Reset()

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			_pipeline.Prepare(_build.GetRepo().GetOrg(), _build.GetRepo().GetName(), _build.GetNumber(), false)

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
				WithLogger(logger),
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithVelaClient(_client),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("%s unable to create build: %v", test.name, err)
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

			loggedError := false

			for _, logEntry := range loggerHook.AllEntries() {
				// Many errors during StreamBuild get logged and ignored.
				// So, Make sure there are no errors logged during StreamBuild.
				if logEntry.Level == logrus.ErrorLevel {
					loggedError = true

					if !test.logError {
						t.Errorf("%s StreamBuild for %s logged an Error: %v", test.name, test.pipeline, logEntry.Message)
					}
				}
			}

			if test.logError && !loggedError {
				t.Errorf("%s StreamBuild for %s did not log an Error but should have", test.name, test.pipeline)
			}
		})
	}
}

func TestLinux_AssembleBuild(t *testing.T) {
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
	if err != nil {
		t.Errorf("unable to create pipeline compiler: %v", err)
	}

	_build := testBuild()

	testLogger := logrus.New()
	loggerHook := logrusTest.NewLocal(testLogger)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	// a valid SCM token, so the kubernetes runtime's forced refresh is observable
	_client.Authentication.SetBuildTokenAuth("build-token", "initial-token", time.Now().Add(2*time.Hour).Unix(), _build.GetRepo().GetFullName(), _build.GetNumber())

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	tests := []struct {
		name     string
		failure  bool
		logError bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-secrets pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-secrets pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/img_notfound.yml",
		//},
		{
			name:     "docker-secrets pipeline with ignoring image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-secrets pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/img_ignorenotfound.yml",
		//},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-services pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/img_notfound.yml",
		//},
		{
			name:     "docker-services pipeline with ignoring image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/img_ignorenotfound.yml",
		//},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-steps pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/img_notfound.yml",
		//},
		{
			name:     "docker-steps pipeline with ignoring image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/img_ignorenotfound.yml",
		//},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-stages pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/img_notfound.yml",
		//},
		{
			name:     "docker-stages pipeline with ignoring image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/img_ignorenotfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with ignoring image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/img_ignorenotfound.yml",
		//},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := testLogger.WithFields(logrus.Fields{"test": test.name})

			defer loggerHook.Reset()

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			_pipeline.Prepare(_build.GetRepo().GetOrg(), _build.GetRepo().GetName(), _build.GetNumber(), false)

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
				WithLogger(logger),
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithVelaClient(_client),
				WithOutputCtn(testOutputsCtn()),
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
					// This handles waiting until runtime.AssembleBuild has prepared and started the PodTracker.
					_mockRuntime.WaitForPodTrackerReady()

					_pod := testPodFor(_pipeline)

					// Now wait until the pod is created at the end of runtime.AssembleBuild.
					_mockRuntime.WaitForPodCreate(_pod.GetNamespace(), _pod.GetName())

					// Mark services running and secrets as completed, but no steps have started.
					err := _mockRuntime.SimulateStatusUpdate(_pod,
						testContainerStatuses(_pipeline, true, 0, 0),
					)
					if err != nil {
						t.Errorf("%s - failed to simulate pod update: %s", test.name, err)
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

			// the kubernetes pod environment is fixed at creation, so the refreshed
			// SCM token must already be in every runnable container of the pod
			if test.runtime == constants.DriverKubernetes {
				_pod := _runtime.(kubernetes.MockKubernetesRuntime).GetPod()

				for _, ctn := range scmAuthContainers(_engine, _pipeline) {
					if got, want := podContainerEnv(_pod, ctn.ID, "VELA_NETRC_PASSWORD"), "ghs_123"; got != want {
						t.Errorf("%s pod container %s has VELA_NETRC_PASSWORD %q, want %q", test.name, ctn.ID, got, want)
					}
				}
			}

			loggedError := false

			for _, logEntry := range loggerHook.AllEntries() {
				// Many errors during StreamBuild get logged and ignored.
				// So, Make sure there are no errors logged during StreamBuild.
				if logEntry.Level == logrus.ErrorLevel {
					loggedError = true

					if !test.logError {
						t.Errorf("%s StreamBuild for %s logged an Error: %v", test.name, test.pipeline, logEntry.Message)
					}
				}
			}

			if test.logError && !loggedError {
				t.Errorf("%s StreamBuild for %s did not log an Error but should have", test.name, test.pipeline)
			}
		})
	}
}

// podContainerEnv returns the value of an environment variable in a pod container.
func podContainerEnv(pod *v1.Pod, container, name string) string {
	for _, ctn := range pod.Spec.Containers {
		if ctn.Name != container {
			continue
		}

		for _, env := range ctn.Env {
			if env.Name == name {
				return env.Value
			}
		}
	}

	return ""
}

// scmAuthContainers returns every container that receives SCM credentials.
func scmAuthContainers(c *client, p *pipeline.Build) []*pipeline.Container {
	containers := []*pipeline.Container{}

	containers = append(containers, p.Services...)

	for _, _stage := range p.Stages {
		if _stage.Name == constants.InitName {
			continue
		}

		for _, _step := range _stage.Steps {
			if !c.neverRuns(_step) {
				containers = append(containers, _step)
			}
		}
	}

	for _, _step := range p.Steps {
		if _step.Name == constants.InitName || c.neverRuns(_step) {
			continue
		}

		containers = append(containers, _step)
	}

	for _, _secret := range p.Secrets {
		if _secret.Origin.Empty() {
			continue
		}

		containers = append(containers, _secret.Origin)
	}

	return containers
}

func TestLinux_ExecBuild(t *testing.T) {
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
	if err != nil {
		t.Errorf("unable to create compiler from CLI command: %v", err)
	}

	_build := testBuild()

	testLogger := logrus.New()
	loggerHook := logrusTest.NewLocal(testLogger)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		logError bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false, // fixed
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-services pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/img_notfound.yml",
		//},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-steps pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		{
			name:     "docker-basic steps with lazy and regular secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic_secrets.yml",
		},
		{
			name:     "kubernetes-basic steps with lazy and regular secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic_secrets.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/img_notfound.yml",
		//},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-stages pipeline with image not found",
			failure:  true,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with image not found",
		//	failure:  true, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/img_notfound.yml",
		//},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := testLogger.WithFields(logrus.Fields{"test": test.name})

			defer loggerHook.Reset()

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			_pipeline.Prepare(_build.GetRepo().GetOrg(), _build.GetRepo().GetName(), _build.GetNumber(), false)

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

			streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
			defer done()

			_engine, err := New(
				WithLogger(logger),
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithVelaClient(_client),
				WithOutputCtn(testOutputsCtn()),
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
					_runtime.(kubernetes.MockKubernetesRuntime).SimulateResync(nil)

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

			loggedError := false

			for _, logEntry := range loggerHook.AllEntries() {
				// Many errors during StreamBuild get logged and ignored.
				// So, Make sure there are no errors logged during StreamBuild.
				if logEntry.Level == logrus.ErrorLevel {
					loggedError = true

					if !test.logError {
						t.Errorf("%s StreamBuild for %s logged an Error: %v", test.name, test.pipeline, logEntry.Message)
					}
				}
			}

			if test.logError && !loggedError {
				t.Errorf("%s StreamBuild for %s did not log an Error but should have", test.name, test.pipeline)
			}
		})
	}
}

func TestLinux_StreamBuild(t *testing.T) {
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
	if err != nil {
		t.Errorf("unable to create compiler from CLI command: %v", err)
	}

	_build := testBuild()

	testLogger := logrus.New()
	loggerHook := logrusTest.NewLocal(testLogger)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	type planFuncType = func(context.Context, *pipeline.Container) error

	// planNothing is a planFuncType that does nothing
	planNothing := func(_ context.Context, _ *pipeline.Container) error {
		return nil
	}

	tests := []struct {
		name           string
		failure        bool
		earlyExecExit  bool
		earlyBuildDone bool
		logError       bool
		runtime        string
		pipeline       string
		msgCount       int
		messageKey     string
		ctn            *pipeline.Container
		streamFunc     func(*client) message.StreamFunc
		planFunc       func(*client) planFuncType
	}{
		{
			name:       "docker-basic services pipeline",
			failure:    false,
			logError:   false,
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
			logError:   false,
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
			logError:   true,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(_ *client) planFuncType {
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
			logError:   true,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(_ *client) planFuncType {
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
			logError:   false,
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
			logError:   false,
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
			logError:   true,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(_ *client) planFuncType {
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
			logError:   true,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(_ *client) planFuncType {
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
			logError:   false,
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
			logError:   false,
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
			logError:   false,
			runtime:    constants.DriverDocker,
			pipeline:   "testdata/build/secrets/basic.yml",
			messageKey: "secret",
			streamFunc: func(c *client) message.StreamFunc {
				return c.secret.stream
			},
			planFunc: func(_ *client) planFuncType {
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
			logError:   false,
			runtime:    constants.DriverKubernetes,
			pipeline:   "testdata/build/secrets/basic.yml",
			messageKey: "secret",
			streamFunc: func(c *client) message.StreamFunc {
				return c.secret.stream
			},
			planFunc: func(_ *client) planFuncType {
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
		{
			name:          "docker-early exit from ExecBuild",
			failure:       false,
			earlyExecExit: true,
			logError:      false,
			runtime:       constants.DriverDocker,
			pipeline:      "testdata/build/steps/basic.yml",
			messageKey:    "step",
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
			name:          "kubernetes-early exit from ExecBuild",
			failure:       false,
			earlyExecExit: true,
			logError:      false,
			runtime:       constants.DriverKubernetes,
			pipeline:      "testdata/build/steps/basic.yml",
			messageKey:    "step",
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
			name:           "docker-build complete before ExecBuild called",
			failure:        false,
			earlyBuildDone: true,
			logError:       false,
			runtime:        constants.DriverDocker,
			pipeline:       "testdata/build/steps/basic.yml",
			messageKey:     "step",
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
			name:           "kubernetes-build complete before ExecBuild called",
			failure:        false,
			earlyBuildDone: true,
			logError:       false,
			runtime:        constants.DriverKubernetes,
			pipeline:       "testdata/build/steps/basic.yml",
			messageKey:     "step",
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
			name:          "docker-early exit from ExecBuild and build complete signaled",
			failure:       false,
			earlyExecExit: true,
			logError:      false,
			runtime:       constants.DriverDocker,
			pipeline:      "testdata/build/steps/basic.yml",
			messageKey:    "step",
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
			name:          "kubernetes-early exit from ExecBuild and build complete signaled",
			failure:       false,
			earlyExecExit: true,
			logError:      false,
			runtime:       constants.DriverKubernetes,
			pipeline:      "testdata/build/steps/basic.yml",
			messageKey:    "step",
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buildCtx, done := context.WithCancel(context.Background())
			defer done()

			streamRequests := make(chan message.StreamRequest)

			logger := testLogger.WithFields(logrus.Fields{"test": test.name})

			defer loggerHook.Reset()

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			_pipeline.Prepare(_build.GetRepo().GetOrg(), _build.GetRepo().GetName(), _build.GetNumber(), false)

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
				WithLogger(logger),
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithLogStreamingTimeout(1*time.Second),
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

			// simulate AssembleBuild()/ExecBuild() which run concurrently with StreamBuild()
			go func() {
				// This Kubernetes setup would normally be called within AssembleBuild()
				if test.runtime == constants.DriverKubernetes {
					err = _runtime.(kubernetes.MockKubernetesRuntime).SetupMock()
					if err != nil {
						t.Errorf("Kubernetes runtime SetupMock returned err: %v", err)
					}

					// Runtime.AssembleBuild marks the PodTracker Ready before it starts it
					_runtime.(kubernetes.MockKubernetesRuntime).MarkPodTrackerReady()
				}

				if test.earlyBuildDone {
					// imitate build getting canceled or otherwise finishing before ExecBuild gets called.
					done()
				}

				if test.earlyExecExit {
					// imitate a failure after ExecBuild starts and before it sends a StreamRequest.
					close(streamRequests)
				}

				if test.earlyBuildDone || test.earlyExecExit {
					return
				}

				// simulate two messages of the same type
				for range 2 {
					// ExecBuild calls PlanService()/PlanStep() before ExecService()/ExecStep()
					// (ExecStage() calls PlanStep() before ExecStep()).
					_engine.err = test.planFunc(_engine)(buildCtx, test.ctn)

					// ExecService()/ExecStep()/secret.exec() send this message
					streamRequests <- message.StreamRequest{
						Key:    test.messageKey,
						Stream: test.streamFunc(_engine),
						// in a real pipeline, the second message would be for a different container
						Container: test.ctn,
					}

					// simulate exec build duration
					time.Sleep(100 * time.Microsecond)
				}

				// signal the end of ExecBuild so StreamBuild can finish up
				close(streamRequests)
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

			loggedError := false

			for _, logEntry := range loggerHook.AllEntries() {
				// Many errors during StreamBuild get logged and ignored.
				// So, Make sure there are no errors logged during StreamBuild.
				if logEntry.Level == logrus.ErrorLevel {
					loggedError = true

					if !test.logError {
						t.Errorf("%s StreamBuild for %s logged an Error: %v", test.name, test.pipeline, logEntry.Message)
					}
				}
			}

			if test.logError && !loggedError {
				t.Errorf("%s StreamBuild for %s did not log an Error but should have", test.name, test.pipeline)
			}
		})
	}
}

func TestLinux_DestroyBuild(t *testing.T) {
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
	if err != nil {
		t.Errorf("unable to create compiler from CLI command: %v", err)
	}

	_build := testBuild()

	testLogger := logrus.New()
	loggerHook := logrusTest.NewLocal(testLogger)

	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		logError bool
		runtime  string
		pipeline string
	}{
		{
			name:     "docker-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "kubernetes-basic secrets pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/secrets/basic.yml",
		},
		{
			name:     "docker-secrets pipeline with name not found",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/secrets/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-secrets pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/secrets/name_notfound.yml",
		//},
		{
			name:     "docker-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "kubernetes-basic services pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "docker-services pipeline with name not found",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/services/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-services pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/services/name_notfound.yml",
		//},
		{
			name:     "docker-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "kubernetes-basic steps pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "docker-steps pipeline with name not found",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/steps/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-steps pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/steps/name_notfound.yml",
		//},
		{
			name:     "docker-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "kubernetes-basic stages pipeline",
			failure:  false,
			logError: false,
			runtime:  constants.DriverKubernetes,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "docker-stages pipeline with name not found",
			failure:  false,
			logError: false,
			runtime:  constants.DriverDocker,
			pipeline: "testdata/build/stages/name_notfound.yml",
		},
		//{
		//	name:     "kubernetes-stages pipeline with name not found",
		//	failure:  false, // FIXME: make Kubernetes mock simulate failure similar to Docker mock
		//	logError: false,
		//	runtime:  constants.DriverKubernetes,
		//	pipeline: "testdata/build/stages/name_notfound.yml",
		//},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := testLogger.WithFields(logrus.Fields{"test": test.name})

			defer loggerHook.Reset()

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile %s pipeline %s: %v", test.name, test.pipeline, err)
			}

			_pipeline.Prepare(_build.GetRepo().GetOrg(), _build.GetRepo().GetName(), _build.GetNumber(), false)

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
				WithLogger(logger),
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithVelaClient(_client),
				WithOutputCtn(testOutputsCtn()),
			)
			if err != nil {
				t.Errorf("unable to create %s executor engine: %v", test.name, err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("%s unable to create build: %v", test.name, err)
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

			loggedError := false

			for _, logEntry := range loggerHook.AllEntries() {
				// Many errors during StreamBuild get logged and ignored.
				// So, Make sure there are no errors logged during StreamBuild.
				if logEntry.Level == logrus.ErrorLevel {
					// Ignore error from not mocking something in the VelaClient
					if strings.HasPrefix(logEntry.Message, "unable to upload") ||
						strings.HasSuffix(logEntry.Message, "not found") {
						// unable to upload final step state: Step 0 does not exist
						// unable to upload service snapshot: Service 0 does not exist
						// unable to destroy secret: Error: No such container: secret_github_octocat_1_notfound
						// unable to destroy service: Error: No such container: service_github_octocat_1_notfound
						// unable to destroy step: Error: No such container: github_octocat_1_test_notfound
						// unable to destroy stage: Error: No such container: github_octocat_1_test_notfound
						continue
					}

					loggedError = true

					if !test.logError {
						t.Errorf("%s StreamBuild for %s logged an Error: %v", test.name, test.pipeline, logEntry.Message)
					}
				}
			}

			if test.logError && !loggedError {
				t.Errorf("%s StreamBuild for %s did not log an Error but should have", test.name, test.pipeline)
			}
		})
	}
}

func testOutputsCtn() *pipeline.Container {
	return &pipeline.Container{
		ID:          "outputs_test",
		Environment: make(map[string]string),
		Detach:      true,
	}
}
