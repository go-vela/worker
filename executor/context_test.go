// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/mock/server"

	"github.com/go-vela/worker/executor/linux"

	"github.com/go-vela/worker/runtime/docker"

	"github.com/go-vela/sdk-go/vela"
)

func TestExecutor_FromContext(t *testing.T) {
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

	_engine, err := linux.New(
		linux.WithBuild(_build),
		linux.WithPipeline(_pipeline),
		linux.WithRepo(_repo),
		linux.WithRuntime(_runtime),
		linux.WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create linux engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		context context.Context
		want    Engine
	}{
		{
			name: "valid executor in context",
			//nolint:staticcheck,revive // ignore using string with context value
			context: context.WithValue(context.Background(), key, _engine),
			want:    _engine,
		},
		{
			name:    "executor not in context",
			context: context.Background(),
			want:    nil,
		},
		{
			name: "invalid executor in context",
			//nolint:staticcheck,revive // ignore using string with context value
			context: context.WithValue(context.Background(), key, "foo"),
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FromContext(test.context)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("FromContext is %v, want %v", got, test.want)
			}
		})
	}
}

func TestExecutor_FromGinContext(t *testing.T) {
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

	_engine, err := linux.New(
		linux.WithBuild(_build),
		linux.WithPipeline(_pipeline),
		linux.WithRepo(_repo),
		linux.WithRuntime(_runtime),
		linux.WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create linux engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		context *gin.Context
		value   interface{}
		want    Engine
	}{
		{
			name:    "valid executor in context",
			context: new(gin.Context),
			value:   _engine,
			want:    _engine,
		},
		{
			name:    "executor not in context",
			context: new(gin.Context),
			value:   nil,
			want:    nil,
		},
		{
			name:    "invalid executor in context",
			context: new(gin.Context),
			value:   "foo",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.value != nil {
				test.context.Set(key, test.value)
			}

			got := FromGinContext(test.context)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("FromGinContext is %v, want %v", got, test.want)
			}
		})
	}
}

func TestExecutor_WithContext(t *testing.T) {
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

	_engine, err := linux.New(
		linux.WithBuild(_build),
		linux.WithPipeline(_pipeline),
		linux.WithRepo(_repo),
		linux.WithRuntime(_runtime),
		linux.WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create linux engine: %v", err)
	}

	//nolint:staticcheck,revive // ignore using string with context value
	want := context.WithValue(context.Background(), key, _engine)

	// run test
	got := WithContext(context.Background(), _engine)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithContext is %v, want %v", got, want)
	}
}

func TestExecutor_WithGinContext(t *testing.T) {
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

	_engine, err := linux.New(
		linux.WithBuild(_build),
		linux.WithPipeline(_pipeline),
		linux.WithRepo(_repo),
		linux.WithRuntime(_runtime),
		linux.WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create linux engine: %v", err)
	}

	want := new(gin.Context)
	want.Set(key, _engine)

	// run test
	got := new(gin.Context)
	WithGinContext(got, _engine)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithGinContext is %v, want %v", got, want)
	}
}
