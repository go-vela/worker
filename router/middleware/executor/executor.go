// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/worker/executor"
)

// Retrieve gets the repo in the given context.
func Retrieve(c *gin.Context) executor.Engine {
	return executor.FromGinContext(c)
}

// Establish sets the executor in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("executor")
		if len(param) == 0 {
			msg := "No executor parameter provided"

			c.AbortWithStatusJSON(http.StatusBadRequest, api.Error{Message: &msg})

			return
		}

		number, err := strconv.Atoi(param)
		if err != nil {
			msg := fmt.Sprintf("invalid executor parameter provided: %s", param)

			c.AbortWithStatusJSON(http.StatusBadRequest, api.Error{Message: &msg})

			return
		}

		// capture executors value from context
		value := c.Value("executors")
		if value == nil {
			msg := "no running executors found"

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

			return
		}

		// cast executors value to expected type
		executors, ok := value.(map[int]executor.Engine)
		if !ok {
			msg := "unable to get executors"

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

			return
		}

		logrus.Debugf("Reading executor %s", param)

		e, ok := executors[number]
		if !ok {
			msg := fmt.Sprintf("unable to get executor %s", param)

			c.AbortWithStatusJSON(http.StatusBadRequest, api.Error{Message: &msg})

			return
		}

		executor.WithGinContext(c, e)
		c.Next()
	}
}
