// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// swagger:operation GET /api/v1/metrics router Metrics
//
// Retrieve metrics from the worker
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successful 'ping' of Vela API
//     schema:
//       type: string

// Metrics returns a Prometheus handler for serving go metrics
func Metrics() http.Handler {
	return promhttp.Handler()
}
