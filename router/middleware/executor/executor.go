// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"strconv"

	"github.com/go-vela/types"
	exec "github.com/go-vela/worker/executor"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the repo in the given context
func Retrieve(c *gin.Context) exec.Engine {
	return FromContext(c)
}

// Establish sets the executor in the given context
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("executor")
		if len(param) == 0 {
			msg := "No executor parameter provided"
			c.AbortWithStatusJSON(http.StatusBadRequest, types.Error{Message: &msg})
			return
		}

		number, err := strconv.Atoi(param)
		if err != nil {
			msg := fmt.Sprintf("invalid executor parameter provided: %s", param)
			c.AbortWithStatusJSON(http.StatusBadRequest, types.Error{Message: &msg})
			return
		}

		executors := exec.FromContext(c)
		logrus.Debugf("Reading executor %s", param)
		e, ok := executors[number]
		if !ok {
			msg := fmt.Sprintf("unable to get executor %s", param)
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})
			return
		}

		ToContext(c, e)
		c.Next()
	}
}
