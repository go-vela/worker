// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/mock/server"

	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func TestExecutor_CreateService_Success(t *testing.T) {
	// setup
	r, _ := docker.NewMock()

	// setup context
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())
	c, _ := vela.NewClient(s.URL, nil)

	e, _ := New(c, r)
	e.WithPipeline(&pipeline.Build{
		Version: "1",
		ID:      "__0",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "service_org_repo_0_postgres;",
				Environment: map[string]string{},
				Image:       "postgres:11-alpine",
				Name:        "postgres",
				Ports:       []string{"5432:5432"},
				Pull:        true,
			},
		},
	})

	// run test
	got := e.CreateService(context.Background(), e.pipeline.Services[0])

	if got != nil {
		t.Errorf("CreateService is %v, want nil", got)
	}
}

func TestExecutor_PlanService_Success(t *testing.T) {
	// setup
	r, _ := docker.NewMock()

	// setup context
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())
	c, _ := vela.NewClient(s.URL, nil)

	e, _ := New(c, r)
	e.WithBuild(&library.Build{
		Number:       vela.Int(1),
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
		Branch:       vela.String("master"),
		Ref:          vela.String("refs/heads/master"),
		BaseRef:      vela.String(""),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("linux"),
	})
	e.WithPipeline(&pipeline.Build{
		Version: "1",
		ID:      "__0",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "service_org_repo_0_postgres;",
				Environment: map[string]string{},
				Image:       "postgres:11-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
			},
		},
	})
	e.WithRepo(&library.Repo{
		Org:         vela.String("github"),
		Name:        vela.String("octocat"),
		FullName:    vela.String("github/octocat"),
		Link:        vela.String("https://github.com/github/octocat"),
		Clone:       vela.String("https://github.com/github/octocat.git"),
		Branch:      vela.String("master"),
		Timeout:     vela.Int64(60),
		Visibility:  vela.String("public"),
		Private:     vela.Bool(false),
		Trusted:     vela.Bool(false),
		Active:      vela.Bool(true),
		AllowPull:   vela.Bool(false),
		AllowPush:   vela.Bool(true),
		AllowDeploy: vela.Bool(false),
		AllowTag:    vela.Bool(false),
	})

	// run test
	got := e.PlanService(context.Background(), e.pipeline.Services[0])

	if got != nil {
		t.Errorf("PlanService is %v, want nil", got)
	}
}

func TestExecutor_ExecService_Success(t *testing.T) {
	// setup
	r, _ := docker.NewMock()

	// setup context
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())
	c, _ := vela.NewClient(s.URL, nil)

	e, _ := New(c, r)
	e.WithBuild(&library.Build{
		Number:       vela.Int(1),
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
		Branch:       vela.String("master"),
		Ref:          vela.String("refs/heads/master"),
		BaseRef:      vela.String(""),
		Host:         vela.String("example.company.com"),
		Runtime:      vela.String("docker"),
		Distribution: vela.String("linux"),
	})
	e.WithPipeline(&pipeline.Build{
		Version: "1",
		ID:      "__0",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "service_org_repo_0_postgres;",
				Environment: map[string]string{},
				Image:       "postgres:11-alpine",
				Name:        "postgres",
				Ports:       []string{"5432:5432"},
			},
		},
	})
	e.WithRepo(&library.Repo{
		Org:         vela.String("github"),
		Name:        vela.String("octocat"),
		FullName:    vela.String("github/octocat"),
		Link:        vela.String("https://github.com/github/octocat"),
		Clone:       vela.String("https://github.com/github/octocat.git"),
		Branch:      vela.String("master"),
		Timeout:     vela.Int64(60),
		Visibility:  vela.String("public"),
		Private:     vela.Bool(false),
		Trusted:     vela.Bool(false),
		Active:      vela.Bool(true),
		AllowPull:   vela.Bool(false),
		AllowPush:   vela.Bool(true),
		AllowDeploy: vela.Bool(false),
		AllowTag:    vela.Bool(false),
	})
	e.serviceLogs.Store(e.pipeline.Services[0].ID, new(library.Log))
	e.services.Store(e.pipeline.Services[0].ID, new(library.Service))

	// run test
	got := e.ExecService(context.Background(), e.pipeline.Services[0])

	if got != nil {
		t.Errorf("ExecService is %v, want nil", got)
	}
}

func TestExecutor_DestroyService_Success(t *testing.T) {
	// setup
	r, _ := docker.NewMock()

	// setup context
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())
	c, _ := vela.NewClient(s.URL, nil)

	e, _ := New(c, r)
	e.WithPipeline(&pipeline.Build{
		Version: "1",
		ID:      "__0",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "service_org_repo_0_postgres;",
				Environment: map[string]string{},
				Image:       "postgres:11-alpine",
				Name:        "postgres",
				Ports:       []string{"5432:5432"},
				ExitCode:    0,
			},
		},
	})

	// run test
	got := e.DestroyService(context.Background(), e.pipeline.Services[0])

	if got != nil {
		t.Errorf("DestroyService is %v, want nil", got)
	}
}
