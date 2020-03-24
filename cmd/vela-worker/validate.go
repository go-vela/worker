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
	logrus.Info("validating worker configuration")

	// verify a build limit was provided
	if w.Build.Limit <= 0 {
		return fmt.Errorf("no worker build limit provided")
	}

	// verify a build timeout was provided
	if w.Build.Timeout <= 0 {
		return fmt.Errorf("no worker build timeout provided")
	}

	// verify a server address was provided
	if len(w.Server.Address) == 0 {
		return fmt.Errorf("no worker server address provided")
	}

	// verify a server secret was provided
	if len(w.Server.Secret) == 0 {
		return fmt.Errorf("no worker server secret provided")
	}

	// verify an executor driver was provided
	if len(w.Executor.Driver) == 0 {
		return fmt.Errorf("no worker executor driver provided")
	}

	// verify the queue configuration
	err := w.Queue.Validate()
	if err != nil {
		return err
	}

	// verify the runtime configuration
	err = w.Runtime.Validate()
	if err != nil {
		return err
	}

	return nil
}
