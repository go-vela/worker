// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os"
	"strings"
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
	logrus.Tracef("serving traffic on %s", w.Config.API.Address.Port())

	// if running with HTTPS, check certs are provided and run with TLS.
	if strings.EqualFold(w.Config.API.Address.Scheme, "https") {
		if len(w.Config.Certificate.Cert) > 0 && len(w.Config.Certificate.Key) > 0 {
			// check that the certificate and key are both populated
			_, err := os.Stat(w.Config.Certificate.Cert)
			if err != nil {
				logrus.Fatalf("Expecting certificate file at %s, got %v", w.Config.Certificate.Cert, err)
			}
			_, err = os.Stat(w.Config.Certificate.Key)
			if err != nil {
				logrus.Fatalf("Expecting certificate key at %s, got %v", w.Config.Certificate.Key, err)
			}
		} else {
			logrus.Fatal("Unable to run with TLS: No certificate provided")
		}
		return _server.RunTLS(
			fmt.Sprintf(":%s", w.Config.API.Address.Port()),
			w.Config.Certificate.Cert,
			w.Config.Certificate.Key,
		)
	}

	// else serve over http
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Run
	return _server.Run(fmt.Sprintf(":%s", w.Config.API.Address.Port()))
}
