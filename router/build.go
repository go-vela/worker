// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

// BuildHandlers extends the provided base router group
// by adding a collection of endpoints for handling
// build related requests.
//
// GET     /api/v1/executors/:executor/build
// DELETE  /api/v1/executors/:executor/build/cancel
// .
func BuildHandlers(base *gin.RouterGroup) {
	// add a collection of endpoints for handling build related requests
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.Group
	build := base.Group("/build")
	{
		// add an endpoint for capturing the build
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.GET
		build.GET("", api.GetBuild)

		// add an endpoint for canceling the build
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.DELETE
		build.DELETE("/cancel", api.CancelBuild)
	}
}
