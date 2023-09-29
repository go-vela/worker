// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// create the context for controlling the worker subprocesses
	ctx, done := context.WithCancel(context.Background())
	// create the errgroup for managing worker subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	g, gctx := errgroup.WithContext(ctx)

	httpHandler, tlsCfg := w.server()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", w.Config.API.Address.Port()),
		Handler:           httpHandler,
		TLSConfig:         tlsCfg,
		ReadHeaderTimeout: 60 * time.Second,
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
			logrus.Info("Closing signal goroutine")
			err := server.Shutdown(ctx)
			if err != nil {
				logrus.Error(err)
			}
			return gctx.Err()
		}

		return nil
	})

	// spawn goroutine for starting the server
	g.Go(func() error {
		var err error
		logrus.Info("starting worker server")
		if tlsCfg != nil {
			if err := server.ListenAndServeTLS(w.Config.Certificate.Cert, w.Config.Certificate.Key); !errors.Is(err, http.ErrServerClosed) {
				// log a message indicating the start of the server
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
				logrus.Errorf("failing worker server: %v", err)
			}
		} else {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				// log a message indicating the start of the server
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
				logrus.Errorf("failing worker server: %v", err)
			}
		}

		return err
	})

	// spawn goroutine for starting the operator
	g.Go(func() error {
		logrus.Info("starting worker operator")
		// start the operator for the worker
		err := w.operate(gctx)
		if err != nil {
			// log the error received from the operator
			//
			// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
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
