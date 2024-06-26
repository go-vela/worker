// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
)

// helper function to setup the queue from the CLI arguments.
func setupClient(s *Server, token string) (*vela.Client, error) {
	logrus.Debug("creating vela client from worker configuration")

	// create a new Vela client from the server configuration
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#NewClient
	vela, err := vela.NewClient(s.Address, "", nil)
	if err != nil {
		return nil, err
	}
	// set token for authentication with the server
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#AuthenticationService.SetTokenAuth
	vela.Authentication.SetTokenAuth(token)

	return vela, nil
}
