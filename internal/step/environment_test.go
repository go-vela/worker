// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/raw"
)

func TestStep_Environment(t *testing.T) {
	// setup types
	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
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
	b.SetBranch("master")
	b.SetRef("refs/heads/master")
	b.SetBaseRef("")
	b.SetHeadRef("changes")
	b.SetHost("example.company.com")
	b.SetRuntime("docker")
	b.SetDistribution("linux")

	c := &pipeline.Container{
		ID:          "step_github_octocat_1_init",
		Directory:   "/home/github/octocat",
		Environment: map[string]string{"FOO": "bar"},
		Image:       "#init",
		Name:        "init",
		Number:      1,
		Pull:        "always",
	}

	r := new(library.Repo)
	r.SetID(1)
	r.SetOrg("github")
	r.SetName("octocat")
	r.SetFullName("github/octocat")
	r.SetLink("https://github.com/github/octocat")
	r.SetClone("https://github.com/github/octocat.git")
	r.SetBranch("master")
	r.SetTimeout(30)
	r.SetVisibility("public")
	r.SetPrivate(false)
	r.SetTrusted(false)
	r.SetActive(true)
	r.SetAllowPull(false)
	r.SetAllowPush(true)
	r.SetAllowDeploy(false)
	r.SetAllowTag(false)
	r.SetAllowComment(false)

	s := new(library.Step)
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
		failure   bool
		build     *library.Build
		container *pipeline.Container
		repo      *library.Repo
		step      *library.Step
	}{
		{
			failure:   false,
			build:     b,
			container: c,
			repo:      r,
			step:      s,
		},
		{
			failure:   true,
			build:     nil,
			container: nil,
			repo:      nil,
			step:      nil,
		},
	}

	// run tests
	for _, test := range tests {
		err := Environment(test.container, test.build, test.repo, test.step, "v0.0.0")

		if test.failure {
			if err == nil {
				t.Errorf("Environment should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Environment returned err: %v", err)
		}
	}
}
