// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/router/middleware/token"
)

// swagger:operation POST /queue-key system
//
// Fill public signing key channel in worker to continue operation
//
// ---
// produces:
// - application/json
// parameters:
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully passed public key to worker
//     schema:
//       type: string
//   '401':
//     description: No public key was passed
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to pass public key to worker
//     schema:
//       "$ref": "#/definitions/Error"

// QueueKey will pass the pubkey given in the request header to the queue-signing-key
// channel of the worker. This will unblock operation if the queue has not been setup.
func QueueKey(c *gin.Context) { // extract the public key channel that was packed into gin context
	v, ok := c.Get("queue-signing-key")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no queue signing key channel in the context")
		return
	}

	// make sure we configured the channel properly
	rChan, ok := v.(chan string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "queue signing key channel in the context is the wrong type")
		return
	}

	// if key is present in the channel, deny registration
	// this will likely never happen as the channel is offloaded immediately
	if len(rChan) > 0 {
		c.JSON(http.StatusOK, "queue key already provided")
		return
	}

	// retrieve public key from header
	t, err := token.Retrieve(c.Request)
	if err != nil {
		// an error occurs when no token was passed
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	publicKeyDecoded, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad public key was provided")
		return
	}

	if len(publicKeyDecoded) == 0 {
		c.JSON(http.StatusBadRequest, "Provided public key is empty")
		return
	}
	// write pubkey token to queue-signing-key channel
	rChan <- t

	c.JSON(http.StatusOK, "successfully passed public key to worker")
}
