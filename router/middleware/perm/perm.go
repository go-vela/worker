// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/types"
	"github.com/go-vela/worker/router/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MustServer ensures the user is the vela server
func MustServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := user.Retrieve(c)

		if strings.EqualFold(u.GetName(), "vela-server") {
			return
		}

		msg := fmt.Sprintf("User %s is not a platform admin", u.GetName())
		err := c.Error(fmt.Errorf(msg))
		if err != nil {
			logrus.Error(err)
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &msg})
	}
}
