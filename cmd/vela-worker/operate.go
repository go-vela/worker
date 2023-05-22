// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"time"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// operate is a helper function to initiate all
// subprocesses for the operator to poll the
// queue and execute Vela pipelines.
func (w *Worker) operate(ctx context.Context) error {
	var err error

	// create the errgroup for managing operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	executors, gctx := errgroup.WithContext(ctx)

	// Define the database representation of the worker
	// and register itself in the database
	registryWorker := new(library.Worker)
	registryWorker.SetHostname(w.Config.API.Address.Hostname())
	registryWorker.SetAddress(w.Config.API.Address.String())
	registryWorker.SetRoutes(w.Config.Queue.Routes)
	registryWorker.SetActive(true)
	registryWorker.SetBuildLimit(int64(w.Config.Build.Limit))

	// pull registration token from configuration if provided; wait if not
	logrus.Trace("waiting for register token")

	token := <-w.RegisterToken

	logrus.Trace("received register token")

	// setup the vela client with the token
	w.VelaClient, err = setupClient(w.Config.Server, token)
	if err != nil {
		return err
	}

	// spawn goroutine for phoning home
	executors.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				logrus.Info("completed looping on worker registration")
				return nil
			default:
				// check in attempt loop
				for {
					// register or update the worker
					//nolint:contextcheck // ignore passing context
					w.CheckedIn, token, err = w.checkIn(registryWorker)
					// check in failed
					if err != nil {
						// check if token is expired
						expired, expiredErr := w.VelaClient.Authentication.IsTokenAuthExpired()
						if expiredErr != nil {
							logrus.Error("unable to check token expiration")
							return expiredErr
						}

						// token has expired
						if expired && len(w.Config.Server.Secret) == 0 {
							// wait on new registration token, return to check in attempt
							logrus.Trace("check-in token has expired, waiting for new register token")

							token = <-w.RegisterToken

							// setup the vela client with the token
							w.VelaClient, err = setupClient(w.Config.Server, token)
							if err != nil {
								return err
							}

							continue
						}

						// check in failed, token is still valid, retry
						logrus.Errorf("unable to check-in worker %s on the server: %v", registryWorker.GetHostname(), err)
						logrus.Info("retrying check-in...")

						time.Sleep(5 * time.Second)

						continue
					}

					// successful check in breaks the loop
					break
				}

				// setup the vela client with the token
				w.VelaClient, err = setupClient(w.Config.Server, token)
				if err != nil {
					return err
				}

				// sleep for the configured time
				time.Sleep(w.Config.CheckIn)
			}
		}
	})

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#New
	//nolint:contextcheck // ignore passing context
	w.Queue, err = queue.New(w.Config.Queue)
	if err != nil {
		registryWorker.SetStatus(constants.WorkerStatusError)
		_, res, ers := w.VelaClient.Worker.Update(registryWorker.GetHostname(), registryWorker)
		if ers != nil {
			// log the error instead of returning so the operation doesn't block worker deployment
			logrus.Errorf("status code: %v, unable to update worker %s status with the server: %w", res.StatusCode, registryWorker.GetHostname(), ers)
		}
		return err
	}

	// iterate till the configured build limit
	for i := 0; i < w.Config.Build.Limit; i++ {
		// evaluate and capture i at each iteration
		//
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		id := i

		// log a message indicating the start of an operator thread
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
		logrus.Infof("thread ID %d listening to queue...", id)

		// spawn errgroup routine for operator subprocess
		//
		// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Go
		executors.Go(func() error {
			// create an infinite loop to poll for builds
			for {
				// do not pull from queue unless worker is checked in with server
				if !w.CheckedIn {
					time.Sleep(5 * time.Second)
					logrus.Info("worker not checked in, skipping queue read")
					continue
				}
				select {
				case <-gctx.Done():
					logrus.WithFields(logrus.Fields{
						"id": id,
					}).Info("completed looping on worker executor")
					return nil
				default:
					logrus.WithFields(logrus.Fields{
						"id": id,
					}).Info("running worker executor exec")

					// exec operator subprocess to poll and execute builds
					// (do not pass the context to avoid errors in one
					// executor+build inadvertently canceling other builds)
					//nolint:contextcheck // ignore passing context
					err = w.exec(id)
					if err != nil {
						// log the error received from the executor
						//
						// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
						logrus.Errorf("failing worker executor: %v", err)
						registryWorker.SetStatus(constants.WorkerStatusError)
						_, res, ers := w.VelaClient.Worker.Update(registryWorker.GetHostname(), registryWorker)
						if ers != nil {
							// log the error instead of returning so the operation doesn't block worker deployment
							logrus.Errorf("status code: %v, unable to update worker %s status with the server: %w", res.StatusCode, registryWorker.GetHostname(), ers)
						}
						return err
					}
				}
			}
		})
	}

	// wait for errors from operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Wait
	return executors.Wait()
}
