// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
	"github.com/go-vela/worker/router/middleware/executor"
)

// executorHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// GET    /api/v1/executors
// GET    /api/v1/executors/:executor
// GET    /api/v1/executors/:executor/build
// DELETE /api/v1/executors/:executor/build/kill
// GET    /api/v1/executors/:executor/pipeline
// GET    /api/v1/executors/:executor/repo
func executorHandlers(base *gin.RouterGroup) {
	// executors endpoints
	executors := base.Group("/executors")
	{
		executors.GET("", api.GetExecutors)

		// exector endpoints
		executor := executors.Group("/:executor", executor.Establish())
		{
			executor.GET("", api.GetExecutor)

			// build endpoints
			buildHandlers(executor)

			// pipeline endpoints
			pipelineHandlers(executor)

			// build endpoints
			repoHandlers(executor)
		} // end of executor endpoints
	} // end of executors endpoints
}
