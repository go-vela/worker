// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
)

// executorHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// GET    	/api/v1/executors
// GET    	/api/v1/executors/:executor
// GET   	/api/v1/executors/:executor/builds/:build
// PATCH    /api/v1/executors/:executor/builds/:build/kill
func executorHandlers(base *gin.RouterGroup) {

	// executors endpoints
	executors := base.Group("/executors")
	{

		executors.GET("", api.GetExecutors)

		// exector endpoints
		executor := executors.Group("/:executor")
		{
			executor.GET("", api.GetExecutor)

			// build endpoints
			buildHandlers(executor)

		} // end of executor endpoints

	} // end of executors endpoints
}
