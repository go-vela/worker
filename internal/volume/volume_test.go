// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package volume

import (
	"reflect"
	"testing"
)

func TestVolume_Parse(t *testing.T) {
	// setup tests
	tests := []struct {
		name   string
		volume string
		want   *Volume
	}{
		{
			name:   "same src and dest",
			volume: "/foo",
			want: &Volume{
				Source:      "/foo",
				Destination: "/foo",
				AccessMode:  "ro",
			},
		},
		{
			name:   "different src and dest",
			volume: "/foo:/bar",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			name:   "read-only different src and dest",
			volume: "/foo:/bar:ro",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			name:   "read-write different src and dest",
			volume: "/foo:/bar:rw",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "rw",
			},
		},
		{
			name:   "invalid",
			volume: "/foo:/bar:/foo:bar",
			want:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Parse(test.volume)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Parse is %v, want %v", got, test.want)
			}
		})
	}
}

func TestImage_ParseWithError(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		volume  string
		want    *Volume
	}{
		{
			name:    "same src and dest",
			failure: false,
			volume:  "/foo",
			want: &Volume{
				Source:      "/foo",
				Destination: "/foo",
				AccessMode:  "ro",
			},
		},
		{
			name:    "different src and dest",
			failure: false,
			volume:  "/foo:/bar",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			name:    "read-only different src and dest",
			failure: false,
			volume:  "/foo:/bar:ro",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			name:    "read-write different src and dest",
			failure: false,
			volume:  "/foo:/bar:rw",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "rw",
			},
		},
		{
			name:    "invalid",
			failure: true,
			volume:  "/foo:/bar:/foo:bar",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseWithError(test.volume)

			if test.failure {
				if err == nil {
					t.Errorf("ParseWithError should have returned err")
				}

				if !reflect.DeepEqual(got, test.want) {
					t.Errorf("ParseWithError is %s want %s", got, test.want)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ParseWithError returned err: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParseWithError is %v, want %v", got, test.want)
			}
		})
	}
}
