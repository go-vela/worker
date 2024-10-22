// SPDX-License-Identifier: Apache-2.0

package step

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Snapshot creates a moment in time record of the
// step and attempts to upload it to the server.
func Snapshot(ctn *pipeline.Container, b *api.Build, c *vela.Client, l *logrus.Entry, s *api.Step) {
	// check if the build is not in a canceled status or error status
	logrus.Debugf("Snapshot s: %s %s", s.GetName(), s.GetStatus())

	if !strings.EqualFold(s.GetStatus(), constants.StatusCanceled) &&
		!strings.EqualFold(s.GetStatus(), constants.StatusError) {
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
		// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
		l = logrus.NewEntry(logrus.StandardLogger())
	}

	// check if the Vela client provided is empty
	if c != nil {
		l.Debug("uploading step snapshot")

		// send API call to update the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#StepService.Update
		_, _, err := c.Step.Update(b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), s)
		if err != nil {
			l.Errorf("unable to upload step snapshot: %v", err)
		}
	}
}

// SnapshotInit creates a moment in time record of the
// init step and attempts to upload it to the server.
func SnapshotInit(ctn *pipeline.Container, b *api.Build, c *vela.Client, l *logrus.Entry, s *api.Step, lg *api.Log) {
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
		// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
		l = logrus.NewEntry(logrus.StandardLogger())
	}

	// check if the Vela client provided is empty
	if c != nil {
		l.Debug("uploading step snapshot")

		// send API call to update the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#StepService.Update
		_, _, err := c.Step.Update(b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), s)
		if err != nil {
			l.Errorf("unable to upload step snapshot: %v", err)
		}

		l.Debug("uploading step logs")

		// send API call to update the logs for the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.UpdateStep
		_, err = c.Log.UpdateStep(b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), s.GetNumber(), lg)
		if err != nil {
			l.Errorf("unable to upload step logs: %v", err)
		}
	}
}
