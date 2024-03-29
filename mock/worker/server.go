// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FakeHandler returns an http.Handler that is capable of handling
// Vela API requests and returning mock responses.
func FakeHandler() http.Handler {
	gin.SetMode(gin.TestMode)

	e := gin.New()

	// mock endpoints for executor calls
	e.GET("/api/v1/executors", getExecutors)
	e.GET("/api/v1/executors/:executor", getExecutor)

	// mock endpoints for build calls
	e.GET("/api/v1/executors/:executor/build", getBuild)
	e.DELETE("/api/v1/executors/:executor/build/cancel", cancelBuild)

	// mock endpoints for pipeline calls
	e.GET("/api/v1/executors/:executor/pipeline", getPipeline)

	// mock endpoints for repo calls
	e.GET("/api/v1/executors/:executor/repo", getRepo)

	// mock endpoint for register call
	e.POST("/register", postRegister)

	return e
}
