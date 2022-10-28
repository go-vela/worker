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

// checkIn is a helper function to to phone home to the server.
func (w *Worker) checkIn(config *library.Worker) error {
	// refresh the server token present in the worker
	logrus.Infof("refreshing token for worker %s", config.GetHostname())

	err := w.refreshToken(config)
	if err != nil {
		return fmt.Errorf("unable to refresh token for worker %s: %w", config.GetHostname(), err)
	}

	// check to see if the worker already exists in the database
	logrus.Infof("retrieving worker %s from the server", config.GetHostname())

	_, resp, err := w.VelaClient.Worker.Get(config.GetHostname())
	if err != nil {
		respErr := fmt.Errorf("unable to retrieve worker %s from the server: %w", config.GetHostname(), err)
		if resp == nil {
			return respErr
		}
		// if we receive a 404 the worker needs to be registered
		if resp.StatusCode == http.StatusNotFound {
			return w.register(config)
		}

		return respErr
	}

	// if we were able to GET the worker, update it
	logrus.Infof("checking worker %s into the server", config.GetHostname())

	_, _, err = w.VelaClient.Worker.Update(config.GetHostname(), config)
	if err != nil {
		return fmt.Errorf("unable to update worker %s on the server: %w", config.GetHostname(), err)
	}

	return nil
}

// register is a helper function to register the worker with the server.
func (w *Worker) register(config *library.Worker) error {
	logrus.Infof("worker %s not found, registering it with the server", config.GetHostname())

	// add worker to db
	_, _, err := w.VelaClient.Worker.Add(config)
	if err != nil {
		return fmt.Errorf("unable to register worker %s with the server: %w", config.GetHostname(), err)
	}

	// successfully added the worker so return nil
	return nil
}

// refreshToken is a helper function to refresh the token with the server for accessing the server.
func (w *Worker) refreshToken(config *library.Worker) error {
	// refresh server token
	t, _, err := w.VelaClient.Token.Refresh(w.Config.ServerToken)
	if err != nil {
		return fmt.Errorf("unable to refresh token for worker %s: %w", config.GetHostname(), err)
	}

	w.Config.ServerToken = t

	// successfully refreshed the token so return nil
	return nil
}
