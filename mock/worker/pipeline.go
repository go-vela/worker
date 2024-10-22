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
	// PipelineResp represents a JSON return for a single pipeline.
	PipelineResp = `{
  "id": 1,
  "repo_id": 1,
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "flavor": "",
  "platform": "",
  "ref": "refs/heads/main",
  "type": "yaml",
  "version": "1",
  "external_secrets": false,
  "internal_secrets": false,
  "services": false,
  "stages": false,
  "steps": true,
  "templates": false,
  "data": "LS0tCnZlcnNpb246ICIxIgoKc3RlcHM6CiAgLSBuYW1lOiBlY2hvCiAgICBpbWFnZTogYWxwaW5lOmxhdGVzdAogICAgY29tbWFuZHM6IFtlY2hvIGZvb10="
}`
)

// getPipeline has a param :pipeline returns mock YAML for a http GET.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func getPipeline(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(PipelineResp)

	var body api.Pipeline
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
