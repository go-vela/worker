// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/worker/router"
	"github.com/go-vela/worker/router/middleware"

	"github.com/sirupsen/logrus"
)

// server is a helper function to listen and serve
// traffic for web and API requests for the Worker.
func (w *Worker) server() error {
	// create the worker router to listen and serve traffic
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#Load
	_server := router.Load(
		middleware.RequestVersion,
		middleware.Executors(w.Executors),
		middleware.Secret(w.Config.Server.Secret),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
	)

	// start serving traffic on the provided worker port
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Run
	return _server.Run(w.Config.API.Port)
}
