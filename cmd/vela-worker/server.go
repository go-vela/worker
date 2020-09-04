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
	// log a message indicating the setup of the server handlers
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Trace
	logrus.Trace("loading router with server handlers")

	// create the worker router to listen and serve traffic
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#Load
	_server := router.Load(
		middleware.RequestVersion,
		middleware.Executors(w.Executors),
		middleware.Secret(w.Config.Server.Secret),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
	)

	// log a message indicating the start of serving traffic
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Tracef
	logrus.Tracef("serving traffic on %s", w.Config.API.Port)

	// start serving traffic with TLS on the provided worker port
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.RunTLS
	if len(w.Config.Certificate.Cert) > 0 {
		return _server.RunTLS(w.Config.API.Port, w.Config.Certificate.Cert, w.Config.Certificate.Key)
	}

	// if no certs are provided, run without TLS
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Run
	return _server.Run(w.Config.API.Port)
}
