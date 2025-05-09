// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/worker/internal/image"
	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
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
							Image: _container.Image,
						},
					},
				},
			},
			container: _container,
		},
		{
			name:    "build stops before container execution with raw pauseImage",
			failure: false,
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
							// container not patched yet with correct image
							Image: pauseImage,
						},
					},
				},
			},
			container: _container,
		},
		{
			name:    "build stops before container execution with canonical pauseImage",
			failure: false,
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
							// container not patched yet with correct image
							Image: image.Parse(pauseImage),
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
		oldPod    *v1.Pod
		newPod    *v1.Pod
	}{
		{
			name:      "default order in ContainerStatuses",
			failure:   false,
			container: _container,
			oldPod:    _pod,
			newPod:    _pod,
		},
		{
			name:      "inverted order in ContainerStatuses",
			failure:   false,
			container: _container,
			oldPod:    _pod,
			newPod: &v1.Pod{
				ObjectMeta: _pod.ObjectMeta,
				TypeMeta:   _pod.TypeMeta,
				Spec:       _pod.Spec,
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					ContainerStatuses: []v1.ContainerStatus{
						// alternate order
						{
							Name: "step-github-octocat-1-echo",
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									Reason:   "Completed",
									ExitCode: 0,
								},
							},
							Image: "alpine:latest",
						},
						{
							Name: "step-github-octocat-1-clone",
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									Reason:   "Completed",
									ExitCode: 0,
								},
							},
							Image: "target/vela-git:v0.4.0",
						},
					},
				},
			},
		},
		{
			name:      "container goes from running to terminated",
			failure:   false,
			container: _container,
			oldPod: &v1.Pod{
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
			newPod: _pod,
		},
		{
			name:      "if client.Pod.Spec is empty podTracker fails",
			failure:   true,
			container: _container,
			oldPod:    _pod,
			newPod: &v1.Pod{
				ObjectMeta: _pod.ObjectMeta,
				TypeMeta:   _pod.TypeMeta,
				Status:     _pod.Status,
				// if client.Pod.Spec is empty, podTracker will fail
				//Spec:       _pod.Spec,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// set up the fake k8s clientset so that it returns the final/updated state
			_engine, err := NewMock(test.newPod)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			go func() {
				// simulate a re-sync/PodUpdate event
				_engine.SimulateResync(test.oldPod)
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

func Test_podTracker_inspectContainerStatuses(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		name       string
		trackedPod string
		ctnName    string
		terminated bool
		pod        *v1.Pod
	}{
		{
			name:       "container is terminated",
			trackedPod: "test/github-octocat-1",
			ctnName:    "step-github-octocat-1-clone",
			terminated: true,
			pod:        _pod,
		},
		{
			name:       "pod is pending",
			trackedPod: "test/github-octocat-1",
			ctnName:    "step-github-octocat-1-clone",
			terminated: false,
			pod: &v1.Pod{
				ObjectMeta: _pod.ObjectMeta,
				TypeMeta:   _pod.TypeMeta,
				Spec:       _pod.Spec,
				Status: v1.PodStatus{
					Phase: v1.PodPending,
				},
			},
		},
		{
			name:       "container is running",
			trackedPod: "test/github-octocat-1",
			ctnName:    "step-github-octocat-1-clone",
			terminated: false,
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
		},
		{
			name:       "pod has an untracked container",
			trackedPod: "test/github-octocat-1",
			ctnName:    "step-github-octocat-1-clone",
			terminated: true,
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
								Terminated: &v1.ContainerStateTerminated{
									Reason:   "Completed",
									ExitCode: 0,
								},
							},
						},
						{
							Name: "injected-by-admissions-controller",
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctnTracker := containerTracker{
				Name:       test.ctnName,
				Terminated: make(chan struct{}),
			}
			podTracker := podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				Containers: map[string]*containerTracker{},
				// other fields not used by inspectContainerStatuses
				// if they're needed, use newPodTracker
			}
			podTracker.Containers[test.ctnName] = &ctnTracker

			podTracker.inspectContainerStatuses(test.pod)

			func() {
				defer func() {
					//nolint:errcheck // repeat close() panics (otherwise it won't)
					recover()
				}()

				close(ctnTracker.Terminated)

				// this will only run if close() did not panic
				if test.terminated {
					t.Error("inspectContainerStatuses should have signaled termination")
				}
			}()
		})
	}
}
