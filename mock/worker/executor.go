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
	// ExecutorResp represents a JSON return for a single worker.
	ExecutorResp = `
		{
			"id": 1,
			"host": "worker_1",
			"runtime": "docker",
			"distribution": "linux",
			"build": ` + BuildResp + `,
			"pipeline": ` + PipelineResp + `,
			"repo": ` + RepoResp + `
		}`

	// ExecutorsResp represents a JSON return for one to many workers.
	ExecutorsResp = `[ ` + ExecutorResp + `,` + ExecutorResp + `]`
)

// getExecutors returns mock JSON for a http GET.
func getExecutors(c *gin.Context) {
	data := []byte(ExecutorsResp)

	var body []library.Executor
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getExecutor has a param :executor returns mock JSON for a http GET.
func getExecutor(c *gin.Context) {
	w := c.Param("executor")

	if strings.EqualFold(w, "0") {
		msg := fmt.Sprintf("Executor %s does not exist", w)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(ExecutorResp)

	var body library.Executor
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
