// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
		volume string
		want   *Volume
	}{
		{
			volume: "/foo",
			want: &Volume{
				Source:      "/foo",
				Destination: "/foo",
				AccessMode:  "ro",
			},
		},
		{
			volume: "/foo:/bar",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			volume: "/foo:/bar:ro",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			volume: "/foo:/bar:rw",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "rw",
			},
		},
		{
			volume: "/foo:/bar:/foo:bar",
			want:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := Parse(test.volume)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Parse is %v, want %v", got, test.want)
		}
	}
}

func TestImage_ParseWithError(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		volume  string
		want    *Volume
	}{
		{
			failure: false,
			volume:  "/foo",
			want: &Volume{
				Source:      "/foo",
				Destination: "/foo",
				AccessMode:  "ro",
			},
		},
		{
			failure: false,
			volume:  "/foo:/bar",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			failure: false,
			volume:  "/foo:/bar:ro",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "ro",
			},
		},
		{
			failure: false,
			volume:  "/foo:/bar:rw",
			want: &Volume{
				Source:      "/foo",
				Destination: "/bar",
				AccessMode:  "rw",
			},
		},
		{
			failure: true,
			volume:  "/foo:/bar:/foo:bar",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := ParseWithError(test.volume)

		if test.failure {
			if err == nil {
				t.Errorf("ParseWithError should have returned err")
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParseWithError is %s want %s", got, test.want)
			}

			continue
		}

		if err != nil {
			t.Errorf("ParseWithError returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ParseWithError is %v, want %v", got, test.want)
		}
	}
}
