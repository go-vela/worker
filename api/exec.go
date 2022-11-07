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

// swagger:operation POST /api/v1/exec system Exec
//
// Perform a manual execution on the worker
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

// Exec represents the API handler to shutdown a
// executors currently running on an worker.
//
// This function performs a soft shut down of a worker.
// Any build running during this time will safely complete, then
// the worker will safely shut itself down.
func Exec(c *gin.Context) {
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

	pkgBytes, err := io.ReadAll(body)
	if err != nil {
		msg := "unable to bind item"

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	// TODO: vader: this should be a build package with secrets
	// (for now) it is the item with faked secrets
	pkg := new(types.BuildPackage)
	err = json.Unmarshal(pkgBytes, pkg)
	if err != nil {
		msg := "unable to bind item"

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	logrus.Info("Sending package over channel.")
	w.PackageChannel <- pkg

	c.JSON(http.StatusOK, "Executing build package.")
}
