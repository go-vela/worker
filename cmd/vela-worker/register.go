// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/http"

	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// checkIn is a helper function to phone home to the server.
func (w *Worker) checkIn(config *library.Worker) (bool, error) {
	// check to see if the worker already exists in the database
	logrus.Infof("retrieving worker %s from the server", config.GetHostname())

	_, resp, err := w.VelaClient.Worker.Get(config.GetHostname())
	if err != nil {
		respErr := fmt.Errorf("unable to retrieve worker %s from the server: %w", config.GetHostname(), err)
		if resp == nil {
			return false, respErr
		}
		// if we receive a 404 the worker needs to be registered
		if resp.StatusCode == http.StatusNotFound {
			return w.register(config)
		}

		return false, respErr
	}

	// if we were able to GET the worker, update it
	logrus.Infof("checking worker %s into the server", config.GetHostname())

	wrkCheckIn, _, err := w.VelaClient.Worker.Update(config.GetHostname(), config)
	if err != nil {
		return false, fmt.Errorf("unable to update worker %s on the server: %w", config.GetHostname(), err)
	}

	// write the token to the auth token channel
	w.AuthToken <- wrkCheckIn.Token.GetToken()

	return true, nil
}

// register is a helper function to register the worker with the server.
func (w *Worker) register(config *library.Worker) (bool, error) {
	logrus.Infof("worker %s not found, registering it with the server", config.GetHostname())

	tkn, _, err := w.VelaClient.Worker.Add(config)
	if err != nil {
		// log the error instead of returning so the operation doesn't block worker deployment
		return false, fmt.Errorf("unable to register worker %s with the server: %w", config.GetHostname(), err)
	}

	logrus.Infof("register function TOKEN from add worker: %s", tkn)

	// write the token to the auth token channel
	w.AuthToken <- tkn.GetToken()

	// successfully added the worker so return nil
	return true, nil
}
