// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/sdk-go/vela"

	"github.com/sirupsen/logrus"
)

// helper function to setup the queue from the CLI arguments.
func setupClient(s *Server, token string) (*vela.Client, error) {
	logrus.Debug("creating vela client from worker configuration")

	// create a new Vela client from the server configuration
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#NewClient
	vela, err := vela.NewClient(s.Address, "", nil)
	if err != nil {
		return nil, err
	}
	// set token for authentication with the server
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#AuthenticationService.SetTokenAuth
	vela.Authentication.SetTokenAuth(token)

	return vela, nil
}
