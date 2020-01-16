// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics returns a Prometheus handler for serving go metrics
func Metrics() http.Handler {
	return promhttp.Handler()
}
