// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/router/middleware/token"
)

// swagger:operation POST /register system Register
//
// Register the worker with the Vela server
//
// ---
// produces:
// - application/json
// parameters:
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully registered worker
//     schema:
//       type: string

// Health check the status of the application.
func Register(c *gin.Context) {
	// extract the auth token channel that was packed into gin context
	v, ok := c.Get("auth-token")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no auth token channel in the context")
		return
	}

	// make sure we configured the channel properly
	authChannel, ok := v.(chan string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "auth token channel in the context is the wrong type")
		return
	}

	// if auth token is present in the channel, deny registration
	if len(authChannel) > 0 {
		c.JSON(http.StatusOK, "worker already registered")
		return
	}

	// retrieve auth token from header
	token, err := token.Retrieve(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// write registration token to auth token channel
	authChannel <- token

	// somehow we need to make sure the registration worked
	// maybe a second channel for registration results?
	c.JSON(http.StatusOK, "successfully registered the worker")
}
