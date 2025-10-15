// SPDX-License-Identifier: Apache-2.0

package step

import (
	"fmt"
	"strings"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/storage"
)

// Skip creates the ruledata from the build and repository
// information and returns true if the data does not match
// the ruleset for the given container.
func Skip(c *pipeline.Container, b *api.Build, status string, storage storage.Storage) (bool, error) {
	fmt.Printf("evaluating step %s for skip\n", c.Name)
	// check if the container provided is empty
	if c == nil {
		return true, nil
	}

	if !c.TestReport.Empty() {
		fmt.Printf("test report step detected\n")
		if !storage.StorageEnable() {
			fmt.Printf("Skipping step %s, storage is nil", c.Name)
			return true, nil
		}
	}

	event := b.GetEvent()
	action := b.GetEventAction()

	// if the build has an event action, concatenate event and event action for matching
	if !strings.EqualFold(action, "") {
		event = event + ":" + action
	}

	// create ruledata from build and repository information
	ruledata := &pipeline.RuleData{
		Branch: b.GetBranch(),
		Event:  event,
		Repo:   b.GetRepo().GetFullName(),
		Status: status,
		Env:    c.Environment,
	}

	// check if the build event is tag
	if strings.EqualFold(b.GetEvent(), constants.EventTag) {
		// add tag information to ruledata with refs/tags prefix removed
		ruledata.Tag = strings.TrimPrefix(b.GetRef(), "refs/tags/")
	}

	// check if the build event is deployment
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		// handle when deployment event is for a tag
		if strings.HasPrefix(b.GetRef(), "refs/tags/") {
			// add tag information to ruledata with refs/tags prefix removed
			ruledata.Tag = strings.TrimPrefix(b.GetRef(), "refs/tags/")
		}
		// add deployment target information to ruledata
		ruledata.Target = b.GetDeploy()
	}

	// check if the build event is schedule
	if strings.EqualFold(b.GetEvent(), constants.EventSchedule) {
		// add schedule target information to ruledata
		ruledata.Target = b.GetDeploy()
	}

	// return the inverse of container execute
	exec, err := c.Execute(ruledata)

	return !exec, err
}
