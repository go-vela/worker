// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

const (
	BuildResp = `{
  "id": 1,
  "repo": {
	"id": 1,
    "owner": {
	  	"id": 1,
	  	"name": "octocat",
	  	"favorites": [],
		"active": true,
    	"admin": false
  	},
  	"org": "github",
	"counter": 10,
	"name": "octocat",
	"full_name": "github/octocat",
	"link": "https://github.com/github/octocat",
	"clone": "https://github.com/github/octocat",
	"branch": "main",
	"build_limit": 10,
	"timeout": 60,
	"visibility": "public",
	"private": false,
	"trusted": true,
	"pipeline_type": "yaml",
	"topics": [],
	"active": true,
	"allow_events": {
		"push": {
			"branch": true,
			"tag": true
		},
		"pull_request": {
			"opened": true,
			"synchronize": true,
			"reopened": true,
			"edited": false
		},
		"deployment": {
			"created": true
		},
		"comment": {
			"created": false,
			"edited": false
		}
  	},
  "approve_build": "fork-always",
  "previous_name": ""
 },
  "pipeline_id": 1,
  "number": 1,
  "parent": 1,
  "event": "push",
  "event_action": "",
  "status": "created",
  "error": "",
  "enqueued": 1563474077,
  "created": 1563474076,
  "started": 1563474077,
  "finished": 0,
  "deploy": "",
  "deploy_number": 1,
  "deploy_payload": {},
  "clone": "https://github.com/github/octocat.git",
  "source": "https://github.com/github/octocat/commit/48afb5bdc41ad69bf22588491333f7cf71135163",
  "title": "push received from https://github.com/github/octocat",
  "message": "First commit...",
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "sender": "OctoKitty",
  "author": "OctoKitty",
  "email": "octokitty@github.com",
  "link": "https://vela.example.company.com/github/octocat/1",
  "branch": "main",
  "ref": "refs/heads/main",
  "head_ref": "",
  "base_ref": "",
  "host": "example.company.com",
  "runtime": "docker",
  "distribution": "linux",
  "approved_at": 0,
  "approved_by": ""
}`
)

// getBuild has a param :build returns mock JSON for a http GET.
func getBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(BuildResp)

	var body api.Build
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

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, BuildResp)
}
