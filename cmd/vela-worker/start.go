// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"

	tomb "gopkg.in/tomb.v2"
)

// Start initiates all subprocesses for the Worker
// from the provided configuration. The server
// subprocess enables the Worker to listen and
// serve traffic for web and API requests. The
// operator subprocess enables the Worker to
// poll the queue and execute Vela pipelines.
func (w *Worker) Start() error {
	// create the tomb for managing worker subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb
	tomb := new(tomb.Tomb)

	// spawn a tomb goroutine to manage the worker subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Go
	tomb.Go(func() error {
		// spawn goroutine for starting the server
		go func() {
			// log a message indicating the start of the server
			//
			// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
			logrus.Info("starting worker server")

			// start the server for the worker
			err := w.server()
			if err != nil {
				// log the error received from the server
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
				logrus.Errorf("failing worker server: %v", err)

				// kill the worker subprocesses
				//
				// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Kill
				tomb.Kill(err)
			}
		}()

		// spawn goroutine for starting the operator
		go func() {
			// log a message indicating the start of the operator
			//
			// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
			logrus.Info("starting worker operator")

			// start the operator for the worker
			err := w.operate()
			if err != nil {
				// log the error received from the operator
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
				logrus.Errorf("failing worker operator: %v", err)

				// kill the worker subprocesses
				//
				// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Kill
				tomb.Kill(err)
			}
		}()

		// create an infinite loop to poll for errors
		for {
			// create a select statement to check for errors
			select {
			// check if one of the worker subprocesses died
			case <-tomb.Dying():
				// fatally log that we're shutting down the worker
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Fatal
				logrus.Fatal("shutting down worker")

				return tomb.Err()
			}
		}
	})

	// wait for errors from worker subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Wait
	return tomb.Wait()
}
