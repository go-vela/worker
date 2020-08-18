// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

func TestRouter_Load(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	want := gin.RoutesInfo{
		{
			Method:      "GET",
			Path:        "/health",
			Handler:     "github.com/go-vela/worker/api.Health",
			HandlerFunc: api.Health,
		},
		{
			Method:      "GET",
			Path:        "/metrics",
			Handler:     "github.com/go-vela/worker/api.Metrics",
			HandlerFunc: gin.WrapH(api.Metrics()),
		},
		{
			Method:      "POST",
			Path:        "/api/v1/shutdown",
			Handler:     "github.com/go-vela/worker/api.Shutdown",
			HandlerFunc: api.Shutdown,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/executors",
			Handler:     "github.com/go-vela/worker/api.GetExecutors",
			HandlerFunc: api.GetExecutors,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/executors/:executor",
			Handler:     "github.com/go-vela/worker/api.GetExecutor",
			HandlerFunc: api.GetExecutor,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/executors/:executor/build",
			Handler:     "github.com/go-vela/worker/api.GetBuild",
			HandlerFunc: api.GetBuild,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/executors/:executor/build/cancel",
			Handler:     "github.com/go-vela/worker/api.CancelBuild",
			HandlerFunc: api.CancelBuild,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/executors/:executor/pipeline",
			Handler:     "github.com/go-vela/worker/api.GetPipeline",
			HandlerFunc: api.GetPipeline,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/executors/:executor/repo",
			Handler:     "github.com/go-vela/worker/api.GetRepo",
			HandlerFunc: api.GetRepo,
		},
	}

	// run test
	got := Load()

	if len(got.Routes()) != len(want) {
		t.Errorf("Load is %v, want %v", got.Routes(), want)
	}
}
