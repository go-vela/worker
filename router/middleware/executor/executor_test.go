// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/executor/linux"
)

func TestExecutor_Retrieve(t *testing.T) {
	// setup types
	want, _ := linux.New(nil, nil)

	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := Retrieve(context)

	if got != want {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestExecutor_Establish(t *testing.T) {
	// setup types
	want := make(map[int]executor.Engine)
	want[0], _ = linux.New(nil, nil)
	got := want[0]

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/0", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { executor.ToContext(c, want) })
	engine.Use(Establish())
	engine.GET("/executors/:executor", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want[0]) {
		t.Errorf("Establish is %v, want %v", got, want)
	}
}

func TestExecutor_Establish_NoExecutor(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/0", nil)

	// setup mock server
	engine.Use(Establish())
	engine.GET("/executors/:executor", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}
}
