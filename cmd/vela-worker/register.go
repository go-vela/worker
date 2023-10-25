// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// checkIn is a helper function to phone home to the server.
func (w *Worker) checkIn(config *library.Worker) (bool, string, error) {
	// check to see if the worker already exists in the database
	logrus.Infof("retrieving worker %s from the server", config.GetHostname())

	_, resp, err := w.VelaClient.Worker.Get(config.GetHostname())
	if err != nil {
		respErr := fmt.Errorf("unable to retrieve worker %s from the server: %w", config.GetHostname(), err)
		// if server is down, the worker status will not be updated
		if resp == nil {
			return false, "", respErr
		}
		// if we receive a 404 the worker needs to be registered
		if resp.StatusCode == http.StatusNotFound {
			return w.register(config)
		}

		return false, "", respErr
	}

	// if we were able to GET the worker, update it
	logrus.Infof("checking worker %s into the server", config.GetHostname())

	tkn, _, err := w.VelaClient.Worker.RefreshAuth(config.GetHostname())
	if err != nil {
		// set to error when check in fails
		w.updateWorkerStatus(config, constants.WorkerStatusError)
		return false, "", fmt.Errorf("unable to refresh auth for worker %s on the server: %w", config.GetHostname(), err)
	}
	// update worker status to Idle when checkIn is successful.
	w.updateWorkerStatus(config, constants.WorkerStatusIdle)

	return true, tkn.GetToken(), nil
}

// register is a helper function to register the worker with the server.
func (w *Worker) register(config *library.Worker) (bool, string, error) {
	logrus.Infof("worker %s not found, registering it with the server", config.GetHostname())

	// status Idle will be set for worker upon first time registration
	// if worker cannot be registered, no status will be set.
	config.SetStatus(constants.WorkerStatusIdle)

	tkn, _, err := w.VelaClient.Worker.Add(config)
	if err != nil {
		// log the error instead of returning so the operation doesn't block worker deployment
		return false, "", fmt.Errorf("unable to register worker %s with the server: %w", config.GetHostname(), err)
	}

	logrus.Infof("worker %q status updated successfully to %s", config.GetHostname(), config.GetStatus())

	// successfully added the worker so return nil
	return true, tkn.GetToken(), nil
}

// queueCheckIn is a helper function to phone home to the redis.
func (w *Worker) queueCheckIn(ctx context.Context, registryWorker *library.Worker) (bool, error) {
	pErr := w.Queue.Ping(ctx)
	if pErr != nil {
		logrus.Errorf("worker %s unable to contact the queue: %v", registryWorker.GetHostname(), pErr)
		// set status to error as queue is not available
		w.updateWorkerStatus(registryWorker, constants.WorkerStatusError)

		return false, pErr
	}

	// update worker status to Idle when setup and ping are good.
	w.updateWorkerStatus(registryWorker, constants.WorkerStatusIdle)

	return true, nil
}

// updateWorkerStatus is a helper function to update worker status
// logs the error if it can't update status
func (w *Worker) updateWorkerStatus(config *library.Worker, status string) {
	config.SetStatus(status)
	_, resp, logErr := w.VelaClient.Worker.Update(config.GetHostname(), config)

	if resp == nil {
		// log the error instead of returning so the operation doesn't block worker deployment
		logrus.Error("worker status update response is nil")
	}

	if logErr != nil {
		if resp != nil {
			// log the error instead of returning so the operation doesn't block worker deployment
			logrus.Errorf("status code: %v, unable to update worker %s status with the server: %v",
				resp.StatusCode, config.GetHostname(), logErr)
		}
	}
}
