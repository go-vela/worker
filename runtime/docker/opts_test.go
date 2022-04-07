// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDocker_ClientOpt_WithPrivilegedImages(t *testing.T) {
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

func TestDocker_ClientOpt_WithHostVolumes(t *testing.T) {
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

func TestDocker_ClientOpt_WithLogger(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		logger  *logrus.Entry
	}{
		{
			failure: false,
			logger:  &logrus.Entry{},
		},
		{
			failure: false,
			logger:  nil,
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithLogger(test.logger),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithLogger should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithLogger returned err: %v", err)
		}

		if test.logger == nil && _service.Logger == nil {
			t.Errorf("_engine.Logger should not be nil even if nil is passed to WithLogger")
		}

		if test.logger != nil && !reflect.DeepEqual(_service.Logger, test.logger) {
			t.Errorf("WithLogger set %v, want %v", _service.Logger, test.logger)
		}
	}
}
