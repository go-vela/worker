// SPDX-License-Identifier: Apache-2.0

package perm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPerm_MustServer_ValidateToken200(t *testing.T) {
	// setup types
	tkn := "superSecret"

	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	workerCtx.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// mocked token validation endpoint used in MustServer
	serverEngine.GET("/validate-token", func(c *gin.Context) {
		// token is not expired and matches server token
		c.Status(http.StatusOK)
	})

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", serverMock.URL) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusOK {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusOK)
	}
}

func TestPerm_MustServer_ValidateToken401(t *testing.T) {
	// setup types
	tkn := "superSecret"

	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	workerCtx.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// mocked token validation endpoint used in MustServer
	serverEngine.GET("/validate-token", func(c *gin.Context) {
		// test that validate-token returning a 401 works as expected
		c.Status(http.StatusUnauthorized)
	})

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", serverMock.URL) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusUnauthorized {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustServer_ValidateToken404(t *testing.T) {
	// setup types
	tkn := "superSecret"

	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	workerCtx.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// skip mocked token validation endpoint used in MustServer
	// test that validate-token returning a 404 works as expected

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", serverMock.URL) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusUnauthorized {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustServer_ValidateToken500(t *testing.T) {
	// setup types
	tkn := "superSecret"

	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	workerCtx.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// mocked token validation endpoint used in MustServer
	serverEngine.GET("/validate-token", func(c *gin.Context) {
		// validate-token returning a server error
		c.Status(http.StatusInternalServerError)
	})

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", serverMock.URL) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusUnauthorized {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustServer_BadServerAddress(t *testing.T) {
	// setup types
	tkn := "superSecret"
	badServerAddress := "test.example.com"

	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	workerCtx.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// mocked token validation endpoint used in MustServer
	serverEngine.GET("/validate-token", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", badServerAddress) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusInternalServerError {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusInternalServerError)
	}
}

func TestPerm_MustServer_NoToken(t *testing.T) {
	// setup types
	tkn := ""

	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	workerCtx.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// mocked token validation endpoint used in MustServer
	serverEngine.GET("/validate-token", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", serverMock.URL) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusBadRequest {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusBadRequest)
	}
}

func TestPerm_MustServer_NoAuth(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock worker router
	workerResp := httptest.NewRecorder()
	workerCtx, workerEngine := gin.CreateTestContext(workerResp)

	// fake request made to the worker router
	workerCtx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/build/cancel", nil)
	// test that skipping adding an authorization header is handled properly

	// setup mock server router
	// the URL of the mock server router is injected into the mock worker router
	serverResp := httptest.NewRecorder()
	_, serverEngine := gin.CreateTestContext(serverResp)

	// mocked token validation endpoint used in MustServer
	serverEngine.GET("/validate-token", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	serverMock := httptest.NewServer(serverEngine)
	defer serverMock.Close()

	workerEngine.Use(func(c *gin.Context) { c.Set("server-address", serverMock.URL) })

	// attach perm middleware that we are testing
	workerEngine.Use(MustServer())
	workerEngine.GET("/build/cancel", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	workerMock := httptest.NewServer(workerEngine)
	defer workerMock.Close()

	// run test
	workerEngine.ServeHTTP(workerCtx.Writer, workerCtx.Request)

	if workerResp.Code != http.StatusBadRequest {
		t.Errorf("MustServer returned %v, want %v", workerResp.Code, http.StatusBadRequest)
	}
}
