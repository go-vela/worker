// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/http"

	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// register is a helper function to register
// the worker in the database, updating the item
// if the worker already exists
func (w *Worker) register(config *library.Worker) {
	// check to see if the worker already exists in the database
	_, resp, err := w.VelaClient.Worker.Get(config.GetHostname())
	if resp.StatusCode == http.StatusNotFound {
		// worker does not exist, create it
		logrus.Info("registering worker with the server")
		_, _, err := w.VelaClient.Worker.Add(config)
		if err != nil {
			// log the error instead of returning so the operation doesn't block worker deployment
			logrus.Error("unable to register worker with the server: %w", err)
		}
		return
	}

	// if there was an error other than a 404, log the error.
	if err != nil {
		// log the error instead of returning so the operation doesn't block worker deployment
		logrus.Error("unable to get worker from server: %w", err)
	}

	// the worker exists in the db, update it with the new config
	logrus.Info("worker previously registered with server, updating information")
	_, _, err = w.VelaClient.Worker.Update(config.GetHostname(), config)
	if err != nil {
		// log the error instead of returning so the operation doesn't block worker deployment
		logrus.Error("unable to update worker on the server: %w", err)
	}
}
