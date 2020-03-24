// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"

	tomb "gopkg.in/tomb.v2"
)

// Start does stuff...
func (w *Worker) Start() error {
	// create the tomb for managing worker processes
	tomb := new(tomb.Tomb)

	// spawn a tomb goroutine to manage the worker processes
	tomb.Go(func() error {
		// spawn goroutine for starting the worker
		go func() {
			logrus.Info("starting worker server")
			// start the server for the worker
			err := w.server()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		// spawn goroutine for starting the operator
		go func() {
			logrus.Info("starting worker operator")
			// start the operator for the worker
			err := w.operate()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		for {
			select {
			case <-tomb.Dying():
				logrus.Fatal("shutting down worker")
				return tomb.Err()
			}
		}
	})

	// watch for errors from worker processes
	return tomb.Wait()
}
