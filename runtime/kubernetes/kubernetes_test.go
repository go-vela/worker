// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"testing"

	"github.com/go-vela/types/pipeline"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubernetes_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure   bool
		namespace string
		path      string
	}{
		{
			failure:   false,
			namespace: "test",
			path:      "testdata/config",
		},
		{
			failure:   true,
			namespace: "test",
			path:      "testdata/config_empty",
		},
		// An empty path implies that we are running in kubernetes,
		// so we should use InClusterConfig. Tests, however, do not
		// run in kubernetes, so we would need a way to mock the
		// return value of rest.InClusterConfig(), but how?
		//{
		//	failure:   false,
		//	namespace: "test",
		//	path:      "",
		//},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithConfigFile(test.path),
			WithNamespace(test.namespace),
		)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
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
				},
				{
					Name: "step-github-octocat-1-echo",
					State: v1.ContainerState{
						Terminated: &v1.ContainerStateTerminated{
							Reason:   "Completed",
							ExitCode: 0,
						},
					},
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
)
