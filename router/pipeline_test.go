// SPDX-License-Identifier: Apache-2.0

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
