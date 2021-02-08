// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Start initiates all subprocesses for the Worker
// from the provided configuration. The server
// subprocess enables the Worker to listen and
// serve traffic for web and API requests. The
// operator subprocess enables the Worker to
// poll the queue and execute Vela pipelines.
func (w *Worker) Start() error {
	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", w.Config.API.Address.Port()),
		Handler: w.server(),
	}

	// goroutine to check for signals to gracefully finish all functions
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			logrus.Infof("Received signal: %s", sig)
			err := server.Shutdown(ctx)
			if err != nil {
				logrus.Error(err)
			}
			done()
		case <-gctx.Done():
			logrus.Info("closing signal goroutine")
			return gctx.Err()
		}

		return nil
	})

	g.Go(func() error {
		var err error
		logrus.Info("starting worker server")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Errorf("failing worker server: %v", err)
		}

		return err
	})

	g.Go(func() error {
		logrus.Info("starting worker operator")
		err := w.operate(gctx)
		if err != nil {
			logrus.Errorf("failing worker operator: %v", err)
			return err
		}

		return err
	})

	// wait for errors from worker subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Wait
	return g.Wait()
}
