// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
	"github.com/go-vela/worker/router/middleware"
)

const (
	base = "/api/v1"
)

// Load is a server function that returns the engine for processing web requests
// on the host it's running on
func Load(options ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.RequestVersion)
	r.Use(middleware.NoCache)
	r.Use(middleware.Options)
	r.Use(middleware.Secure)

	r.Use(options...)

	r.GET("/health", api.Health)
	r.GET("/metrics", gin.WrapH(api.Metrics()))
	r.POST("/shutdown", api.Shutdown)

	// api endpoints
	baseAPI := r.Group(base)
	{
		// executor endpoints
		executorHandlers(baseAPI)

	} // end of api

	return r
}
