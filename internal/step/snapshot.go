// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"strings"
	"time"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// Snapshot creates a moment in time record of the
// step and attempts to upload it to the server.
func Snapshot(ctn *pipeline.Container, b *library.Build, c *vela.Client, l *logrus.Entry, r *library.Repo, s *library.Step) {
	// check if the build is not in a canceled status
	if !strings.EqualFold(s.GetStatus(), constants.StatusCanceled) {
		// check if the container is running in headless mode
		if !ctn.Detach {
			// update the step fields to indicate a success
			s.SetStatus(constants.StatusSuccess)
			s.SetFinished(time.Now().UTC().Unix())
		}

		// check if the container has an unsuccessful exit code
		if ctn.ExitCode != 0 {
			// check if container failures should be ignored
			if !ctn.Ruleset.Continue {
				// set build status to failure
				b.SetStatus(constants.StatusFailure)
			}

			// update the step fields to indicate a failure
			s.SetExitCode(ctn.ExitCode)
			s.SetStatus(constants.StatusFailure)
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
		l.Debug("uploading step snapshot")

		// send API call to update the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.Update
		_, _, err := c.Step.Update(r.GetOrg(), r.GetName(), b.GetNumber(), s)
		if err != nil {
			l.Errorf("unable to upload step snapshot: %v", err)
		}
	}
}

// SnapshotInit creates a moment in time record of the
// init step and attempts to upload it to the server.
func SnapshotInit(ctn *pipeline.Container, b *library.Build, c *vela.Client, l *logrus.Entry, r *library.Repo, s *library.Step, lg *library.Log) {
	// check if the build is not in a canceled status
	if !strings.EqualFold(s.GetStatus(), constants.StatusCanceled) {
		// check if the container has an unsuccessful exit code
		if ctn.ExitCode != 0 {
			// check if container failures should be ignored
			if !ctn.Ruleset.Continue {
				// set build status to failure
				b.SetStatus(constants.StatusFailure)
			}

			// update the step fields to indicate a failure
			s.SetExitCode(ctn.ExitCode)
			s.SetStatus(constants.StatusFailure)
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
		l.Debug("uploading step snapshot")

		// send API call to update the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.Update
		_, _, err := c.Step.Update(r.GetOrg(), r.GetName(), b.GetNumber(), s)
		if err != nil {
			l.Errorf("unable to upload step snapshot: %v", err)
		}

		l.Debug("uploading step logs")

		// send API call to update the logs for the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#LogService.UpdateStep
		_, _, err = c.Log.UpdateStep(r.GetOrg(), r.GetName(), b.GetNumber(), s.GetNumber(), lg)
		if err != nil {
			l.Errorf("unable to upload step logs: %v", err)
		}
	}
}
