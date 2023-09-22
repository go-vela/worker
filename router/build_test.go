// SPDX-License-Identifier: Apache-2.0

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
			Path:        "/build/cancel",
			Handler:     "github.com/go-vela/worker/api.CancelBuild",
			HandlerFunc: api.CancelBuild,
		},
	}

	// run test
	BuildHandlers(&_engine.RouterGroup)

	got := _engine.Routes()

	if len(got) != len(want) {
		t.Errorf("BuildHandlers is %v, want %v", got, want)
	}
}
