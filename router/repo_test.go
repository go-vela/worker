// SPDX-License-Identifier: Apache-2.0

package router

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

func TestRouter_RepoHandlers(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	_engine := gin.New()

	want := gin.RoutesInfo{
		{
			Method:      "GET",
			Path:        "/repo",
			Handler:     "github.com/go-vela/worker/api.GetRepo",
			HandlerFunc: api.GetRepo,
		},
	}

	// run test
	RepoHandlers(&_engine.RouterGroup)

	got := _engine.Routes()

	if len(got) != len(want) {
		t.Errorf("RepoHandlers is %v, want %v", got, want)
	}
}
