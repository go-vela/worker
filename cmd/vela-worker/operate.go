// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/types/constants"
)

// operate is a helper function to initiate all
// subprocesses for the operator to poll the
// queue and execute Vela pipelines.
//
//nolint:funlen // refactor candidate
func (w *Worker) operate(ctx context.Context) error {
	var err error
	// create the errgroup for managing operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#Group
	executors, gctx := errgroup.WithContext(ctx)
	// Define the database representation of the worker
	// and register itself in the database
	registryWorker := new(api.Worker)
	registryWorker.SetHostname(w.Config.API.Address.Hostname())
	registryWorker.SetAddress(w.Config.API.Address.String())
	registryWorker.SetActive(true)
	registryWorker.SetBuildLimit(int64(w.Config.Build.Limit))

	// set routes from config if set or defaulted to `vela`
	if (len(w.Config.Queue.Routes) > 0) && (w.Config.Queue.Routes[0] != "NONE" && w.Config.Queue.Routes[0] != "") {
		registryWorker.SetRoutes(w.Config.Queue.Routes)
	}

	// pull registration token from configuration if provided; wait if not
	logrus.Trace("waiting for register token")

	token := <-w.RegisterToken

	logrus.Trace("received register token")
	logrus.Trace("setting up vela client")
	// setup the vela client with the token
	w.VelaClient, err = setupClient(w.Config.Server, token)
	if err != nil {
		return err
	}

	logrus.Trace("getting queue creds")
	// fetching queue credentials using registration token
	creds, _, err := w.VelaClient.Queue.GetInfo()
	if err != nil {
		logrus.Trace("error getting creds")
		return err
	}

	// if an address was given at start up, use that — else use what is returned from server
	if len(w.Config.Queue.Address) == 0 {
		w.Config.Queue.Address = creds.GetQueueAddress()
	}

	// set public key in queue config
	w.Config.Queue.PublicKey = creds.GetPublicKey()

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue#New
	w.Queue, err = queue.New(w.Config.Queue)
	if err != nil {
		logrus.Error("queue setup failed")
		// set to error as queue setup fails
		w.updateWorkerStatus(registryWorker, constants.WorkerStatusError)
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

					w.QueueCheckedIn, err = w.queueCheckIn(gctx, registryWorker)

					if err != nil {
						// queue check in failed, retry
						logrus.Errorf("unable to ping queue %v", err)
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

	// iterate till the configured build limit
	for i := 0; i < w.Config.Build.Limit; i++ {
		// evaluate and capture i at each iteration
		//
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		id := i

		// log a message indicating the start of an operator thread
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Info
		logrus.Infof("thread ID %d listening to queue...", id)

		// spawn errgroup routine for operator subprocess
		//
		// https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Go
		executors.Go(func() error {
			// create an infinite loop to poll for builds
			for {
				// do not pull from queue unless worker is checked in with server
				if !w.CheckedIn {
					time.Sleep(5 * time.Second)
					logrus.Info("worker not checked in, skipping queue read")

					continue
				}
				// do not pull from queue unless queue setup is done and connected
				if !w.QueueCheckedIn {
					time.Sleep(5 * time.Second)
					logrus.Info("queue ping failed, skipping queue read")

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
					err = w.exec(id, registryWorker)
					if err != nil {
						// log the error received from the executor
						//
						// https://pkg.go.dev/github.com/sirupsen/logrus#Errorf
						logrus.Errorf("failing worker executor: %v", err)
						registryWorker.SetStatus(constants.WorkerStatusError)
						_, resp, logErr := w.VelaClient.Worker.Update(registryWorker.GetHostname(), registryWorker)

						if resp == nil {
							// log the error instead of returning so the operation doesn't block worker deployment
							logrus.Error("status update response is nil")
						}

						if logErr != nil {
							if resp != nil {
								// log the error instead of returning so the operation doesn't block worker deployment
								logrus.Errorf("status code: %v, unable to update worker %s status with the server: %v", resp.StatusCode, registryWorker.GetHostname(), logErr)
							}
						}

						return err
					}
				}
			}
		})
	}

	// wait for errors from operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Wait
	return executors.Wait()
}
