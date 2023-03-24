// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types"
	"github.com/go-vela/worker/router/middleware/token"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MustServer ensures the caller is the vela server.
func MustServer() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		vela, err := vela.NewClient(addr, "", nil)
		if err != nil {
			msg := fmt.Sprintf("error creating vela client: %s", err)

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}

		vela.Authentication.SetTokenAuth(tkn)

		_, err = vela.Authentication.ValidateToken()
		if err != nil {
			msg := fmt.Sprintf("error validating token: %s", err)

			err := c.Error(fmt.Errorf(msg))
			if err != nil {
				logrus.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

			return
		}
	}
}
