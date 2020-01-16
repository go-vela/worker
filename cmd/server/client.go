// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/sdk-go/vela"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// helper function to setup the queue from the CLI arguments.
func setupClient(c *cli.Context) (*vela.Client, error) {
	log.Debug("Creating vela client from CLI configuration")

	vela, err := vela.NewClient(c.String("server-addr"), nil)
	if err != nil {
		return nil, err
	}
	// set token for auth
	vela.Authentication.SetTokenAuth(c.String("vela-secret"))

	return vela, nil
}
