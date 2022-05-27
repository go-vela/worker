// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

const (
	// BuildResp represents a JSON return for a single build.
	BuildResp = `{
  "id": 1,
  "repo_id": 1,
  "number": 1,
  "parent": 1,
  "event": "push",
  "status": "created",
  "error": "",
  "enqueued": 1563474077,
  "created": 1563474076,
  "started": 1563474077,
  "finished": 0,
  "deploy": "",
  "clone": "https://github.com/github/octocat.git",
  "source": "https://github.com/github/octocat/commit/48afb5bdc41ad69bf22588491333f7cf71135163",
  "title": "push received from https://github.com/github/octocat",
  "message": "First commit...",
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "sender": "OctoKitty",
  "author": "OctoKitty",
  "email": "octokitty@github.com",
  "link": "https://vela.example.company.com/github/octocat/1",
  "branch": "master",
  "ref": "refs/heads/master",
  "base_ref": "",
  "host": "example.company.com",
  "runtime": "docker",
  "distribution": "linux"
}`
)

// getBuild has a param :build returns mock JSON for a http GET.
func getBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(BuildResp)

	var body library.Build
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// cancelBuild has a param :build returns mock JSON for a http DELETE.
//
// Pass "0" to :build to test receiving a http 404 response.
func cancelBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, BuildResp)
}
