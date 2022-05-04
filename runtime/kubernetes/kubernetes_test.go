// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"fmt"
	"testing"

	"github.com/go-vela/types/pipeline"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubernetes_New(t *testing.T) {
	// setup tests
	tests := []struct {
		name      string
		failure   bool
		namespace string
		path      string
	}{
		{
			name:      "valid config file",
			failure:   false,
			namespace: "test",
			path:      "testdata/config",
		},
		{
			name:      "invalid config file",
			failure:   true,
			namespace: "test",
			path:      "testdata/config_empty",
		},
		// An empty path implies that we are running in kubernetes,
		// so we should use InClusterConfig. Tests, however, do not
		// run in kubernetes, so we would need a way to mock the
		// return value of rest.InClusterConfig(), but how?
		//{
		//	name:      "InClusterConfig file",
		//	failure:   false,
		//	namespace: "test",
		//	path:      "",
		//},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(
				WithConfigFile(test.path),
				WithNamespace(test.namespace),
			)

			if test.failure {
				if err == nil {
					t.Errorf("New should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("New returned err: %v", err)
			}
		})
	}
}

// setup global variables used for testing.
var (
	_container = &pipeline.Container{
		ID:          "step-github-octocat-1-clone",
		Directory:   "/vela/src/github.com/octocat/helloworld",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "target/vela-git:v0.4.0",
		Name:        "clone",
		Number:      2,
		Pull:        "always",
	}

	_stagesContainer = &pipeline.Container{
		ID:          "step-github-octocat-1-clone-clone",
		Directory:   _container.Directory,
		Environment: _container.Environment,
		Image:       _container.Image,
		Name:        _container.Name,
		Number:      _container.Number,
		Pull:        _container.Pull,
	}

	_pod = &v1.Pod{
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
					Name: "step-github-octocat-1-clone",
					State: v1.ContainerState{
						Terminated: &v1.ContainerStateTerminated{
							Reason:   "Completed",
							ExitCode: 0,
						},
					},
					Image: "target/vela-git:v0.4.0",
				},
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
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "step-github-octocat-1-clone",
					Image:           "target/vela-git:v0.4.0",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "step-github-octocat-1-echo",
					Image:           "alpine:latest",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "service-github-octocat-1-postgres",
					Image:           "postgres:12-alpine",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
			},
			HostAliases: []v1.HostAlias{
				{
					IP: "127.0.0.1",
					Hostnames: []string{
						"postgres.local",
						"echo.local",
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "github-octocat-1",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	_stepsPodBeforeRunStep = &v1.Pod{
		ObjectMeta: _pod.ObjectMeta,
		TypeMeta:   _pod.TypeMeta,
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "step-github-octocat-1-clone",
					Image: pauseImage, // step is not running yet
					State: v1.ContainerState{
						// pause is running, not the step image
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "step-github-octocat-1-echo",
					Image: pauseImage, // step is not running yet
					State: v1.ContainerState{
						// pause is running, not the step image
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "service-github-octocat-1-postgres",
					Image: "postgres:12-alpine",
					State: v1.ContainerState{
						// service is running
						Running: &v1.ContainerStateRunning{},
					},
				},
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "step-github-octocat-1-clone",
					Image:           pauseImage, // not running yet
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "step-github-octocat-1-echo",
					Image:           pauseImage, // not running yet
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "service-github-octocat-1-postgres",
					Image:           "postgres:12-alpine",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
			},
			HostAliases: _pod.Spec.HostAliases,
			Volumes:     _pod.Spec.Volumes,
		},
	}

	_stepsPodWithRunningStep = &v1.Pod{
		ObjectMeta: _pod.ObjectMeta,
		TypeMeta:   _pod.TypeMeta,
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "step-github-octocat-1-clone",
					Image: "target/vela-git:v0.4.0",
					State: v1.ContainerState{
						// step is running
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "step-github-octocat-1-echo",
					Image: pauseImage,
					State: v1.ContainerState{
						// pause is running, not the step image
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "service-github-octocat-1-postgres",
					Image: "postgres:12-alpine",
					State: v1.ContainerState{
						// service is running
						Running: &v1.ContainerStateRunning{},
					},
				},
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "step-github-octocat-1-clone",
					Image:           "target/vela-git:v0.4.0", // running
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "step-github-octocat-1-echo",
					Image:           pauseImage, // not running yet
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "service-github-octocat-1-postgres",
					Image:           "postgres:12-alpine",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
			},
			HostAliases: _pod.Spec.HostAliases,
			Volumes:     _pod.Spec.Volumes,
		},
	}

	_stages = &pipeline.Build{
		Version: "1",
		ID:      "github-octocat-1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service-github-octocat-1-postgres",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      4,
				Ports:       []string{"5432:5432"},
			},
		},
		Stages: pipeline.StageSlice{
			{
				Name: "init",
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step-github-octocat-1-init-init",
						Directory:   "/vela/src/github.com/octocat/helloworld",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
				},
			},
			{
				Name:  "clone",
				Needs: []string{"init"},
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step-github-octocat-1-clone-clone",
						Directory:   "/vela/src/github.com/octocat/helloworld",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-git:v0.4.0",
						Name:        "clone",
						Number:      2,
						Pull:        "always",
					},
				},
			},
			{
				Name:  "echo",
				Needs: []string{"clone"},
				Steps: pipeline.ContainerSlice{
					{
						ID:          "step-github-octocat-1-echo-echo",
						Commands:    []string{"echo hello"},
						Detach:      true,
						Directory:   "/vela/src/github.com/octocat/helloworld",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      3,
						Pull:        "always",
					},
				},
			},
		},
	}

	_steps = &pipeline.Build{
		Version: "1",
		ID:      "github-octocat-1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service-github-octocat-1-postgres",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      4,
				Ports:       []string{"5432:5432"},
			},
		},
		Steps: pipeline.ContainerSlice{
			{
				ID:          "step-github-octocat-1-init",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step-github-octocat-1-clone",
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step-github-octocat-1-echo",
				Commands:    []string{"echo hello"},
				Detach:      true,
				Directory:   "/vela/src/github.com/octocat/helloworld",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
			},
		},
	}

	_stagesPod = &v1.Pod{
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
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "step-github-octocat-1-clone-clone",
					Image:           "target/vela-git:v0.4.0",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "step-github-octocat-1-echo-echo",
					Image:           "alpine:latest",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "service-github-octocat-1-postgres",
					Image:           "postgres:12-alpine",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
			},
			HostAliases: []v1.HostAlias{
				{
					IP: "127.0.0.1",
					Hostnames: []string{
						"postgres.local",
						"echo.local",
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "github-octocat-1",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	_stagesPodBeforeRunStep = &v1.Pod{
		ObjectMeta: _stagesPod.ObjectMeta,
		TypeMeta:   _stagesPod.TypeMeta,
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "step-github-octocat-1-clone-clone",
					Image: pauseImage, // step is not running yet
					State: v1.ContainerState{
						// pause is running, not the step image
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "step-github-octocat-1-echo-echo",
					Image: pauseImage, // step is not running yet
					State: v1.ContainerState{
						// pause is running, not the step image
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "service-github-octocat-1-postgres",
					Image: "postgres:12-alpine",
					State: v1.ContainerState{
						// service is running
						Running: &v1.ContainerStateRunning{},
					},
				},
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "step-github-octocat-1-clone-clone",
					Image:           pauseImage, // not running yet
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "step-github-octocat-1-echo-echo",
					Image:           pauseImage, // not running yet
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "service-github-octocat-1-postgres",
					Image:           "postgres:12-alpine",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
			},
			HostAliases: _stagesPod.Spec.HostAliases,
			Volumes:     _stagesPod.Spec.Volumes,
		},
	}

	_stagesPodWithRunningStep = &v1.Pod{
		ObjectMeta: _stagesPod.ObjectMeta,
		TypeMeta:   _stagesPod.TypeMeta,
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "step-github-octocat-1-clone-clone",
					Image: "target/vela-git:v0.4.0",
					State: v1.ContainerState{
						// step is running
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "step-github-octocat-1-echo-echo",
					Image: pauseImage,
					State: v1.ContainerState{
						// pause is running, not the step image
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name:  "service-github-octocat-1-postgres",
					Image: "postgres:12-alpine",
					State: v1.ContainerState{
						// service is running
						Running: &v1.ContainerStateRunning{},
					},
				},
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "step-github-octocat-1-clone-clone",
					Image:           "target/vela-git:v0.4.0", // running
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "step-github-octocat-1-echo-echo",
					Image:           pauseImage, // not running yet
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
				{
					Name:            "service-github-octocat-1-postgres",
					Image:           "postgres:12-alpine",
					WorkingDir:      "/vela/src/github.com/octocat/helloworld",
					ImagePullPolicy: v1.PullAlways,
				},
			},
			HostAliases: _stagesPod.Spec.HostAliases,
			Volumes:     _stagesPod.Spec.Volumes,
		},
	}
)

func mockContainerEvent(pod *v1.Pod, ctnName, reason, message string) *v1.Event {
	return &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: pod.ObjectMeta.Namespace,
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      pod.TypeMeta.Kind,
			Name:      pod.ObjectMeta.Name,
			Namespace: pod.ObjectMeta.Namespace,
			FieldPath: fmt.Sprintf("spec.containers{%s}", ctnName),
		},
		Reason:  reason,
		Message: message,
	}
}
