// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"reflect"
	"testing"

	"github.com/go-vela/pkg-runtime/runtime/docker"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func TestLinux_WithBuild(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := int64(1)
	b := &library.Build{ID: &id}

	want, _ := New(vela, r)
	want.build = b

	// run test
	got, err := New(vela, r)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithBuild(b), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestLinux_WithPipeline(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := string(1)
	b := &pipeline.Build{ID: id}

	want, _ := New(vela, r)
	want.pipeline = b

	// run test
	got, err := New(vela, r)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithPipeline(b), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestLinux_WithRepo(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := int64(1)
	repo := &library.Repo{ID: &id}

	want, _ := New(vela, r)
	want.repo = repo

	// run test
	got, err := New(vela, r)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithRepo(repo), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestLinux_WithUser(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := int64(1)
	u := &library.User{ID: &id}

	want, _ := New(vela, r)
	want.user = u

	// run test
	got, err := New(vela, r)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithUser(u), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestLinux_GetBuild(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := int64(1)
	b := &library.Build{ID: &id}
	want := b

	executor, _ := New(vela, r)
	executor.WithBuild(b)

	// run test
	got, err := executor.GetBuild()
	if err != nil {
		t.Errorf("Unable to get build from compiler: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBuild is %v, want %v", got, want)
	}
}

func TestLinux_GetPipeline(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := "1"
	p := &pipeline.Build{ID: id}
	want := p

	executor, _ := New(vela, r)
	executor.WithPipeline(p)

	// run test
	got, err := executor.GetPipeline()
	if err != nil {
		t.Errorf("Unable to get pipeline from compiler: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetPipeline is %v, want %v", got, want)
	}
}

func TestLinux_GetRepo(t *testing.T) {
	// setup types
	vela, _ := vela.NewClient("http://localhost:8080", nil)
	r, _ := docker.NewMock()

	id := int64(1)
	repo := &library.Repo{ID: &id}
	want := repo

	executor, _ := New(vela, r)
	executor.WithRepo(repo)

	// run test
	got, err := executor.GetRepo()
	if err != nil {
		t.Errorf("Unable to get repo from compiler: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetPipeline is %v, want %v", got, want)
	}
}
