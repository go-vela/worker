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

	// verify an executor driver was provided
	if len(w.Executor.Driver) == 0 {
		return fmt.Errorf("no worker executor driver provided")
	}

	// verify a queue driver was provided
	if len(w.Queue.Driver) == 0 {
		return fmt.Errorf("no worker queue driver provided")
	}

	// verify a runtime driver was provided
	if len(w.Runtime.Driver) == 0 {
		return fmt.Errorf("no worker runtime driver provided")
	}

	return nil
}
