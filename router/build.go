package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
)

// buildHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// GET    /api/v1/executors/:executor/builds/:build
// PUT    /api/v1/executors/:executor/builds/:build/kill
func buildHandlers(base *gin.RouterGroup) {

	// builds endpoints
	builds := base.Group("/builds")
	{

		// build endpoints
		build := builds.Group("/:build")
		{
			build.GET("", api.GetBuild)
			build.PUT("/kill", api.KillBuild)
		} // end of build endpoints

	} // end of builds endpoints
}
