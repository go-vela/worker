// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"

	v1 "k8s.io/api/core/v1"
)

func TestKubernetes_InspectBuild(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			failure:  false,
			pipeline: _stages,
		},
		{
			failure:  false,
			pipeline: _steps,
		},
	}

	// run tests
	for _, test := range tests {
		_, err = _engine.InspectBuild(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("InspectBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectBuild returned err: %v", err)
		}
	}
}

func TestKubernetes_SetupBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
		opts     []ClientOpt
	}{
		{
			failure:  false,
			pipeline: _stages,
			opts:     nil,
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     nil,
		},
		{
			failure:  false,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-empty.yaml")},
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-empty.yaml")},
		},
		{
			failure:  false,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template.yaml")},
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template.yaml")},
		},
		{
			failure:  false,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-security-context.yaml")},
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-security-context.yaml")},
		},
		{
			failure:  false,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-node-selection.yaml")},
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-node-selection.yaml")},
		},
		{
			failure:  false,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-dns.yaml")},
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-dns.yaml")},
		},
		{
			failure:  false,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("mock-pipeline-pods-template", "")},
		},
		{
			failure:  false,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("mock-pipeline-pods-template", "")},
		},
		{
			failure:  true,
			pipeline: _stages,
			opts:     []ClientOpt{WithPodsTemplate("missing-pipeline-pods-template", "")},
		},
		{
			failure:  true,
			pipeline: _steps,
			opts:     []ClientOpt{WithPodsTemplate("missing-pipeline-pods-template", "")},
		},
	}

	// run tests
	for _, test := range tests {
		// setup types
		_engine, err := NewMock(&v1.Pod{}, test.opts...)
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		err = _engine.SetupBuild(context.Background(), test.pipeline)

		// this does not test the resulting pod spec (ie no tests for ObjectMeta, RestartPolicy)

		if test.failure {
			if err == nil {
				t.Errorf("SetupBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("SetupBuild returned err: %v", err)
		}
	}
}

func TestKubernetes_AssembleBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		pipeline *pipeline.Build
		// k8sPod is the pod that the mock Kubernetes client will return
		k8sPod *v1.Pod
		// enginePod is the pod under construction in the Runtime Engine
		enginePod *v1.Pod
	}{
		{
			failure:   false,
			pipeline:  _stages,
			k8sPod:    &v1.Pod{},
			enginePod: _stagesPod,
		},
		{
			failure:   false,
			pipeline:  _steps,
			k8sPod:    &v1.Pod{},
			enginePod: _pod,
		},
		{
			failure:   true,
			pipeline:  _stages,
			k8sPod:    _stagesPod,
			enginePod: _stagesPod,
		},
		{
			failure:   true,
			pipeline:  _steps,
			k8sPod:    _pod,
			enginePod: _pod,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := NewMock(test.k8sPod)
		_engine.Pod = test.enginePod

		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		err = _engine.AssembleBuild(context.Background(), test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("AssembleBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("AssembleBuild returned err: %v", err)
		}
	}
}

func TestKubernetes_RemoveBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		failure    bool
		createdPod bool
		pipeline   *pipeline.Build
		pod        *v1.Pod
	}{
		{
			failure:    false,
			createdPod: true,
			pipeline:   _stages,
			pod:        _pod,
		},
		{
			failure:    false,
			createdPod: true,
			pipeline:   _steps,
			pod:        _pod,
		},
		{
			failure:    false,
			createdPod: false,
			pipeline:   _stages,
			pod:        &v1.Pod{},
		},
		{
			failure:    false,
			pipeline:   _steps,
			pod:        &v1.Pod{},
			createdPod: false,
		},
		{
			failure:    true,
			pipeline:   _stages,
			pod:        &v1.Pod{},
			createdPod: true,
		},
		{
			failure:    true,
			pipeline:   _steps,
			pod:        &v1.Pod{},
			createdPod: true,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := NewMock(test.pod)
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		_engine.createdPod = test.createdPod

		err = _engine.RemoveBuild(context.Background(), test.pipeline)
		if test.failure {
			if err == nil {
				t.Errorf("RemoveBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RemoveBuild returned err: %v", err)
		}
	}
}
