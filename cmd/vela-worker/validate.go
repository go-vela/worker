// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// Validate verifies the Worker is properly configured.
func (w *Worker) Validate() error {
	// log a message indicating the configuration verification
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Info
	logrus.Info("validating worker configuration")

	// check that hostname was properly populated
	if len(w.Config.API.Address.Hostname()) == 0 {
		switch strings.ToLower(w.Config.API.Address.Scheme) {
		case "http", "https":
			retErr := "worker server address invalid: %s"
			return fmt.Errorf(retErr, w.Config.API.Address.String())
		default:
			// hostname will be empty if a scheme is not provided
			retErr := "worker server address invalid, no scheme: %s"
			return fmt.Errorf(retErr, w.Config.API.Address.String())
		}
	}

	// verify a build limit was provided
	if w.Config.Build.Limit <= 0 {
		return fmt.Errorf("no worker build limit provided")
	}

	// verify a build timeout was provided
	if w.Config.Build.Timeout <= 0 {
		return fmt.Errorf("no worker build timeout provided")
	}

	// verify a worker address was provided
	if *w.Config.API.Address == (url.URL{}) {
		return fmt.Errorf("no worker address provided")
	}

	// verify a server address was provided
	if len(w.Config.Server.Address) == 0 {
		return fmt.Errorf("no worker server address provided")
	}

	// verify an executor driver was provided
	if len(w.Config.Executor.Driver) == 0 {
		return fmt.Errorf("no worker executor driver provided")
	}

	// verify the runtime configuration
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime#Setup.Validate
	err := w.Config.Runtime.Validate()
	if err != nil {
		return err
	}

	return nil
}
