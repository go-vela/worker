// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
)

// helper function to setup the vela API client for worker check-in and build executable retrieval.
func setupClient(s *Server, token string) (*vela.Client, error) {
	logrus.Debug("creating vela client from worker configuration")

	vela, err := vela.NewClient(s.Address, "", nil)
	if err != nil {
		return nil, err
	}

	vela.Authentication.SetTokenAuth(token)

	return vela, nil
}

// helper function to setup the vela API client for executor calls to update build resources.
func setupExecClient(s *Server, buildToken, scmToken string) (*vela.Client, error) {
	logrus.Debug("creating vela client from worker configuration")

	vela, err := vela.NewClient(s.Address, "", nil)
	if err != nil {
		return nil, err
	}

	vela.Authentication.SetBuildTokenAuth(buildToken, scmToken)

	return vela, nil
}
