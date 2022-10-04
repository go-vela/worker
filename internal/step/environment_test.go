// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/raw"
	"github.com/google/go-cmp/cmp"
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
		name      string
		failure   bool
		build     *library.Build
		container *pipeline.Container
		repo      *library.Repo
		step      *library.Step
		want      *pipeline.Container
	}{
		{
			name:      "success",
			failure:   false,
			build:     b,
			container: c,
			repo:      r,
			step:      s,
			want: &pipeline.Container{
				ID:        "step_github_octocat_1_init",
				Directory: "/home/github/octocat",
				Environment: map[string]string{
					"BUILD_AUTHOR":             "OctoKitty",
					"BUILD_AUTHOR_EMAIL":       "OctoKitty@github.com",
					"BUILD_BASE_REF":           "",
					"BUILD_BRANCH":             "master",
					"BUILD_CHANNEL":            "vela",
					"BUILD_CLONE":              "https://github.com/github/octocat.git",
					"BUILD_COMMIT":             "48afb5bdc41ad69bf22588491333f7cf71135163",
					"BUILD_CREATED":            "1563474076",
					"BUILD_ENQUEUED":           "1563474077",
					"BUILD_EVENT":              "push",
					"BUILD_HOST":               "example.company.com",
					"BUILD_LINK":               "https://example.company.com/github/octocat/1",
					"BUILD_MESSAGE":            "First commit...",
					"BUILD_NUMBER":             "1",
					"BUILD_PARENT":             "1",
					"BUILD_REF":                "refs/heads/master",
					"BUILD_SENDER":             "OctoKitty",
					"BUILD_SOURCE":             "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
					"BUILD_STARTED":            "1563474078",
					"BUILD_STATUS":             "running",
					"BUILD_TITLE":              "push received from https://github.com/github/octocat",
					"BUILD_WORKSPACE":          "/vela/src",
					"FOO":                      "bar",
					"REPOSITORY_ACTIVE":        "true",
					"REPOSITORY_ALLOW_COMMENT": "false",
					"REPOSITORY_ALLOW_DEPLOY":  "false",
					"REPOSITORY_ALLOW_PULL":    "false",
					"REPOSITORY_ALLOW_PUSH":    "true",
					"REPOSITORY_ALLOW_TAG":     "false",
					"REPOSITORY_BRANCH":        "master",
					"REPOSITORY_CLONE":         "https://github.com/github/octocat.git",
					"REPOSITORY_FULL_NAME":     "github/octocat",
					"REPOSITORY_LINK":          "https://github.com/github/octocat",
					"REPOSITORY_NAME":          "octocat",
					"REPOSITORY_ORG":           "github",
					"REPOSITORY_PRIVATE":       "false",
					"REPOSITORY_TIMEOUT":       "30",
					"REPOSITORY_TRUSTED":       "false",
					"REPOSITORY_VISIBILITY":    "public",
					"VELA_BUILD_AUTHOR":        "OctoKitty",
					"VELA_BUILD_AUTHOR_EMAIL":  "OctoKitty@github.com",
					"VELA_BUILD_BASE_REF":      "",
					"VELA_BUILD_BRANCH":        "master",
					"VELA_BUILD_CHANNEL":       "vela",
					"VELA_BUILD_CLONE":         "https://github.com/github/octocat.git",
					"VELA_BUILD_COMMIT":        "48afb5bdc41ad69bf22588491333f7cf71135163",
					"VELA_BUILD_CREATED":       "1563474076",
					"VELA_BUILD_DISTRIBUTION":  "linux",
					"VELA_BUILD_ENQUEUED":      "1563474077",
					"VELA_BUILD_EVENT":         "push",
					"VELA_BUILD_EVENT_ACTION":  "",
					"VELA_BUILD_HOST":          "example.company.com",
					"VELA_BUILD_LINK":          "https://example.company.com/github/octocat/1",
					"VELA_BUILD_MESSAGE":       "First commit...",
					"VELA_BUILD_NUMBER":        "1",
					"VELA_BUILD_PARENT":        "1",
					"VELA_BUILD_REF":           "refs/heads/master",
					"VELA_BUILD_RUNTIME":       "docker",
					"VELA_BUILD_SENDER":        "OctoKitty",
					"VELA_BUILD_SOURCE":        "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
					"VELA_BUILD_STARTED":       "1563474078",
					"VELA_BUILD_STATUS":        "running",
					"VELA_BUILD_TITLE":         "push received from https://github.com/github/octocat",
					"VELA_BUILD_WORKSPACE":     "/vela/src",
					"VELA_DISTRIBUTION":        "linux",
					"VELA_HOST":                "example.company.com",
					"VELA_REPO_ACTIVE":         "true",
					"VELA_REPO_ALLOW_COMMENT":  "false",
					"VELA_REPO_ALLOW_DEPLOY":   "false",
					"VELA_REPO_ALLOW_PULL":     "false",
					"VELA_REPO_ALLOW_PUSH":     "true",
					"VELA_REPO_ALLOW_TAG":      "false",
					"VELA_REPO_BRANCH":         "master",
					"VELA_REPO_BUILD_LIMIT":    "0",
					"VELA_REPO_CLONE":          "https://github.com/github/octocat.git",
					"VELA_REPO_FULL_NAME":      "github/octocat",
					"VELA_REPO_LINK":           "https://github.com/github/octocat",
					"VELA_REPO_NAME":           "octocat",
					"VELA_REPO_ORG":            "github",
					"VELA_REPO_PIPELINE_TYPE":  "",
					"VELA_REPO_PRIVATE":        "false",
					"VELA_REPO_TIMEOUT":        "30",
					"VELA_REPO_TRUSTED":        "false",
					"VELA_REPO_VISIBILITY":     "public",
					"VELA_RUNTIME":             "docker",
					"VELA_STEP_CREATED":        "1563474076",
					"VELA_STEP_DISTRIBUTION":   "linux",
					"VELA_STEP_EXIT_CODE":      "0",
					"VELA_STEP_HOST":           "example.company.com",
					"VELA_STEP_IMAGE":          "target/vela-git:v0.3.0",
					"VELA_STEP_NAME":           "clone",
					"VELA_STEP_NUMBER":         "1",
					"VELA_STEP_RUNTIME":        "docker",
					"VELA_STEP_STAGE":          "",
					"VELA_STEP_STARTED":        "1563474078",
					"VELA_STEP_STATUS":         "running",
					"VELA_VERSION":             "v0.0.0",
				},
				Image:  "#init",
				Name:   "init",
				Number: 1,
				Pull:   "always",
			},
		},
		{
			name:      "nil failure",
			failure:   true,
			build:     nil,
			container: nil,
			repo:      nil,
			step:      nil,
			want:      nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Environment(test.container, test.build, test.repo, test.step, "v0.0.0")

			if test.failure {
				if err == nil {
					t.Errorf("%s Environment should have returned err", test.name)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("%s Environment returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, test.container); diff != "" {
				t.Errorf("%s Environment mismatch (-want +got):\n%v", test.name, diff)
			}
		})
	}
}
