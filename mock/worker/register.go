// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"encoding/base64"
	"fmt"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// postRegister returns mock JSON for a http POST.
//
// Do not pass an auth token to fail the request.
func postRegister(c *gin.Context) {
	res := new(library.WorkerRegistration)
	// Binding request body into QueueRegistration struct
	err := c.Bind(res)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for worker-registration details %w", err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}
	// validate encoded public key from the JSON body
	err = validatePubKey(res.GetPublicKey())
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	// validate encoded public key from the JSON body
	err = validateQueueAddress(res.GetQueueAddress())
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, "successfully passed token to worker")
}

// validatePubKey is a helper function to validate
// the provided pubkey
func validatePubKey(s string) error {
	// Decode public key to validate if key is base64 encoded
	publicKeyDecoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("bad public key was provided")
	}

	if len(publicKeyDecoded) == 0 {
		return fmt.Errorf("provided public key is empty")
	}
	// validate decoded public key length
	if len(publicKeyDecoded) != 32 {
		return fmt.Errorf("no valid signing public key provided")
	}

	return nil
}

// validateQueueAddress is a helper function to validate
// the provided queue address
func validateQueueAddress(s string) error {
	// verify a queue address was provided
	if len(s) == 0 {
		return fmt.Errorf("no queue address provided")
	}

	// check if the queue address has a scheme
	if !strings.Contains(s, "://") {
		return fmt.Errorf("queue address must be fully qualified (<scheme>://<host>)")
	}

	// check if the queue address has a trailing slash
	if strings.HasSuffix(s, "/") {
		return fmt.Errorf("queue address must not have trailing slash")
	}

	return nil
}
