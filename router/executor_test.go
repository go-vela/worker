// SPDX-License-Identifier: Apache-2.0

package router

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

func TestRouter_ExecutorHandlers(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	_engine := gin.New()

	want := gin.RoutesInfo{
		{
			Method:      "GET",
			Path:        "/executors",
			Handler:     "github.com/go-vela/worker/api.GetExecutors",
			HandlerFunc: api.GetExecutors,
		},
		{
			Method:      "GET",
			Path:        "/executors/:executor",
			Handler:     "github.com/go-vela/worker/api.GetExecutor",
			HandlerFunc: api.GetExecutor,
		},
		{
			Method:      "GET",
			Path:        "/executors/:executor/build",
			Handler:     "github.com/go-vela/worker/api.GetBuild",
			HandlerFunc: api.GetBuild,
		},
		{
			Method:      "DELETE",
			Path:        "/executors/:executor/build/cancel",
			Handler:     "github.com/go-vela/worker/api.CancelBuild",
			HandlerFunc: api.CancelBuild,
		},
		{
			Method:      "GET",
			Path:        "/executors/:executor/pipeline",
			Handler:     "github.com/go-vela/worker/api.GetPipeline",
			HandlerFunc: api.GetPipeline,
		},
		{
			Method:      "GET",
			Path:        "/executors/:executor/repo",
			Handler:     "github.com/go-vela/worker/api.GetRepo",
			HandlerFunc: api.GetRepo,
		},
	}

	// run test
	ExecutorHandlers(&_engine.RouterGroup)

	got := _engine.Routes()

	if len(got) != len(want) {
		t.Errorf("ExecutorHandlers is %v, want %v", got, want)
	}
}
