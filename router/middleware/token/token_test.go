// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestToken_Retrieve(t *testing.T) {
	// setup types
	want := "foobar"

	header := fmt.Sprintf("Bearer %s", want)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", header)

	// run test
	got, err := Retrieve(request)
	if err != nil {
		t.Errorf("Retrieve returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestToken_Retrieve_Error(t *testing.T) {
	// setup types
	request, _ :=  http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)

	// run test
	got, err := Retrieve(request)
	if err == nil {
		t.Errorf("Retrieve should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}
