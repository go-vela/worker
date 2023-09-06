// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/go-vela/types/library"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_WorkerRegistration(t *testing.T) {

	// setup types
	want := make(chan library.WorkerRegistration, 1)
	got := make(chan library.WorkerRegistration, 1)

	want <- library.WorkerRegistration{}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(WorkerRegistration(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("worker-registration").(chan library.WorkerRegistration)
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("QueueRegistration returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("QueueRegistration is %v, want foo", got)
	}
}
