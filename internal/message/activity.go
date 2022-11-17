// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package message

import (
	"sync"

	"github.com/go-vela/types/library"
)

type Action interface{}

type AddBuild struct{}

type RemoveBuild struct{}

// BuildActivity is the message used to update activity tracking for
// build that is being executed by the worker.
type BuildActivity struct {
	Build *library.Build
	Action
}

type Activity struct {
	ActiveBuilds []*library.Build
	Mutex        *sync.Mutex
	Channel      chan BuildActivity
}

func (a *Activity) GetBuild(build *library.Build) (*library.Build, int) {
	var _build *library.Build
	idx := -1
	for i, b := range a.ActiveBuilds {
		if b.GetID() == build.GetID() {
			_build = b
			idx = i
		}
	}
	return _build, idx
}

func (a *Activity) AddBuild(build *library.Build) {
	// check activity for incoming build
	_build, idx := a.GetBuild(build)

	// build found
	if _build != nil || idx != -1 {
		return
	}

	// add build
	a.ActiveBuilds = append(a.ActiveBuilds, build)
}

func (a *Activity) RemoveBuild(build *library.Build) {
	// check activity for incoming build
	_build, idx := a.GetBuild(build)

	// build not found
	if _build == nil || idx == -1 {
		return
	}

	// remove build
	a.ActiveBuilds[idx] = a.ActiveBuilds[len(a.ActiveBuilds)-1]
	a.ActiveBuilds = a.ActiveBuilds[:len(a.ActiveBuilds)-1]
}
