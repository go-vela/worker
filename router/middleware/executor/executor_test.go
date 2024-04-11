// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime/docker"
)

func TestExecutor_Retrieve(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	_repo := new(api.Repo)
	_repo.SetOwner(new(library.User))

	want, err := executor.New(&executor.Setup{
		Driver:     constants.DriverLinux,
		MaxLogSize: 2097152,
		Client:     new(vela.Client),
		Runtime:    _runtime,
		Build:      new(library.Build),
		Pipeline:   new(pipeline.Build),
		Repo:       _repo,
	})
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	// setup context
	context := new(gin.Context)
	executor.WithGinContext(context, want)

	// run test
	got := Retrieve(context)

	if got != want {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestExecutor_Establish(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	_repo := new(api.Repo)
	_repo.SetOwner(new(library.User))

	want, err := executor.New(&executor.Setup{
		Driver:     constants.DriverLinux,
		MaxLogSize: 2097152,
		Client:     new(vela.Client),
		Runtime:    _runtime,
		Build:      new(library.Build),
		Pipeline:   new(pipeline.Build),
		Repo:       _repo,
	})
	if err != nil {
		t.Errorf("unable to create executor engine: %v", err)
	}

	_executors := make(map[int]executor.Engine)
	_executors[0] = want

	got := *new(executor.Engine)

	// setup context
	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/0", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("executors", _executors) })
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

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Establish is %v, want %v", got, want)
	}
}

func TestExecutor_Establish_NoParam(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	// setup context
	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/", nil)

	// setup mock server
	engine.Use(Establish())
	engine.GET("/executors/:executor", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestExecutor_Establish_InvalidParam(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	// setup context
	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/foo", nil)

	// setup mock server
	engine.Use(Establish())
	engine.GET("/executors/:executor", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestExecutor_Establish_NoExecutors(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	// setup context
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
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusInternalServerError)
	}
}

func TestExecutor_Establish_InvalidExecutors(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	// setup context
	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/0", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("executors", "invalid") })
	engine.Use(Establish())
	engine.GET("/executors/:executor", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusInternalServerError)
	}
}

func TestExecutor_Establish_ExecutorNotFound(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	// setup context
	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/executors/0", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("executors", make(map[int]executor.Engine)) })
	engine.Use(Establish())
	engine.GET("/executors/:executor", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}
