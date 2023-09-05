// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// swagger:operation POST /register system Register
//
// Fill registration token channel in worker to continue operation
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the details to register worker
//   required: true
//   schema:
//     "$ref": "#/definitions/WorkerRegistration"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully passed details to worker
//     schema:
//       type: string
//   '401':
//     description: No details was passed
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to pass details to worker
//     schema:
//       "$ref": "#/definitions/Error"

// Register will pass the token given in the request header to the register token
// channel of the worker. This will unblock operation if the worker has not been
// registered and the provided registration token is valid.
func Register(c *gin.Context) {
	res := new(library.WorkerRegistration)
	// extract the worker hostname that was packed into gin context
	v, ok := c.Get("worker-registration")

	if !ok {
		c.JSON(http.StatusInternalServerError, "no worker-registrationn channel in the context")
		return
	}
	// make sure we configured the channel properly
	rChan, ok := v.(chan library.WorkerRegistration)
	if !ok {
		c.JSON(http.StatusInternalServerError, "worker-registration channel in the context is the wrong type")
		return
	}
	// if key is present in the channel, deny registration
	// this will likely never happen as the channel is offloaded immediately
	if len(rChan) > 0 {
		c.JSON(http.StatusOK, "queue details already provided")
		return
	}
	// Binding request body into QueueRegistration struct
	err := c.Bind(res)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for worker-registration details %w", err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	w, ok := c.Get("worker-hostname")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no worker hostname in the context")
		return
	}
	token := res.GetRegistrationToken()
	// extract the subject from the token
	sub, err := getSubjectFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	// make sure we configured the hostname properly
	hostname, ok := w.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "worker hostname in the context is the wrong type")
		return
	}

	// if the subject doesn't match the worker hostname return an error
	if sub != hostname {
		c.JSON(http.StatusUnauthorized, "worker hostname is invalid")
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

	// write registration token to auth token channel
	rChan <- *res

	c.JSON(http.StatusOK, "successfully passed registration details to worker")
}

// getSubjectFromToken is a helper function to extract
// the subject from the token claims.
func getSubjectFromToken(token string) (string, error) {
	// create a new JWT parser
	j := jwt.NewParser()

	// parse the payload
	t, _, err := j.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("unable to parse token")
	}

	sub, err := t.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("unable to get subject from token")
	}

	// make sure there was a subject defined
	if len(sub) == 0 {
		return "", fmt.Errorf("no subject defined in token")
	}

	return sub, nil
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
