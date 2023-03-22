// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-vela/server/util"
	"github.com/go-vela/worker/router/middleware/token"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the user in the given context.
func Retrieve(c *gin.Context) *library.User {
	return FromContext(c)
}

// Establish sets the user in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := new(library.User)

		t, err := token.Retrieve(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		if secret, ok := c.Value("secret").(string); ok {
			if strings.EqualFold(t, secret) {
				u.SetName("vela-server")
				u.SetActive(true)
				u.SetAdmin(true)

				ToContext(c, u)
				c.Next()

				return
			}
		}

		// prepare the request to the worker
		client := http.DefaultClient
		client.Timeout = 30 * time.Second

		// set the API endpoint path we send the request to
		url := fmt.Sprintf("%s/validate-token", c.MustGet("server"))

		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			retErr := fmt.Errorf("unable to form a request to %s: %w", u, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add the token to authenticate to the worker
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))

		// perform the request to the server
		resp, err := client.Do(req)
		if err != nil {
			logrus.Debug("token validation for server token failed, adding nil user to context")
			ToContext(c, u)
			c.Next()

			return
		}
		defer resp.Body.Close()

		u.SetName("vela-server")
		u.SetActive(true)
		u.SetAdmin(true)

		ToContext(c, u)
		c.Next()
	}
}
