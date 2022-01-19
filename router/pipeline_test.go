// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

func TestRouter_PipelineHandlers(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	_engine := gin.New()

	want := gin.RoutesInfo{
		{
			Method:      "GET",
			Path:        "/pipeline",
			Handler:     "github.com/go-vela/worker/api.GetPipeline",
			HandlerFunc: api.GetPipeline,
		},
	}

	// run test
	PipelineHandlers(&_engine.RouterGroup)

	got := _engine.Routes()

	if len(got) != len(want) {
		t.Errorf("PipelineHandlers is %v, want %v", got, want)
	}
}
