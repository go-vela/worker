// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"reflect"
	"testing"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
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

func TestKubernetes_ClientOpt_WithPodsTemplate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure          bool
		podsTemplateName string
		podsTemplatePath string
		wantName         string
		wantTemplate     *velav1alpha1.PipelinePodTemplate
	}{
		{
			failure:          false,
			podsTemplateName: "foo-bar-name",
			podsTemplatePath: "",
			wantName:         "foo-bar-name",
			wantTemplate:     nil,
		},
		{
			failure:          false,
			podsTemplateName: "",
			podsTemplatePath: "",
			wantName:         "",
			wantTemplate:     nil,
		},
		{
			failure:          false, // ignores missing files; can be added later
			podsTemplateName: "",
			podsTemplatePath: "testdata/does-not-exist.yaml",
			wantName:         "",
			wantTemplate:     nil,
		},
		{
			failure:          false,
			podsTemplateName: "",
			podsTemplatePath: "testdata/pipeline-pods-template-empty.yaml",
			wantName:         "",
			wantTemplate:     &velav1alpha1.PipelinePodTemplate{},
		},
		{
			failure:          false,
			podsTemplateName: "",
			podsTemplatePath: "testdata/pipeline-pods-template.yaml",
			wantName:         "",
			wantTemplate: &velav1alpha1.PipelinePodTemplate{
				Metadata: velav1alpha1.PipelinePodTemplateMeta{
					Annotations: map[string]string{"annotation/foo": "bar"},
					Labels:      map[string]string{"foo": "bar"},
				},
			},
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
