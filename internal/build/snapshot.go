// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"strings"
	"time"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Snapshot creates a moment in time record of the build
// and attempts to upload it to the server.
func Snapshot(b *library.Build, c *vela.Client, e error, l *logrus.Entry, r *library.Repo) {
	// check if the build is not in a canceled status
	if !strings.EqualFold(b.GetStatus(), constants.StatusCanceled) {
		// check if the error provided is empty
		if e != nil {
			// populate build fields with error based values
			b.SetError(e.Error())
			b.SetStatus(constants.StatusError)
			b.SetFinished(time.Now().UTC().Unix())
		}
	}

	// check if the logger provided is empty
	if l == nil {
		// create new logger
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
		l = logrus.NewEntry(logrus.StandardLogger())
	}

	// check if the Vela client provided is empty
	if c != nil {
		l.Debug("uploading build snapshot")

		// send API call to update the build
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#BuildService.Update
		_, _, err := c.Build.Update(r.GetOrg(), r.GetName(), b)
		if err != nil {
			l.Errorf("unable to upload build snapshot: %v", err)
		}
	}
}
