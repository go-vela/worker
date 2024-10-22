// SPDX-License-Identifier: Apache-2.0

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/worker/router/middleware/token"
)

// MustServer ensures the caller is the vela server.
func MustServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// retrieve the callers token from the request headers
		tkn, err := token.Retrieve(c.Request)
		if err != nil {
			msg := fmt.Sprintf("error parsing token: %v", err)

			logrus.Error(msg)

			c.AbortWithStatusJSON(http.StatusBadRequest, api.Error{Message: &msg})

			return
		}

		// retrieve the configured server address from the context
		addr := c.MustGet("server-address").(string)

		// create a temporary client to validate the incoming request
		vela, err := vela.NewClient(addr, "vela-worker", nil)
		if err != nil {
			msg := fmt.Sprintf("error creating vela client: %s", err)

			logrus.Error(msg)

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

			return
		}

		// validate a token was provided
		if strings.EqualFold(tkn, "") {
			msg := "missing token"

			logrus.Error(msg)

			c.AbortWithStatusJSON(http.StatusBadRequest, api.Error{Message: &msg})

			return
		}

		// set the token auth provided in the callers request header
		vela.Authentication.SetTokenAuth(tkn)

		// validate the token with the configured vela server
		resp, err := vela.Authentication.ValidateToken()
		if err != nil {
			msg := fmt.Sprintf("error validating token: %s", err)

			logrus.Error(msg)

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Error{Message: &msg})

			return
		}

		// if ValidateToken returned anything other than 200 consider the token invalid
		if resp.StatusCode != http.StatusOK {
			msg := "unable to validate token"

			logrus.Error(msg)

			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Error{Message: &msg})

			return
		}
	}
}
