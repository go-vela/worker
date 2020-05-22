// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Validate verifies the Worker is properly configured.
func (w *Worker) Validate() error {
	// log a message indicating the configuration verification
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
	logrus.Info("validating worker configuration")

	// verify a build limit was provided
	if w.Config.Build.Limit <= 0 {
		return fmt.Errorf("no worker build limit provided")
	}

	// verify a build timeout was provided
	if w.Config.Build.Timeout <= 0 {
		return fmt.Errorf("no worker build timeout provided")
	}

	// verify a hostname was provided
	if len(w.Config.Hostname) == 0 {
		return fmt.Errorf("no worker hostname provided")
	}

	// verify a server address was provided
	if len(w.Config.Server.Address) == 0 {
		return fmt.Errorf("no worker server address provided")
	}

	// verify a server secret was provided
	if len(w.Config.Server.Secret) == 0 {
		return fmt.Errorf("no worker server secret provided")
	}

	// verify an executor driver was provided
	if len(w.Config.Executor.Driver) == 0 {
		return fmt.Errorf("no worker executor driver provided")
	}

	// verify the queue configuration
	//
	// https://godoc.org/github.com/go-vela/pkg-queue/queue#Setup.Validate
	err := w.Config.Queue.Validate()
	if err != nil {
		return err
	}

	// verify the runtime configuration
	//
	// https://godoc.org/github.com/go-vela/pkg-runtime/runtime#Setup.Validate
	err = w.Config.Runtime.Validate()
	if err != nil {
		return err
	}

	return nil
}
