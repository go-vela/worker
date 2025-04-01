// SPDX-License-Identifier: Apache-2.0

package outputs

import (
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/worker/internal/step"
)

func TestOutputs_Sanitize(t *testing.T) {
	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("github")
	r.SetName("octocat")
	r.SetFullName("github/octocat")
	r.SetLink("https://github.com/github/octocat")
	r.SetClone("https://github.com/github/octocat.git")
	r.SetBranch("main")
	r.SetTimeout(30)
	r.SetVisibility("public")
	r.SetPrivate(false)
	r.SetTrusted(false)
	r.SetActive(true)
	r.SetAllowEvents(api.NewEventsFromMask(1))

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetParent(1)
	b.SetEvent("push")
	b.SetStatus("running")
	b.SetError("")
	b.SetEnqueued(1563474077)
	b.SetCreated(1563474076)
	b.SetStarted(1563474078)
	b.SetFinished(1563474079)
	b.SetDeploy("")
	b.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	b.SetClone("https://github.com/github/octocat.git")
	b.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetTitle("push received from https://github.com/github/octocat")
	b.SetMessage("First commit...")
	b.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetSender("OctoKitty")
	b.SetAuthor("OctoKitty")
	b.SetEmail("OctoKitty@github.com")
	b.SetLink("https://example.company.com/github/octocat/1")
	b.SetBranch("main")
	b.SetRef("refs/heads/main")
	b.SetBaseRef("")
	b.SetHeadRef("changes")
	b.SetHost("example.company.com")
	b.SetRuntime("docker")
	b.SetDistribution("linux")

	s := new(api.Step)
	s.SetID(1)
	s.SetBuildID(1)
	s.SetRepoID(1)
	s.SetNumber(1)
	s.SetName("clone")
	s.SetImage("target/vela-git:v0.3.0")
	s.SetStatus("running")
	s.SetExitCode(0)
	s.SetCreated(1563474076)
	s.SetStarted(1563474078)
	s.SetFinished(1563474079)
	s.SetHost("example.company.com")
	s.SetRuntime("docker")
	s.SetDistribution("linux")

	// setup tests
	tests := []struct {
		name      string
		container *pipeline.Container
		outputs   map[string]string
		want      map[string]string
	}{
		{
			name: "defined env with new outputs",
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			outputs: map[string]string{
				"COLOR": "blue",
				"SIZE":  "large",
			},
			want: map[string]string{
				"COLOR": "blue",
				"SIZE":  "large",
			},
		},
		{
			name: "overwrite outputs",
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			outputs: map[string]string{
				"VELA_BUILD_EVENT": "blue",
				"SIZE":             "large",
			},
			want: map[string]string{
				"SIZE": "large",
			},
		},
		{
			name: "defined env with overwrite outputs",
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			outputs: map[string]string{
				"FOO": "baz",
			},
			want: map[string]string{},
		},
		{
			name:      "empty container",
			container: nil,
			outputs: map[string]string{
				"FOO": "baz",
			},
			want: nil,
		},
		{
			name: "empty outputs",
			container: &pipeline.Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			outputs: map[string]string{},
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = step.Environment(test.container, b, s, "v0.0.0", "ey123")

			got := Sanitize(test.container, test.outputs)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Sanitize is %v, want %v", got, test.want)
			}
		})
	}
}
