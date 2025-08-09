// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	stdcontext "context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_ServerAddress(t *testing.T) {
	// setup types
	got := ""
	want := "foobar"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequestWithContext(stdcontext.Background(), http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(ServerAddress(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("server-address").(string)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("ServerAddress returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ServerAddress is %v, want %v", got, want)
	}
}
