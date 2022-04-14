// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
)

func TestKubernetes_ClientOpt_WithConfigFile(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		file    string
		want    string
	}{
		{
			name:    "valid config file",
			failure: false,
			file:    "testdata/config",
			want:    "testdata/config",
		},
		{
			name:    "invalid config file",
			failure: true,
			file:    "testdata/config_empty",
			want:    "testdata/config_empty",
		},
		{
			name:    "missing config file",
			failure: true,
			file:    "testdata/config_missing",
			want:    "testdata/config_missing",
		},
		{
			name:    "InClusterConfig file missing",
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
		name      string
		failure   bool
		namespace string
		want      string
	}{
		{
			name:      "namespace",
			failure:   false,
			namespace: "foo",
			want:      "foo",
		},
		{
			name:      "empty namespace fails",
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
		name    string
		volumes []string
		want    []string
	}{
		{
			name:    "defined",
			volumes: []string{"/foo/bar.txt:/foo/bar.txt", "/tmp/baz.conf:/tmp/baz.conf"},
			want:    []string{"/foo/bar.txt:/foo/bar.txt", "/tmp/baz.conf:/tmp/baz.conf"},
		},
		{
			name:    "empty",
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
		name   string
		images []string
		want   []string
	}{
		{
			name:   "defined",
			images: []string{"alpine", "golang"},
			want:   []string{"alpine", "golang"},
		},
		{
			name:   "empty",
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
		name    string
		failure bool
		logger  *logrus.Entry
	}{
		{
			name:    "provided logger",
			failure: false,
			logger:  &logrus.Entry{},
		},
		{
			name:    "nil logger",
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

func TestKubernetes_ClientOpt_WithPodsTemplate(t *testing.T) {
	// setup tests
	tests := []struct {
		name             string
		failure          bool
		podsTemplateName string
		podsTemplatePath string
		wantName         string
		wantTemplate     *velav1alpha1.PipelinePodTemplate
	}{
		{
			name:             "name",
			failure:          false,
			podsTemplateName: "foo-bar-name",
			podsTemplatePath: "",
			wantName:         "foo-bar-name",
			wantTemplate:     nil,
		},
		{
			name:             "no name or path",
			failure:          false,
			podsTemplateName: "",
			podsTemplatePath: "",
			wantName:         "",
			wantTemplate:     nil,
		},
		{
			name:             "ignores missing files",
			failure:          false, // ignores missing files; can be added later
			podsTemplateName: "",
			podsTemplatePath: "testdata/does-not-exist.yaml",
			wantName:         "",
			wantTemplate:     nil,
		},
		{
			name:             "path-empty",
			failure:          false,
			podsTemplateName: "",
			podsTemplatePath: "testdata/pipeline-pods-template-empty.yaml",
			wantName:         "",
			wantTemplate:     &velav1alpha1.PipelinePodTemplate{},
		},
		{
			name:             "path",
			failure:          false,
			podsTemplateName: "",
			podsTemplatePath: "testdata/pipeline-pods-template.yaml",
			wantName:         "",
			wantTemplate: &velav1alpha1.PipelinePodTemplate{
				Metadata: velav1alpha1.PipelinePodTemplateMeta{
					Annotations: map[string]string{"annotation/foo": "bar"},
					Labels: map[string]string{
						"foo":      "bar",
						"pipeline": "this-is-ignored", // loaded in opts. Ignored in SetupBuild.
					},
				},
			},
		},
		{
			name:             "path-malformed",
			failure:          true,
			podsTemplateName: "",
			podsTemplatePath: "testdata/pipeline-pods-template-malformed.yaml",
			wantName:         "",
			wantTemplate:     nil,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := New(
			WithConfigFile("testdata/config"),
			WithNamespace("foo"),
			WithPodsTemplate(test.podsTemplateName, test.podsTemplatePath),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithPodsTemplate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithPodsTemplate returned err: %v", err)
		}

		if !reflect.DeepEqual(_engine.config.PipelinePodsTemplateName, test.wantName) {
			t.Errorf("WithPodsTemplate is %v, wantName %v", _engine.config.PipelinePodsTemplateName, test.wantName)
		}

		if test.wantTemplate != nil && !reflect.DeepEqual(_engine.PipelinePodTemplate, test.wantTemplate) {
			t.Errorf("WithPodsTemplate is %v, wantTemplate %v", _engine.PipelinePodTemplate, test.wantTemplate)
		}
	}
}
