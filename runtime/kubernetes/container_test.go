// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestKubernetes_InspectContainer(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
			failure:   false,
			container: new(pipeline.Container),
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.InspectContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("InspectContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("InspectContainer returned err: %v", err)
		}
	}
}

func TestKubernetes_RemoveContainer(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.RemoveContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("RemoveContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RemoveContainer returned err: %v", err)
		}
	}
}

func TestKubernetes_RunContainer(t *testing.T) {
	// TODO: include VolumeMounts?
	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
		pipeline  *pipeline.Build
		pod       *v1.Pod
		volumes   []string
	}{
		{
			failure:   false,
			container: _container,
			pipeline:  _stages,
			pod:       _pod,
		},
		{
			failure:   false,
			container: _container,
			pipeline:  _steps,
			pod:       _pod,
		},
	}

	// run tests
	for _, test := range tests {
		_engine, err := NewMock(test.pod)
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		if len(test.volumes) > 0 {
			_engine.config.Volumes = test.volumes
		}

		err = _engine.RunContainer(context.Background(), test.container, test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("RunContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("RunContainer returned err: %v", err)
		}
	}
}

func TestKubernetes_SetupContainer(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
	}{
		{
			failure:   false,
			container: _container,
		},
		{
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo", "hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      2,
				Pull:        "always",
			},
		},
		{
			failure: false,
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo", "hello"},
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Image:       "target/vela-docker:latest",
				Name:        "echo",
				Number:      2,
				Pull:        "always",
			},
		},
	}

	// run tests
	for _, test := range tests {
		err = _engine.SetupContainer(context.Background(), test.container)

		// this does not test the resulting pod spec (ie no tests for ImagePullPolicy, VolumeMounts)

		if test.failure {
			if err == nil {
				t.Errorf("SetupContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("SetupContainer returned err: %v", err)
		}
	}
}

// TODO: implement this once they resolve the bug
//
// https://github.com/kubernetes/kubernetes/issues/84203
func TestKubernetes_TailContainer(t *testing.T) {
	// Unfortunately, we can't implement this test using
	// the native Kubernetes fake. This is because there
	// is a bug in that code where an "empty" request is
	// always returned when calling the GetLogs function.
	//
	// https://github.com/kubernetes/kubernetes/issues/84203
	// fixed in k8s.io/client-go v0.19.0; we already have v0.22.2
}

func TestKubernetes_WaitContainer(t *testing.T) {
	// setup tests
	tests := []struct {
		failure   bool
		container *pipeline.Container
		object    runtime.Object
	}{
		{
			failure:   false,
			container: _container,
			object:    _pod,
		},
		{
			failure:   false,
			container: _container,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "github-octocat-1",
					Namespace: "test",
					Labels: map[string]string{
						"pipeline": "github-octocat-1",
					},
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "step-github-octocat-1-echo",
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									Reason:   "Completed",
									ExitCode: 0,
								},
							},
						},
						{
							Name: "step-github-octocat-1-clone",
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									Reason:   "Completed",
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
		},
		{
			failure:   true,
			container: _container,
			object:    new(v1.PodTemplate),
		},
	}

	// run tests
	for _, test := range tests {
		// setup types
		_engine, _watch, err := newMockWithWatch(_pod, "pods")
		if err != nil {
			t.Errorf("unable to create runtime engine: %v", err)
		}

		go func() {
			// simulate adding a pod to the watcher
			_watch.Add(test.object)
		}()

		err = _engine.WaitContainer(context.Background(), test.container)

		if test.failure {
			if err == nil {
				t.Errorf("WaitContainer should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WaitContainer returned err: %v", err)
		}
	}
}
