// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/worker/api"
)

// RepoHandlers extends the provided base router group
// by adding a collection of endpoints for handling
// repo related requests.
//
// GET  /api/v1/executors/:executor/repo
func RepoHandlers(base *gin.RouterGroup) {
	// add a collection of endpoints for handling repo related requests
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.Group
	repo := base.Group("/repo")
	{
		// add an endpoint for capturing the repo
		//
		// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#RouterGroup.GET
		repo.GET("", api.GetRepo)
	}
}
