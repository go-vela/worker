// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"bytes"
	stdcontext "context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestMiddleware_Logger(t *testing.T) {
	// setup types
	payload, _ := json.Marshal(`{"foo": "bar"}`)
	wantLevel := logrus.InfoLevel
	wantMessage := ""

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequestWithContext(stdcontext.Background(), http.MethodPost, "/foobar", bytes.NewBuffer(payload))

	// setup mock server
	engine.Use(Payload())
	engine.Use(Logger(logger, time.RFC3339, true))
	engine.POST("/foobar", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := hook.LastEntry().Level
	gotMessage := hook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Logger Message is %v, want %v", gotMessage, wantMessage)
	}
}

func TestMiddleware_Logger_Error(t *testing.T) {
	// setup types
	wantLevel := logrus.ErrorLevel
	wantMessage := "Error #01: test error\n"

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequestWithContext(stdcontext.Background(), http.MethodGet, "/foobar", nil)

	// setup mock server
	engine.Use(Logger(logger, time.RFC3339, true))
	engine.GET("/foobar", func(c *gin.Context) {
		//nolint:errcheck // ignore checking error
		c.Error(fmt.Errorf("test error"))
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := hook.LastEntry().Level
	gotMessage := hook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Logger Message is %v, want %v", gotMessage, wantMessage)
	}
}
