// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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

	// mock endpoints for utility calls
	//e.GET("/health", getHealth)
	//e.GET("/metrics", getMetrics)
	//e.GET("/version", getVersion)
	//e.POST("/api/v1/shutdown", postShutdown)

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

	return e
}
