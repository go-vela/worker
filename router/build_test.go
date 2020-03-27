// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

func TestRouter_BuildHandlers(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	_engine := gin.New()

	want := gin.RoutesInfo{
		{
			Method:      "GET",
			Path:        "/build",
			Handler:     "github.com/go-vela/worker/api.GetBuild",
			HandlerFunc: api.GetBuild,
		},
		{
			Method:      "DELETE",
			Path:        "/build/kill",
			Handler:     "github.com/go-vela/worker/api.KillBuild",
			HandlerFunc: api.KillBuild,
		},
	}

	// run test
	BuildHandlers(&_engine.RouterGroup)

	got := _engine.Routes()

	if len(got) != len(want) {
		t.Errorf("BuildHandlers is %v, want %v", got, want)
	}
}
