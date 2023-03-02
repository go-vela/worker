// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"net/http"
)

// Retrieve gets the worker in the given context.
func Retrieve(c *gin.Context) *library.Worker {
	return FromContext(c)
}

// Establish sets the worker in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		//wHostName := os.Getenv("WORKER_ADDR")
		wParam := util.PathParameter(c, "worker")
		if len(wParam) == 0 {
			retErr := fmt.Errorf("no worker parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
		//if wParam != wHostName {
		//	logrus.Infof("provided hostname %s does not match intended hostname %s", wParam, wHostName)
		//	retErr := fmt.Errorf("provided hostname %s does not match intended hostname", wParam)
		//	util.HandleError(c, http.StatusBadRequest, retErr)
		//
		//	return
		//}
		w := new(library.Worker)
		w.SetHostname(wParam)

		ToContext(c, w)
		c.Next()
	}
}
