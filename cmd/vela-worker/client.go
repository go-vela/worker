// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/sdk-go/vela"

	"github.com/sirupsen/logrus"
)

// helper function to setup the queue from the CLI arguments.
func setupClient(s *Server) (*vela.Client, error) {
	logrus.Debug("creating vela client from worker configuration")

	vela, err := vela.NewClient(s.Address, nil)
	if err != nil {
		return nil, err
	}
	// set token for auth
	vela.Authentication.SetTokenAuth(s.Secret)

	return vela, nil
}
