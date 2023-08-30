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
)

// swagger:operation POST /queue-registration system Queue Registration
//
// Fill queue registration channel in worker to continue operation
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing queue address and queue public key
//   required: true
//   schema:
//     "$ref": "#/definitions/Queue"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully passed queue address and queue public key to worker
//     schema:
//       type: string
//   '401':
//     description: No queue address and queue public key was passed
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to pass queue address and queue public key to worker
//     schema:
//       "$ref": "#/definitions/Error"

// QueueRegistration will pass the json body of queue address and queue public key to the queue registration
// channel of the worker. This will unblock operation if the queue configuration details are not setup
func QueueRegistration(c *gin.Context) {
	res := new(library.QueueRegistration)
	v, ok := c.Get("queue-registration")

	if !ok {
		c.JSON(http.StatusInternalServerError, "no queue registration channel in the context")
		return
	}

	// make sure we configured the channel properly
	rChan, ok := v.(chan library.QueueRegistration)
	if !ok {
		c.JSON(http.StatusInternalServerError, "queue signing key channel in the context is the wrong type")
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
		retErr := fmt.Errorf("unable to decode JSON for queue details %w", err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}
	// Decode public key to validate if key is base64 encoded
	publicKeyDecoded, err := base64.StdEncoding.DecodeString(res.GetPublicKey())
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad public key was provided")
		return
	}

	if len(publicKeyDecoded) == 0 {
		c.JSON(http.StatusBadRequest, "Provided public key is empty")
		return
	}
	// validate decoded public key length
	if len(publicKeyDecoded) != 32 {
		c.JSON(http.StatusBadRequest, "no valid signing public key provided")
		return
	}

	// verify a queue address was provided
	if len(res.GetQueueAddress()) == 0 {
		c.JSON(http.StatusBadRequest, "no queue address provided")
		return
	}

	// check if the queue address has a scheme
	if !strings.Contains(res.GetQueueAddress(), "://") {
		c.JSON(http.StatusBadRequest, "queue address must be fully qualified (<scheme>://<host>)")
		return
	}

	// check if the queue address has a trailing slash
	if strings.HasSuffix(res.GetQueueAddress(), "/") {
		c.JSON(http.StatusBadRequest, "queue address must not have trailing slash")
		return
	}

	if res.GetQueueAddress() != "" && res.GetPublicKey() != "" {
		rChan <- *res
		c.JSON(http.StatusOK, "successfully passed queue details to worker")
		return

	} else {
		c.JSON(http.StatusBadRequest, "both public key and queue address are required")
		return
	}
}
