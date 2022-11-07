// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/worker/worker"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/challenge system Challenge
//
// Initiate a manual execution on the worker
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '501':
//     description: Endpoint is not yet implemented
//     schema:
//       type: string

// TODO:VADER: fillme

func Challenge(c *gin.Context) {
	// var err error

	// capture worker value from context
	value := c.Value("worker")
	if value == nil {
		msg := "no running worker found"

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	// cast executors value to expected type
	w, ok := value.(worker.Worker)
	if !ok {
		msg := "unable to get worker"

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	// read incoming body from the request
	body := c.Request.Body

	challengeBody, err := io.ReadAll(body)
	if err != nil {
		msg := "unable to bind item"

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}
	type Challenge struct {
		Challenge string `json:"challenge"`
		Token     string `json:"token"`
	}

	// TODO: vader: make this more secure
	challenge := new(Challenge)
	err = json.Unmarshal(challengeBody, challenge)
	if err != nil {
		msg := "unable to bind item"

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}
	challenge.Token = w.Config.Server.Secret

	logrus.Info("Responding to server challenge.")

	c.JSON(http.StatusOK, challenge)
}
