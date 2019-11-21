package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
)

// executorHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// GET    /api/v1/executors         --> github.com/go-vela/worker/api.GetExecutors (8 handlers)
// GET    /api/v1/executors/:executor --> github.com/go-vela/worker/api.GetExecutor (8 handlers)
// GET    /api/v1/executors/:executor/builds/:build --> github.com/go-vela/worker/api.GetBuild (8 handlers)
// PUT    /api/v1/executors/:executor/builds/:build/kill --> github.com/go-vela/worker/api.KillBuild (8 handlers)
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
