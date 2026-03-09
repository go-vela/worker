// SPDX-License-Identifier: Apache-2.0

package step

import (
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

func TestStep_Skip(t *testing.T) {
	// setup types
	_repo := &api.Repo{
		ID:          new(int64(1)),
		Org:         new("github"),
		Name:        new("octocat"),
		FullName:    new("github/octocat"),
		Link:        new("https://github.com/github/octocat"),
		Clone:       new("https://github.com/github/octocat.git"),
		Branch:      new("main"),
		Timeout:     new(int32(60)),
		Visibility:  new("public"),
		Private:     new(false),
		Trusted:     new(false),
		Active:      new(true),
		AllowEvents: api.NewEventsFromMask(1),
	}

	_build := &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("push"),
		EventAction:  new(""),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_comment := &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("comment"),
		EventAction:  new("created"),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_deploy := &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("deployment"),
		EventAction:  new(""),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_deployFromTag := &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("deployment"),
		EventAction:  new(""),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/tags/v1.0.0"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_schedule := &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("schedule"),
		EventAction:  new(""),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_tag := &api.Build{
		ID:           new(int64(1)),
		Number:       new(int64(1)),
		Repo:         _repo,
		Parent:       new(int64(1)),
		Event:        new("tag"),
		EventAction:  new(""),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}

	_container := &pipeline.Container{
		ID:          "step_github_octocat_1_init",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "#init",
		Name:        constants.InitName,
		Number:      1,
		Pull:        "always",
	}

	tests := []struct {
		name      string
		build     *api.Build
		container *pipeline.Container
		want      bool
	}{
		{
			name:      "build",
			build:     _build,
			container: _container,
			want:      false,
		},
		{
			name:      "comment",
			build:     _comment,
			container: _container,
			want:      false,
		},
		{
			name:      "deploy",
			build:     _deploy,
			container: _container,
			want:      false,
		},
		{
			name:      "deployFromTag",
			build:     _deployFromTag,
			container: _container,
			want:      false,
		},
		{
			name:      "schedule",
			build:     _schedule,
			container: _container,
			want:      false,
		},
		{
			name:      "tag",
			build:     _tag,
			container: _container,
			want:      false,
		},
		{
			name:      "skip nil",
			build:     nil,
			container: nil,
			want:      true,
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Skip(test.container, test.build, test.build.GetStatus())
			if err != nil {
				t.Errorf("Skip returned error: %s", err)
			}

			if got != test.want {
				t.Errorf("Skip is %v, want %v", got, test.want)
			}
		})
	}
}
