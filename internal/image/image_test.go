// SPDX-License-Identifier: Apache-2.0

package image

import (
	"strings"
	"testing"
)

func TestImage_Parse(t *testing.T) {
	// setup tests
	tests := []struct {
		name  string
		image string
		want  string
	}{
		{
			name:  "image only",
			image: "golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "image and tag",
			image: "golang:latest",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "repo and image",
			image: "library/golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "repo image and tag",
			image: "library/golang:1.14",
			want:  "docker.io/library/golang:1.14",
		},
		{
			name:  "hub repo and image",
			image: "docker.io/library/golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "hub repo image and tag",
			image: "docker.io/library/golang:latest",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "alt hub with repo and image",
			image: "index.docker.io/library/golang",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "alt hub with repo image and tag",
			image: "index.docker.io/library/golang:latest",
			want:  "docker.io/library/golang:latest",
		},
		{
			name:  "gcr hub with repo and image",
			image: "gcr.io/library/golang",
			want:  "gcr.io/library/golang:latest",
		},
		{
			name:  "gcr hub with repo image and tag",
			image: "gcr.io/library/golang:latest",
			want:  "gcr.io/library/golang:latest",
		},
		{
			name:  "garbage in garbage out",
			image: "!@#$%^&*()",
			want:  "!@#$%^&*()",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Parse(test.image)

			if !strings.EqualFold(got, test.want) {
				t.Errorf("Parse is %s want %s", got, test.want)
			}
		})
	}
}

func TestImage_ParseWithError(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		image   string
		want    string
	}{
		{
			name:    "image only",
			failure: false,
			image:   "golang",
			want:    "docker.io/library/golang:latest",
		},
		{
			name:    "image and tag",
			failure: false,
			image:   "golang:latest",
			want:    "docker.io/library/golang:latest",
		},
		{
			name:    "image and tag",
			failure: false,
			image:   "golang:1.14",
			want:    "docker.io/library/golang:1.14",
		},
		{
			name:    "fails with bad image",
			failure: true,
			image:   "!@#$%^&*()",
			want:    "!@#$%^&*()",
		},
		{
			name:    "fails with image sha",
			failure: true,
			image:   "1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
			want:    "sha256:1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseWithError(test.image)

			if test.failure {
				if err == nil {
					t.Errorf("ParseWithError should have returned err")
				}

				if !strings.EqualFold(got, test.want) {
					t.Errorf("ParseWithError is %s want %s", got, test.want)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ParseWithError returned err: %v", err)
			}

			if !strings.EqualFold(got, test.want) {
				t.Errorf("ParseWithError is %s want %s", got, test.want)
			}
		})
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
