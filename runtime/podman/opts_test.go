// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"reflect"
	"testing"
)

func TestPodman_ClientOpt_WithPrivilegedImages(t *testing.T) {
	// setup tests
	tests := []struct {
		images []string
		want   []string
	}{
		{
			images: []string{"alpine", "golang"},
			want:   []string{"alpine", "golang"},
		},
		{
			images: []string{},
			want:   []string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithPrivilegedImages(test.images),
		)

		if err != nil {
			t.Errorf("WithPrivilegedImages returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Images, test.want) {
			t.Errorf("WithPrivilegedImages is %v, want %v", _service.config.Images, test.want)
		}
	}
}

func TestPodman_ClientOpt_WithHostVolumes(t *testing.T) {
	// setup tests
	tests := []struct {
		volumes []string
		want    []string
	}{
		{
			volumes: []string{"/foo/bar.txt:/foo/bar.txt", "/tmp/baz.conf:/tmp/baz.conf"},
			want:    []string{"/foo/bar.txt:/foo/bar.txt", "/tmp/baz.conf:/tmp/baz.conf"},
		},
		{
			volumes: []string{},
			want:    []string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithHostVolumes(test.volumes),
		)

		if err != nil {
			t.Errorf("WithHostVolumes returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Volumes, test.want) {
			t.Errorf("WithHostVolumes is %v, want %v", _service.config.Volumes, test.want)
		}
	}
}
