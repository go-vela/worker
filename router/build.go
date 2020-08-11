// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
// DELETE  /api/v1/executors/:executor/build/kill
func BuildHandlers(base *gin.RouterGroup) {
	// add a collection of endpoints for handling build related requests
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.Group
	build := base.Group("/build")
	{
		// add an endpoint for capturing the build
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
		build.GET("", api.GetBuild)

		// add an endpoint for killing the build
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.DELETE
		build.DELETE("/kill", api.CancelBuild)
	}
}
