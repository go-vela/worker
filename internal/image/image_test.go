// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package image

import (
	"strings"
	"testing"
)

func TestImage_Parse(t *testing.T) {
	// setup tests
	tests := []struct {
		image string
		want  string
	}{
		{
			image: "golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "golang:latest",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "library/golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "library/golang:1.14",
			want:  "docker.io/library/golang:1.14",
		},
		{
			image: "docker.io/library/golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "docker.io/library/golang:latest",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "index.docker.io/library/golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "index.docker.io/library/golang:latest",
			want:  "docker.io/library/golang:latest",
		},
		{
			image: "gcr.io/library/golang",
			want:  "gcr.io/library/golang:latest",
		},
		{
			image: "gcr.io/library/golang:latest",
			want:  "gcr.io/library/golang:latest",
		},
		{
			image: "!@#$%^&*()",
			want:  "!@#$%^&*()",
		},
	}

	// run tests
	for _, test := range tests {
		got := Parse(test.image)

		if !strings.EqualFold(got, test.want) {
			t.Errorf("Parse is %s want %s", got, test.want)
		}
	}
}

func TestImage_ParseWithError(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		image   string
		want    string
	}{
		{
			failure: false,
			image:   "golang",
			want:    "docker.io/library/golang:latest",
		},
		{
			failure: false,
			image:   "golang:latest",
			want:    "docker.io/library/golang:latest",
		},
		{
			failure: false,
			image:   "golang:1.14",
			want:    "docker.io/library/golang:1.14",
		},
		{
			failure: true,
			image:   "!@#$%^&*()",
			want:    "!@#$%^&*()",
		},
		{
			failure: true,
			image:   "1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
			want:    "sha256:1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
		},
	}

	// run tests
	for _, test := range tests {
		got, err := ParseWithError(test.image)

		if test.failure {
			if err == nil {
				t.Errorf("ParseWithError should have returned err")
			}

			if !strings.EqualFold(got, test.want) {
				t.Errorf("ParseWithError is %s want %s", got, test.want)
			}

			continue
		}

		if err != nil {
			t.Errorf("ParseWithError returned err: %v", err)
		}

		if !strings.EqualFold(got, test.want) {
			t.Errorf("ParseWithError is %s want %s", got, test.want)
		}
	}
}

func TestImage_IsPrivilegedImage(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		image   string
		pattern string
		want    bool
	}{
		{
			name:    "test privileged image without tag",
			image:   "docker.company.com/foo/bar",
			pattern: "docker.company.com/foo/bar",
			want:    true,
		},
		{
			name:    "test privileged image with tag",
			image:   "docker.company.com/foo/bar:v0.1.0",
			pattern: "docker.company.com/foo/bar",
			want:    true,
		},
		{
			name:    "test privileged image with tag",
			image:   "docker.company.com/foo/bar",
			pattern: "docker.company.com/foo/bar:v0.1.0",
			want:    false,
		},
		{
			name:    "test privileged with bad image",
			image:   "!@#$%^&*()",
			pattern: "docker.company.com/foo/bar",
			want:    false,
		},
		{
			name:    "test privileged with bad pattern",
			image:   "docker.company.com/foo/bar",
			pattern: "!@#$%^&*()",
			want:    false,
		},
		{
			name:    "test privileged with on extended path image",
			image:   "docker.company.com/foo/bar",
			pattern: "docker.company.com/foo",
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := IsPrivilegedImage(test.image, test.pattern)
			if got != test.want {
				t.Errorf("IsPrivilegedImage is %v want %v", got, test.want)
			}
		})
	}
}
