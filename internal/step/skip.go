// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Skip creates the ruledata from the build and repository
// information and returns true if the data does not match
// the ruleset for the given container.
func Skip(c *pipeline.Container, b *library.Build, r *library.Repo) bool {
	// check if the container provided is empty
	if c == nil {
		return true
	}

	// create ruledata from build and repository information
	//
	// https://pkg.go.dev/github.com/go-vela/types/pipeline#RuleData
	ruledata := &pipeline.RuleData{
		Branch: b.GetBranch(),
		Event:  b.GetEvent(),
		Repo:   r.GetFullName(),
		Status: b.GetStatus(),
	}

	// check if the build event is tag
	if strings.EqualFold(b.GetEvent(), constants.EventTag) {
		// add tag information to ruledata with refs/tags prefix removed
		ruledata.Tag = strings.TrimPrefix(b.GetRef(), "refs/tags/")
	}

	// check if the build event is deployment
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		// add deployment target information to ruledata
		ruledata.Target = b.GetDeploy()
	}

	// return the inverse of container execute
	//
	// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.Execute
	return !c.Execute(ruledata)
}
