// SPDX-License-Identifier: Apache-2.0

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
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)

	// run test
	got, err := Retrieve(request)
	if err == nil {
		t.Errorf("Retrieve should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}
