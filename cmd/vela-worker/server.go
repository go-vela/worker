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

// server is a helper function to ...
func (w *Worker) server() error {
	// create the worker router to listen and serve traffic
	router := router.Load(
		middleware.RequestVersion,
		// TODO: make this do stuff
		// middleware.Executor(w.Executors),
		middleware.Secret(w.Server.Secret),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
	)

	// start serving traffic on the provided worker port
	return router.Run(w.API.Port)
}
