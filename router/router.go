// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// Package router Vela worker
//
// API for a Vela worker
//
//	Version: 0.0.0-dev
//	Schemes: http, https
//	Host: localhost
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	  ApiKeyAuth:
//	    description: Bearer token
//	    type: apiKey
//	    in: header
//	    name: Authorization
//
// swagger:meta
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
	"github.com/go-vela/worker/router/middleware"
	"github.com/go-vela/worker/router/middleware/perm"
	"github.com/go-vela/worker/router/middleware/user"
)

const (
	base = "/api/v1"
)

// Load creates the gin.Engine with the provided
// options (middleware functions) for processing
// web and API requests for the worker.
func Load(options ...gin.HandlerFunc) *gin.Engine {
	// create an empty gin engine with no middleware
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#New
	r := gin.New()

	// attach a middleware that recovers from any panics
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Recovery
	r.Use(gin.Recovery())

	// attach a middleware that injects the Vela version into the request
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware?tab=doc#RequestVersion
	r.Use(middleware.RequestVersion)

	// attach a middleware that prevents the client from caching
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware?tab=doc#NoCache
	r.Use(middleware.NoCache)

	// attach a middleware capable of handling options requests
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware?tab=doc#Options
	r.Use(middleware.Options)

	// attach a middleware for adding extra security measures
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware?tab=doc#Secure
	r.Use(middleware.Secure)

	// attach all other provided middleware
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Use
	r.Use(options...)

	// add an endpoint for reporting the health of the worker
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
	r.GET("/health", api.Health)

	// add an endpoint for reporting metrics for the worker
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
	r.GET("/metrics", gin.WrapH(api.Metrics()))

	// add an endpoint for reporting version information for the worker
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
	r.GET("/version", api.Version)

	// add a collection of endpoints for handling API related requests
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.Group
	baseAPI := r.Group(base, user.Establish(), perm.MustServer())
	{
		// add an endpoint for shutting down the worker
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.POST
		baseAPI.POST("/shutdown", api.Shutdown)

		// add a collection of endpoints for handling executor related requests
		//
		// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#ExecutorHandlers
		ExecutorHandlers(baseAPI)
	}

	// endpoint for passing a new registration token to the deadloop running operate.go
	r.POST("/register", api.Register)

	return r
}
