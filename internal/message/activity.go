// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package message

import "github.com/go-vela/types/library"

type Action interface{}

type AddBuild struct{}

type RemoveBuild struct{}

// BuildActivity is the message used to update activity tracking for
// build that is being executed by the worker.
type BuildActivity struct {
	Build *library.Build
	Action
}
