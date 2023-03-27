// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types"
	"github.com/go-vela/worker/router/middleware/token"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MustServer ensures the caller is the vela server.
func MustServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// retrieve the callers token from the request headers
		tkn, err := token.Retrieve(c.Request)
		if err != nil {
			msg := fmt.Sprintf("error parsing token: %s", err)

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())

			return
		}

		// retrieve the configured server address from the context
		addr, ok := c.MustGet("server-address").(string)
		if !ok {
			msg := "error retrieving server address"

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		// create a temporary client to validate the incoming request
		vela, err := vela.NewClient(addr, "vela-worker", nil)
		if err != nil {
			msg := fmt.Sprintf("error creating vela client: %s", err)

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		// validate a token was provided
		if strings.EqualFold(tkn, "") {
			msg := "missing token"

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, types.Error{Message: &msg})

			return
		}

		// set the token auth provided in the callers request header
		vela.Authentication.SetTokenAuth(tkn)

		// validate the token with the configured vela server
		resp, err := vela.Authentication.ValidateToken()
		if err != nil {
			msg := fmt.Sprintf("error validating token: %s", err)

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		// if ValidateToken returned anything other than 200 consider the token invalid
		if resp.StatusCode != http.StatusOK {
			msg := "unable to validate token"

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &msg})

			return
		}
	}
}
