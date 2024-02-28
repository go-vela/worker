// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// swagger:operation GET /metrics system Metrics
//
// Retrieve metrics from the worker
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successful retrieval of worker metrics
//     schema:
//       type: string

// Metrics returns a Prometheus handler for serving go metrics.
func Metrics() http.Handler {
	return promhttp.Handler()
}
