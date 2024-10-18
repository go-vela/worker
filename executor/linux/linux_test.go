// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"math"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/mock/server"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/worker/runtime/docker"
)

func TestEqual(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	_linux, err := New(
		WithBuild(testBuild()),
		WithHostname("localhost"),
		WithPipeline(testSteps(constants.DriverDocker)),
		WithRuntime(_runtime),
		WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create linux executor: %v", err)
	}

	_alternate, err := New(
		WithBuild(testBuild()),
		WithHostname("a.different.host"),
		WithPipeline(testSteps(constants.DriverDocker)),
		WithRuntime(_runtime),
		WithVelaClient(_client),
	)
	if err != nil {
		t.Errorf("unable to create alternate local executor: %v", err)
	}

	tests := []struct {
		name string
		a    *client
		b    *client
		want bool
	}{
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "left nil",
			a:    nil,
			b:    _linux,
			want: false,
		},
		{
			name: "right nil",
			a:    _linux,
			b:    nil,
			want: false,
		},
		{
			name: "equal",
			a:    _linux,
			b:    _linux,
			want: true,
		},
		{
			name: "not equal",
			a:    _linux,
			b:    _alternate,
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Equal(test.a, test.b); got != test.want {
				t.Errorf("Equal() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLinux_New(t *testing.T) {
	// setup types
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	_client, err := vela.NewClient(s.URL, "", nil)
	if err != nil {
		t.Errorf("unable to create Vela API client: %v", err)
	}

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		build   *api.Build
	}{
		{
			name:    "with build",
			failure: false,
			build:   testBuild(),
		},
		{
			name:    "nil build",
			failure: true,
			build:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(
				WithBuild(test.build),
				WithHostname("localhost"),
				WithPipeline(testSteps(constants.DriverDocker)),
				WithRuntime(_runtime),
				WithVelaClient(_client),
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

// testBuild is a test helper function to create a Build
// type with all fields set to a fake value.
func testBuild() *api.Build {
	return &api.Build{
		ID:           vela.Int64(1),
		Number:       vela.Int(1),
		Repo:         testRepo(),
		Parent:       vela.Int(1),
		Event:        vela.String("push"),
		Status:       vela.String("success"),
		Error:        vela.String(""),
		Enqueued:     vela.Int64(1563474077),
		Created:      vela.Int64(1563474076),
		Started:      vela.Int64(1563474077),
		Finished:     vela.Int64(0),
		Deploy:       vela.String(""),
		Clone:        vela.String("https://github.com/github/octocat.git"),
		Source:       vela.String("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        vela.String("push received from https://github.com/github/octocat"),
		Message:      vela.String("First commit..."),
		Commit:       vela.String("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       vela.String("OctoKitty"),
		Author:       vela.String("OctoKitty"),
		Branch:       vela.String("main"),
		Ref:          vela.String("refs/heads/main"),
		BaseRef:      vela.String(""),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("linux"),
	}
}

// testRepo is a test helper function to create a Repo
// type with all fields set to a fake value.
func testRepo() *api.Repo {
	return &api.Repo{
		ID:         vela.Int64(1),
		Org:        vela.String("github"),
		Name:       vela.String("octocat"),
		FullName:   vela.String("github/octocat"),
		Link:       vela.String("https://github.com/github/octocat"),
		Clone:      vela.String("https://github.com/github/octocat.git"),
		Branch:     vela.String("main"),
		Timeout:    vela.Int64(60),
		Visibility: vela.String("public"),
		Private:    vela.Bool(false),
		Trusted:    vela.Bool(false),
		Active:     vela.Bool(true),
		Owner:      testUser(),
	}
}

// testUser is a test helper function to create a User
// type with all fields set to a fake value.
func testUser() *api.User {
	return &api.User{
		ID:        vela.Int64(1),
		Name:      vela.String("octocat"),
		Token:     vela.String("superSecretToken"),
		Favorites: vela.Strings([]string{"github/octocat"}),
		Active:    vela.Bool(true),
		Admin:     vela.Bool(false),
	}
}

// testSteps is a test helper function to create a steps
// pipeline with fake steps.
func testSteps(runtime string) *pipeline.Build {
	steps := &pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service_github_octocat_1_postgres",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		Steps: pipeline.ContainerSlice{
			{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo hello"},
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
			},
		},
		Secrets: pipeline.SecretSlice{
			{
				Name:   "foo",
				Key:    "github/octocat/foo",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			{
				Name:   "foo",
				Key:    "github/foo",
				Engine: "native",
				Type:   "org",
				Origin: &pipeline.Container{},
			},
			{
				Name:   "foo",
				Key:    "github/octokitties/foo",
				Engine: "native",
				Type:   "shared",
				Origin: &pipeline.Container{},
			},
		},
	}

	// apply any runtime-specific cleanups
	return steps.Sanitize(runtime)
}

// testPod is a test helper function to create a Pod
// type with all fields set to a fake value.
func testPod(useStages bool) *v1.Pod {
	// https://github.com/go-vela/worker/blob/main/runtime/kubernetes/kubernetes_test.go#L83
	pod := &v1.Pod{
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
		},
		Spec: v1.PodSpec{},
	}

	if useStages {
		pod.Spec.Containers = []v1.Container{
			{
				Name:            "github-octocat-1-clone-clone",
				Image:           "target/vela-git:v0.6.0",
				WorkingDir:      "/vela/src/github.com/octocat/helloworld",
				ImagePullPolicy: v1.PullAlways,
			},
			{
				Name:            "github-octocat-1-echo-echo",
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
		}
		pod.Status.ContainerStatuses = []v1.ContainerStatus{
			{
				Name: "github-octocat-1-clone-clone",
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						Reason:   "Completed",
						ExitCode: 0,
					},
				},
				Image: "target/vela-git:v0.6.0",
			},
			{
				Name: "github-octocat-1-echo-echo",
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						Reason:   "Completed",
						ExitCode: 0,
					},
				},
				Image: "alpine:latest",
			},
		}
	} else { // step
		pod.Spec.Containers = []v1.Container{
			{
				Name:            "step-github-octocat-1-clone",
				Image:           "target/vela-git:v0.6.0",
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
		}
		pod.Status.ContainerStatuses = []v1.ContainerStatus{
			{
				Name: "step-github-octocat-1-clone",
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						Reason:   "Completed",
						ExitCode: 0,
					},
				},
				Image: "target/vela-git:v0.6.0",
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
		}
	}

	return pod
}

// testPodFor is a test helper function to create a Pod
// using container names from a test pipeline.
func testPodFor(build *pipeline.Build) *v1.Pod {
	var (
		pod        *v1.Pod
		containers []v1.Container
	)

	workingDir := "/vela/src/github.com/octocat/helloworld"
	useStages := len(build.Stages) > 0
	pod = testPod(useStages)

	for _, service := range build.Services {
		containers = append(containers, v1.Container{
			Name:       service.ID,
			Image:      service.Image,
			WorkingDir: workingDir,
			// service.Pull should be one of: Always, Never, IfNotPresent
			ImagePullPolicy: v1.PullPolicy(service.Pull),
		})
	}

	if useStages {
		containers = append(containers, v1.Container{
			Name:            "step-github-octocat-1-clone-clone",
			Image:           "target/vela-git:v0.6.0",
			WorkingDir:      workingDir,
			ImagePullPolicy: v1.PullAlways,
		})
	} else { // steps
		containers = append(containers, v1.Container{
			Name:            "step-github-octocat-1-clone",
			Image:           "target/vela-git:v0.6.0",
			WorkingDir:      workingDir,
			ImagePullPolicy: v1.PullAlways,
		})
	}

	for _, stage := range build.Stages {
		for _, step := range stage.Steps {
			if step.Name == "init" {
				continue
			}

			containers = append(containers, v1.Container{
				Name:       step.ID,
				Image:      step.Image,
				WorkingDir: workingDir,
				// step.Pull should be one of: Always, Never, IfNotPresent
				ImagePullPolicy: v1.PullPolicy(step.Pull),
			})
		}
	}

	for _, step := range build.Steps {
		if step.Name == "init" {
			continue
		}

		containers = append(containers, v1.Container{
			Name:       step.ID,
			Image:      step.Image,
			WorkingDir: workingDir,
			// step.Pull should be one of: Always, Never, IfNotPresent
			ImagePullPolicy: v1.PullPolicy(step.Pull),
		})
	}

	for _, secret := range build.Secrets {
		if secret.Origin.Empty() {
			continue
		}

		containers = append(containers, v1.Container{
			Name:       secret.Origin.ID,
			Image:      secret.Origin.Image,
			WorkingDir: workingDir,
			// secret.Origin.Pull should be one of: Always, Never, IfNotPresent
			ImagePullPolicy: v1.PullPolicy(secret.Origin.Pull),
		})
	}

	pod.Spec.Containers = containers
	pod.Status.ContainerStatuses = testContainerStatuses(build, false, 0, 0)

	return pod
}

// countBuildSteps counts the steps in the build.
func countBuildSteps(build *pipeline.Build) int {
	steps := 0

	for _, stage := range build.Stages {
		for _, step := range stage.Steps {
			if step.Name == "init" {
				continue
			}

			steps++
		}
	}

	for _, step := range build.Steps {
		if step.Name == "init" {
			continue
		}

		steps++
	}

	return steps
}

// testContainerStatuses is a test helper function to create a ContainerStatuses list.
func testContainerStatuses(build *pipeline.Build, servicesRunning bool, stepsRunningCount, stepsCompletedPercent int) []v1.ContainerStatus {
	var containerStatuses []v1.ContainerStatus

	useStages := len(build.Stages) > 0
	stepsCompletedCount := 0

	if stepsCompletedPercent > 0 {
		stepsCompletedCount = int(math.Round(float64(stepsCompletedPercent) / 100 * float64(countBuildSteps(build))))
	}

	if servicesRunning {
		for _, service := range build.Services {
			containerStatuses = append(containerStatuses, v1.ContainerStatus{
				Name: service.ID,
				State: v1.ContainerState{
					Running: &v1.ContainerStateRunning{},
				},
				Image: service.Image,
			})
		}
	}

	if useStages {
		containerStatuses = append(containerStatuses, v1.ContainerStatus{
			Name: "step-github-octocat-1-clone-clone",
			State: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{
					Reason:   "Completed",
					ExitCode: 0,
				},
			},
			Image: "target/vela-git:v0.6.0",
		})
	} else { // steps
		containerStatuses = append(containerStatuses, v1.ContainerStatus{
			Name: "step-github-octocat-1-clone",
			State: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{
					Reason:   "Completed",
					ExitCode: 0,
				},
			},
			Image: "target/vela-git:v0.6.0",
		})
	}

	steps := 0

	for _, stage := range build.Stages {
		for _, step := range stage.Steps {
			if step.Name == "init" {
				continue
			}

			steps++
			if steps > stepsCompletedCount+stepsRunningCount {
				break
			}

			if stepsRunningCount > 0 && steps > stepsCompletedCount {
				containerStatuses = append(containerStatuses, v1.ContainerStatus{
					Name: step.ID,
					State: v1.ContainerState{
						Running: &v1.ContainerStateRunning{},
					},
					Image: step.Image,
				})
			} else if steps <= stepsCompletedCount {
				containerStatuses = append(containerStatuses, v1.ContainerStatus{
					Name: step.ID,
					State: v1.ContainerState{
						Terminated: &v1.ContainerStateTerminated{
							Reason:   "Completed",
							ExitCode: 0,
						},
					},
					Image: step.Image,
				})
			}
		}
	}

	for _, step := range build.Steps {
		if step.Name == "init" {
			continue
		}

		steps++
		if steps > stepsCompletedCount+stepsRunningCount {
			break
		}

		if stepsRunningCount > 0 && steps > stepsCompletedCount {
			containerStatuses = append(containerStatuses, v1.ContainerStatus{
				Name: step.ID,
				State: v1.ContainerState{
					Running: &v1.ContainerStateRunning{},
				},
				Image: step.Image,
			})
		} else if steps <= stepsCompletedCount {
			containerStatuses = append(containerStatuses, v1.ContainerStatus{
				Name: step.ID,
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						Reason:   "Completed",
						ExitCode: 0,
					},
				},
				Image: step.Image,
			})
		}
	}

	for _, secret := range build.Secrets {
		if secret.Origin.Empty() {
			continue
		}

		containerStatuses = append(containerStatuses, v1.ContainerStatus{
			Name:  secret.Origin.ID,
			Image: secret.Origin.Image,
			State: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{
					Reason:   "Completed",
					ExitCode: 0,
				},
			},
		})
	}

	return containerStatuses
}
