// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
)

// buildHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// GET    	/api/v1/executors/:executor/build
// DELETE    /api/v1/executors/:executor/build/kill
func buildHandlers(base *gin.RouterGroup) {
	// builds endpoints
	build := base.Group("/build")
	{
		build.GET("", api.GetBuild)
		build.DELETE("/kill", api.KillBuild)
	} // end of builds endpoints
}
