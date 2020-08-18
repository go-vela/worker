// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
	"github.com/go-vela/worker/router/middleware/executor"
)

// ExecutorHandlers extends the provided base router group
// by adding a collection of endpoints for handling
// executor related requests.
//
// GET     /api/v1/executors
// GET     /api/v1/executors/:executor
// GET     /api/v1/executors/:executor/build
// DELETE  /api/v1/executors/:executor/build/cancel
// GET     /api/v1/executors/:executor/pipeline
// GET     /api/v1/executors/:executor/repo
func ExecutorHandlers(base *gin.RouterGroup) {
	// add a collection of endpoints for handling executors related requests
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.Group
	executors := base.Group("/executors")
	{
		// add an endpoint for capturing the executors
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
		executors.GET("", api.GetExecutors)

		// add a collection of endpoints for handling executor related requests
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.Group
		executor := executors.Group("/:executor", executor.Establish())
		{
			// add an endpoint for capturing a executor
			//
			// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
			executor.GET("", api.GetExecutor)

			// add a collection of endpoints for handling build related requests
			//
			// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#BuildHandlers
			BuildHandlers(executor)

			// add a collection of endpoints for handling pipeline related requests
			//
			// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#PipelineHandlers
			PipelineHandlers(executor)

			// add a collection of endpoints for handling repo related requests
			//
			// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#RepoHandlers
			RepoHandlers(executor)
		}
	}
}
