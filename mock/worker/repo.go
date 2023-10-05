// SPDX-License-Identifier: Apache-2.0

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
	// RepoResp represents a JSON return for a single repo.
	RepoResp = `{
  "id": 1,
  "user_id": 1,
  "org": "github",
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
  "active": true,
  "allow_pr": false,
  "allow_push": true,
  "allow_deploy": false,
  "allow_tag": false
}`
)

// getRepo has a param :repo returns mock JSON for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func getRepo(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(RepoResp)

	var body library.Repo
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
