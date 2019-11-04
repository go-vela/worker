// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// helper function to validate all CLI configuration.
func validate(c *cli.Context) error {
	logrus.Debug("Validating CLI configuration")

	// validate core configuration
	err := validateCore(c)
	if err != nil {
		return err
	}

	// validate executor configuration
	err = validateExecutor(c)
	if err != nil {
		return err
	}

	// validate queue configuration
	err = validateQueue(c)
	if err != nil {
		return err
	}

	// validate runtime configuration
	err = validateRuntime(c)
	if err != nil {
		return err
	}

	return nil
}

// helper function to validate the core CLI configuration.
func validateCore(c *cli.Context) error {
	logrus.Trace("Validating core CLI configuration")

	if len(c.String("server-addr")) == 0 {
		return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag not specified")
	}

	if !strings.Contains(c.String("server-addr"), "://") {
		return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag must be <scheme>://<hostname> format")
	}

	if strings.HasSuffix(c.String("server-addr"), "/") {
		return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag must not have trailing slash")
	}

	if len(c.String("vela-secret")) == 0 {
		return fmt.Errorf("vela-secret (VELA_SECRET) flag not specified")
	}

	return nil
}

// helper function to validate the executor CLI configuration.
func validateExecutor(c *cli.Context) error {
	logrus.Trace("Validating executor CLI configuration")

	if len(c.String("executor-driver")) == 0 {
		return fmt.Errorf("executor-driver (VELA_EXECUTOR_DRIVER or EXECUTOR_DRIVER) flag not specified")
	}

	if c.Int("executor-threads") < 1 {
		return fmt.Errorf("executor-threads (VELA_EXECUTOR_THREADS or EXECUTOR_THREADS) flag improperly configured")
	}

	return nil
}

// helper function to validate the queue CLI configuration.
func validateQueue(c *cli.Context) error {
	logrus.Trace("Validating queue CLI configuration")

	if len(c.String("queue-driver")) == 0 {
		return fmt.Errorf("queue-driver (VELA_QUEUE_DRIVER or QUEUE_DRIVER) flag not specified")
	}

	if len(c.String("queue-config")) == 0 {
		return fmt.Errorf("queue-config (VELA_QUEUE_CONFIG or QUEUE_CONFIG) flag not specified")
	}

	return nil
}

// helper function to validate the runtime CLI configuration.
func validateRuntime(c *cli.Context) error {
	logrus.Trace("Validating runtime CLI configuration")

	if len(c.String("runtime-driver")) == 0 {
		return fmt.Errorf("runtime-driver (VELA_RUNTIME_DRIVER or RUNTIME_DRIVER) flag not specified")
	}

	return nil
}
