// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-vela/types/pipeline"
	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestKubernetes_InspectContainer(t *testing.T) {
	// setup tests
	tests := []struct {
		name      string
		failure   bool
		pod       *v1.Pod
		container *pipeline.Container
	}{
		{
			name:      "build container",
			failure:   false,
			pod:       _pod,
			container: _container,
		},
		{
			name:      "empty build container",
			failure:   false,
			pod:       _pod,
			container: new(pipeline.Container),
		},
		{
			name:    "container not terminated",
			failure: true,
			pod: &v1.Pod{
				ObjectMeta: _pod.ObjectMeta,
				TypeMeta:   _pod.TypeMeta,
				Spec:       _pod.Spec,
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "step-github-octocat-1-clone",
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{},
							},
						},
					},
				},
			},
			container: _container,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			_engine, err := NewMock(test.pod)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			err = _engine.InspectContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("InspectContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectContainer returned err: %v", err)
			}
		})
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
		name      string
		failure   bool
		container *pipeline.Container
	}{
		{
			name:      "build container",
			failure:   false,
			container: _container,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = _engine.RemoveContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("RemoveContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RemoveContainer returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_RunContainer(t *testing.T) {
	// TODO: include VolumeMounts?
	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		pipeline  *pipeline.Build
		pod       *v1.Pod
		volumes   []string
	}{
		{
			name:      "stages",
			failure:   false,
			container: _container,
			pipeline:  _stages,
			pod:       _pod,
		},
		{
			name:      "steps",
			failure:   false,
			container: _container,
			pipeline:  _steps,
			pod:       _pod,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RunContainer returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_SetupContainer(t *testing.T) {
	// setup tests
	tests := []struct {
		name             string
		failure          bool
		container        *pipeline.Container
		opts             []ClientOpt
		wantPrivileged   bool
		wantFromTemplate interface{}
	}{
		{
			name:             "step-clone",
			failure:          false,
			container:        _container, // clone step
			opts:             nil,
			wantPrivileged:   false,
			wantFromTemplate: nil,
		},
		{
			name:    "step-echo",
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
			opts:             nil,
			wantPrivileged:   false,
			wantFromTemplate: nil,
		},
		{
			name:    "privileged",
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
			opts:             []ClientOpt{WithPrivilegedImages([]string{"target/vela-docker"})},
			wantPrivileged:   true,
			wantFromTemplate: nil,
		},
		{
			name:           "PipelinePodsTemplate",
			failure:        false,
			container:      _container,
			opts:           []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-security-context.yaml")},
			wantPrivileged: false,
			wantFromTemplate: velav1alpha1.PipelineContainerSecurityContext{
				Capabilities: &v1.Capabilities{
					Drop: []v1.Capability{"ALL"},
					Add:  []v1.Capability{"NET_ADMIN", "SYS_TIME"},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			_engine, err := NewMock(_pod.DeepCopy(), test.opts...)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}
			// actually run the test
			err = _engine.SetupContainer(context.Background(), test.container)

			// this does not (yet) test everything in the resulting pod spec (ie no tests for ImagePullPolicy, VolumeMounts)

			if test.failure {
				if err == nil {
					t.Errorf("SetupContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("SetupContainer returned err: %v", err)
			}

			// SetupContainer added the last pod so get it for inspection
			i := len(_engine.Pod.Spec.Containers) - 1
			ctn := _engine.Pod.Spec.Containers[i]

			// make sure the lookup map is working as expected
			if j := _engine.containersLookup[ctn.Name]; i != j {
				t.Errorf("expected containersLookup[ctn.Name] to be %d, got %d", i, j)
			}

			// Make sure Container has Privileged configured correctly
			if test.wantPrivileged {
				if ctn.SecurityContext == nil {
					t.Errorf("Pod.Containers[%v].SecurityContext is nil", i)
				} else if *ctn.SecurityContext.Privileged != test.wantPrivileged {
					t.Errorf("Pod.Containers[%v].SecurityContext.Privileged is %v, want %v", i, *ctn.SecurityContext.Privileged, test.wantPrivileged)
				}
			} else {
				if ctn.SecurityContext != nil && ctn.SecurityContext.Privileged != nil && *ctn.SecurityContext.Privileged != test.wantPrivileged {
					t.Errorf("Pod.Containers[%v].SecurityContext.Privileged is %v, want %v", i, *ctn.SecurityContext.Privileged, test.wantPrivileged)
				}
			}

			switch test.wantFromTemplate.(type) {
			case velav1alpha1.PipelineContainerSecurityContext:
				want := test.wantFromTemplate.(velav1alpha1.PipelineContainerSecurityContext)

				// PipelinePodsTemplate defined SecurityContext.Capabilities
				if want.Capabilities != nil {
					if ctn.SecurityContext == nil {
						t.Errorf("Pod.Containers[%v].SecurityContext is nil", i)
					} else if !reflect.DeepEqual(ctn.SecurityContext.Capabilities, want.Capabilities) {
						t.Errorf("Pod.Containers[%v].SecurityContext.Capabilities is %v, want %v", i, ctn.SecurityContext.Capabilities, want.Capabilities)
					}
				}
			}
		})
	}
}

func TestKubernetes_TailContainer(t *testing.T) {
	// Unfortunately, we can't test failures using the native Kubernetes fake.
	// k8s.client-go v0.19.0 added a mock GetLogs() response so that
	// it no longer panics with an "empty" request, but now it always returns
	// a successful response with Body: "fake logs".
	//
	// https://github.com/kubernetes/kubernetes/issues/84203
	// https://github.com/kubernetes/kubernetes/pulls/91485
	//
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
	}{
		{
			name:      "got logs",
			failure:   false,
			container: _container,
		},
		// We cannot test failures, because the mock GetLogs() always
		// returns a successful response with logs body: "fake logs"
		//{
		//	name:      "empty build container",
		//	failure:   true,
		//	container: new(pipeline.Container),
		//},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.TailContainer(context.Background(), test.container)

			if test.failure {
				if err == nil {
					t.Errorf("TailContainer should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("TailContainer returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_WaitContainer(t *testing.T) {
	// setup tests
	tests := []struct {
		name      string
		failure   bool
		container *pipeline.Container
		object    runtime.Object
	}{
		{
			name:      "default order in ContainerStatuses",
			failure:   false,
			container: _container,
			object:    _pod,
		},
		{
			name:      "inverted order in ContainerStatuses",
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
			name:      "watch returns invalid type",
			failure:   true,
			container: _container,
			object:    new(v1.PodTemplate),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

				return // continue to next test
			}

			if err != nil {
				t.Errorf("WaitContainer returned err: %v", err)
			}
		})
	}
}
