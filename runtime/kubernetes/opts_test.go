// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func TestKubernetes_ClientOpt_WithConfigFile(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    string
	}{
		{
			failure: false,
			file:    "testdata/config",
			want:    "testdata/config",
		},
		{
			failure: true,
			file:    "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithConfigFile(test.file),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithConfigFile should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithConfigFile returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.config.File, test.want) {
			t.Errorf("WithConfigFile is %v, want %v", _engine.config.File, test.want)
		}
	}
}

func TestKubernetes_ClientOpt_WithNamespace(t *testing.T) {
	// setup tests
	tests := []struct {
		failure   bool
		namespace string
		want      string
	}{
		{
			failure:   false,
			namespace: "foo",
			want:      "foo",
		},
		{
			failure:   true,
			namespace: "",
			want:      "",
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithConfigFile("testdata/config"),
			WithNamespace(test.namespace),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithNamespace should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithNamespace returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.config.Namespace, test.want) {
			t.Errorf("WithNamespace is %v, want %v", _engine.config.Namespace, test.want)
		}
	}
}

func TestKubernetes_ClientOpt_WithHostVolumes(t *testing.T) {
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
		_engine, err := New(
			WithConfigFile("testdata/config"),
			WithHostVolumes(test.volumes),
		)

		if err != nil {
			t.Errorf("WithHostVolumes returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.config.Volumes, test.want) {
			t.Errorf("WithHostVolumes is %v, want %v", _engine.config.Volumes, test.want)
		}
	}
}

func TestKubernetes_ClientOpt_WithPrivilegedImages(t *testing.T) {
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
		_engine, err := New(
			WithConfigFile("testdata/config"),
			WithPrivilegedImages(test.images),
		)

		if err != nil {
			t.Errorf("WithPrivilegedImages returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.config.Images, test.want) {
			t.Errorf("WithPrivilegedImages is %v, want %v", _engine.config.Images, test.want)
		}
	}
}

func TestKubernetes_ClientOpt_WithLogger(t *testing.T) {
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
		_engine, err := New(
			WithConfigFile("testdata/config"),
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

		if test.logger == nil && _engine.Logger == nil {
			t.Errorf("_engine.Logger should not be nil even if nil is passed to WithLogger")
		}

		if test.logger != nil && !reflect.DeepEqual(_engine.Logger, test.logger) {
			t.Errorf("WithLogger set %v, want %v", _engine.Logger, test.logger)
		}
	}
}
