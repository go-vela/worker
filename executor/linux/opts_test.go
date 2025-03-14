// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"
)

func TestLinux_Opt_WithBuild(t *testing.T) {
	// setup types
	_build := testBuild()

	// setup tests
	tests := []struct {
		name    string
		failure bool
		build   *api.Build
	}{
		{
			name:    "build",
			failure: false,
			build:   _build,
		},
		{
			name:    "nil build",
			failure: true,
			build:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithBuild(test.build),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithBuild returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.build, _build) {
				t.Errorf("WithBuild is %v, want %v", _engine.build, _build)
			}
		})
	}
}

func TestLinux_Opt_WithMaxLogSize(t *testing.T) {
	// setup tests
	tests := []struct {
		name       string
		failure    bool
		maxLogSize uint
	}{
		{
			name:       "defined",
			failure:    false,
			maxLogSize: 2097152,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithMaxLogSize(test.maxLogSize),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithMaxLogSize should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithMaxLogSize returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.maxLogSize, test.maxLogSize) {
				t.Errorf("WithMaxLogSize is %v, want %v", _engine.maxLogSize, test.maxLogSize)
			}
		})
	}
}

func TestLinux_Opt_WithLogStreamingTimeout(t *testing.T) {
	// setup tests
	tests := []struct {
		name                string
		failure             bool
		logStreamingTimeout time.Duration
	}{
		{
			name:                "defined",
			failure:             false,
			logStreamingTimeout: 1 * time.Second,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithLogStreamingTimeout(test.logStreamingTimeout),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogStreamingTimeout should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithLogStreamingTimeout returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.logStreamingTimeout, test.logStreamingTimeout) {
				t.Errorf("WithLogStreamingTimeout is %v, want %v", _engine.logStreamingTimeout, test.logStreamingTimeout)
			}
		})
	}
}

func TestLinux_Opt_WithPrivilegedImages(t *testing.T) {
	// setup tests
	tests := []struct {
		name             string
		failure          bool
		privilegedImages []string
	}{
		{
			name:             "empty privileged images",
			failure:          false,
			privilegedImages: []string{},
		},
		{
			name:             "with privileged image",
			failure:          false,
			privilegedImages: []string{"target/vela-docker"},
		},
		{
			name:             "with privileged images",
			failure:          false,
			privilegedImages: []string{"alpine", "target/vela-docker"},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithPrivilegedImages(test.privilegedImages),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithPrivilegedImages should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithPrivilegedImages returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.privilegedImages, test.privilegedImages) {
				t.Errorf("WithPrivilegedImages is %v, want %v", _engine.privilegedImages, test.privilegedImages)
			}
		})
	}
}

func TestLinux_Opt_WithEnforceTrustedRepos(t *testing.T) {
	// setup tests
	tests := []struct {
		name                string
		failure             bool
		enforceTrustedRepos bool
	}{
		{
			name:                "enforce trusted repos enabled",
			failure:             false,
			enforceTrustedRepos: true,
		},
		{
			name:                "enforce trusted repos disabled",
			failure:             false,
			enforceTrustedRepos: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithEnforceTrustedRepos(test.enforceTrustedRepos),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithEnforceTrustedRepos should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithEnforceTrustedRepos returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.enforceTrustedRepos, test.enforceTrustedRepos) {
				t.Errorf("WithEnforceTrustedRepos is %v, want %v", _engine.enforceTrustedRepos, test.enforceTrustedRepos)
			}
		})
	}
}

func TestLinux_Opt_WithHostname(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		hostname string
		want     string
	}{
		{
			name:     "dns hostname",
			hostname: "vela.worker.localhost",
			want:     "vela.worker.localhost",
		},
		{
			name:     "empty hostname is localhost",
			hostname: "",
			want:     "localhost",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithHostname(test.hostname),
			)
			if err != nil {
				t.Errorf("unable to create linux engine: %v", err)
			}

			if !reflect.DeepEqual(_engine.Hostname, test.want) {
				t.Errorf("WithHostname is %v, want %v", _engine.Hostname, test.want)
			}
		})
	}
}

func TestLinux_Opt_WithLogger(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		logger  *logrus.Entry
	}{
		{
			name:    "provided logger",
			failure: false,
			logger:  &logrus.Entry{},
		},
		{
			name:    "nil logger",
			failure: false,
			logger:  nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithLogger(test.logger),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogger should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithLogger returned err: %v", err)
			}

			if test.logger == nil && _engine.Logger == nil {
				t.Errorf("_engine.Logger should not be nil even if nil is passed to WithLogger")
			}

			if test.logger != nil && !reflect.DeepEqual(_engine.Logger, test.logger) {
				t.Errorf("WithLogger set %v, want %v", _engine.Logger, test.logger)
			}
		})
	}
}

func TestLinux_Opt_WithPipeline(t *testing.T) {
	// setup types
	_steps := testSteps(constants.DriverDocker)

	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			name:     "steps pipeline",
			failure:  false,
			pipeline: _steps,
		},
		{
			name:     "nil pipeline",
			failure:  true,
			pipeline: nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithPipeline(test.pipeline),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithPipeline should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithPipeline returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.pipeline, _steps) {
				t.Errorf("WithPipeline is %v, want %v", _engine.pipeline, _steps)
			}
		})
	}
}

func TestLinux_Opt_WithRuntime(t *testing.T) {
	// setup types
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
	}{
		{
			name:    "docker runtime",
			failure: false,
			runtime: _docker,
		},
		{
			name:    "kubernetes runtime",
			failure: false,
			runtime: _kubernetes,
		},
		{
			name:    "nil runtime",
			failure: true,
			runtime: nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithRuntime(test.runtime),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithRuntime should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithRuntime returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.Runtime, test.runtime) {
				t.Errorf("WithRuntime is %v, want %v", _engine.Runtime, test.runtime)
			}
		})
	}
}

func TestLinux_Opt_WithVelaClient(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		client  *vela.Client
	}{
		{
			name:    "vela client",
			failure: false,
			client:  _client,
		},
		{
			name:    "nil vela client",
			failure: true,
			client:  nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithVelaClient(test.client),
			)

			if test.failure {
				if err == nil {
					t.Errorf("WithVelaClient should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WithVelaClient returned err: %v", err)
			}

			if !reflect.DeepEqual(_engine.Vela, _client) {
				t.Errorf("WithVelaClient is %v, want %v", _engine.Vela, _client)
			}
		})
	}
}

func TestLinux_Opt_WithVersion(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version",
			version: "v1.0.0",
			want:    "v1.0.0",
		},
		{
			name:    "empty version",
			version: "",
			want:    "v0.0.0",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := New(
				WithVersion(test.version),
			)
			if err != nil {
				t.Errorf("unable to create linux engine: %v", err)
			}

			if !reflect.DeepEqual(_engine.Version, test.want) {
				t.Errorf("WithVersion is %v, want %v", _engine.Version, test.want)
			}
		})
	}
}
