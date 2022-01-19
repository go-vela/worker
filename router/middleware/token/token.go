// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"fmt"
	"net/http"
	"strings"
)

// nolint: godot // ignore comment ending in period
//
// Retrieve gets the token from the provided request http.Request
// to be parsed and validated. This is called on every request
// to enable capturing the user making the request and validating
// they have the proper access. The following methods of providing
// authentication to Vela are supported:
//
// * Bearer token in `Authorization` header
func Retrieve(r *http.Request) (string, error) {
	// get the token from the `Authorization` header
	token := r.Header.Get("Authorization")
	if len(token) > 0 {
		if strings.Contains(token, "Bearer") {
			return strings.Split(token, "Bearer ")[1], nil
		}
	}

	return "", fmt.Errorf("no token provided in Authorization header")
}
