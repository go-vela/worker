// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/worker/router/middleware/user"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestPerm_MustServer_success(t *testing.T) {
	// setup types
	secret := "superSecret"

	u := new(library.User)
	u.SetID(1)
	u.SetName("vela-server")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/server/users", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(user.Establish())
	engine.Use(MustServer())
	engine.GET("/server/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustServer returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustServer_failure(t *testing.T) {
	// setup types
	secret := "foo"

	u := new(library.User)
	u.SetID(1)
	u.SetName("not-vela-server")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/server/users", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(user.Establish())
	engine.Use(MustServer())
	engine.GET("/server/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustServer returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}
