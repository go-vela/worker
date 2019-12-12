// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this pipelinesitory.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
)

// pipelineHandlers is a function that extends the provided base router group
// with the API handlers for pipeline functionality.
//
// GET    	/api/v1/executors/:executor/pipeline
func pipelineHandlers(base *gin.RouterGroup) {

	// pipelines endpoints
	pipeline := base.Group("/pipeline")
	{
		pipeline.GET("", api.GetPipeline)
	} // end of pipelines endpoints
}
