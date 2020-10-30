// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// register is a helper function to register
// the worker in the database, updating the item
// if the worker already exists
func (w *Worker) register(config *library.Worker) error {
	// check to see if the worker already exists in the database
	_, _, err := w.VelaClient.Worker.Get(config.GetHostname())
	if err != nil {
		// worker does not exist, create it
		logrus.Info("registering worker with the server")
		_, _, err := w.VelaClient.Worker.Add(config)
		if err != nil {
			return err
		}
		return nil
	}

	// the worker exists in the db, update it with the new config
	logrus.Info("worker previously registered with server, updating information")
	_, _, err = w.VelaClient.Worker.Update(config.GetHostname(), config)
	if err != nil {
		return err
	}

	return nil
}
