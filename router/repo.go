// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/api"
)

// repoHandlers is a function that extends the provided base router group
// with the API handlers for repo functionality.
//
// GET    	/api/v1/executors/:executor/repo
func repoHandlers(base *gin.RouterGroup) {
	// repos endpoints
	repo := base.Group("/repo")
	{
		repo.GET("", api.GetRepo)
	} // end of repos endpoints
}
