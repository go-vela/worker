// SPDX-License-Identifier: Apache-2.0

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
	// https://pkg.go.dev/github.com/gin-gonic/gin#New
	r := gin.New()

	// attach a middleware that recovers from any panics
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#Recovery
	r.Use(gin.Recovery())

	// attach a middleware that injects the Vela version into the request
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware#RequestVersion
	r.Use(middleware.RequestVersion)

	// attach a middleware that prevents the client from caching
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware#NoCache
	r.Use(middleware.NoCache)

	// attach a middleware capable of handling options requests
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware#Options
	r.Use(middleware.Options)

	// attach a middleware for adding extra security measures
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router/middleware#Secure
	r.Use(middleware.Secure)

	// attach all other provided middleware
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#Engine.Use
	r.Use(options...)

	// add an endpoint for reporting the health of the worker
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.GET
	r.GET("/health", api.Health)

	// add an endpoint for reporting metrics for the worker
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.GET
	r.GET("/metrics", gin.WrapH(api.Metrics()))

	// add an endpoint for reporting version information for the worker
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.GET
	r.GET("/version", api.Version)

	// add a collection of endpoints for handling API related requests
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.Group
	baseAPI := r.Group(base, perm.MustServer())
	{
		// add an endpoint for shutting down the worker
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.POST
		baseAPI.POST("/shutdown", api.Shutdown)

		// add a collection of endpoints for handling executor related requests
		//
		// https://pkg.go.dev/github.com/go-vela/worker/router#ExecutorHandlers
		ExecutorHandlers(baseAPI)
	}

	// endpoint for passing a new registration token to the deadloop running operate.go
	r.POST("/register", api.Register)

	return r
}
