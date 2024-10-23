// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// checkIn is a helper function to phone home to the server.
func (w *Worker) checkIn(config *api.Worker) (bool, string, error) {
	// check to see if the worker already exists in the database
	logrus.Infof("retrieving worker %s from the server", config.GetHostname())

	var (
		tkn     *api.Token
		retries = 3
	)

	for i := 0; i < retries; i++ {
		logrus.Debugf("check in loop - attempt %d", i+1)
		// check if we're on the first iteration of the loop
		if i > 0 {
			// incrementally sleep in between retries
			time.Sleep(time.Duration(i*10) * time.Second)
		}

		_, resp, err := w.VelaClient.Worker.Get(config.GetHostname())
		if err != nil {
			respErr := fmt.Errorf("unable to retrieve worker %s from the server: %w", config.GetHostname(), err)
			// if server is down, the worker status will not be updated
			if resp == nil || (resp.StatusCode != http.StatusNotFound) {
				return false, "", respErr
			}
			// if we receive a 404 the worker needs to be registered
			if resp.StatusCode == http.StatusNotFound {
				registered, strToken, regErr := w.register(config)
				if regErr != nil {
					if i < retries-1 {
						logrus.WithError(err).Warningf("retrying #%d", i+1)

						// continue to the next iteration of the loop
						continue
					}
				}

				return registered, strToken, regErr
			}
		}

		// if we were able to GET the worker, update it
		logrus.Infof("checking worker %s into the server", config.GetHostname())

		tkn, _, err = w.VelaClient.Worker.RefreshAuth(config.GetHostname())
		if err != nil {
			if i < retries-1 {
				logrus.WithError(err).Warningf("retrying #%d", i+1)

				// continue to the next iteration of the loop
				continue
			}

			// set to error when check in fails
			w.updateWorkerStatus(config, constants.WorkerStatusError)

			return false, "", fmt.Errorf("unable to refresh auth for worker %s on the server: %w", config.GetHostname(), err)
		}

		status := w.getWorkerStatusFromConfig(config)

		// update worker status to Idle when checkIn is successful.
		w.updateWorkerStatus(config, status)

		break
	}

	return true, tkn.GetToken(), nil
}

// register is a helper function to register the worker with the server.
func (w *Worker) register(config *api.Worker) (bool, string, error) {
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
func (w *Worker) queueCheckIn(ctx context.Context, registryWorker *api.Worker) (bool, error) {
	pErr := w.Queue.Ping(ctx)
	if pErr != nil {
		logrus.Errorf("worker %s unable to contact the queue: %v", registryWorker.GetHostname(), pErr)
		// set status to error as queue is not available
		w.updateWorkerStatus(registryWorker, constants.WorkerStatusError)

		return false, pErr
	}

	status := w.getWorkerStatusFromConfig(registryWorker)

	// update worker status to Idle when setup and ping are good.
	w.updateWorkerStatus(registryWorker, status)

	return true, nil
}

// updateWorkerStatus is a helper function to update worker status
// logs the error if it can't update status.
func (w *Worker) updateWorkerStatus(config *api.Worker, status string) {
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

		if resp == nil {
			// log the error instead of returning so the operation doesn't block worker deployment
			logrus.Errorf("worker status update response is nil, unable to update worker %s status with the server: %v",
				config.GetHostname(), logErr)
		}
	}
}
