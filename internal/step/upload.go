// SPDX-License-Identifier: Apache-2.0

package step

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Upload tracks the final state of the step
// and attempts to upload it to the server.
func Upload(ctn *pipeline.Container, b *api.Build, c *vela.Client, l *logrus.Entry, s *api.Step) {
	// handle the step based off the status provided
	switch s.GetStatus() {
	// step is in a canceled state
	case constants.StatusCanceled:
		fallthrough
	// step is in a error state
	case constants.StatusError:
		fallthrough
	// step is in a failure state
	case constants.StatusFailure:
		// if the step is in a canceled, error
		// or failure state we DO NOT want to
		// update the state to be success
		break
	// step is in a pending state
	case constants.StatusPending:
		// if the step is in a pending state
		// then something must have gone
		// drastically wrong because this
		// SHOULD NOT happen
		//
		// TODO: consider making this a constant
		s.SetExitCode(137)
		s.SetFinished(time.Now().UTC().Unix())
		s.SetStatus(constants.StatusKilled)

		// check if the step was not started
		if s.GetStarted() == 0 {
			// set the started time to the finished time
			s.SetStarted(s.GetFinished())
		}
	default:
		// update the step with a success state
		s.SetStatus(constants.StatusSuccess)
	}

	// check if the step finished
	if s.GetFinished() == 0 {
		// update the step with the finished timestamp
		s.SetFinished(time.Now().UTC().Unix())

		// check the container for an unsuccessful exit code
		if ctn.ExitCode != 0 {
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
		l.Debug("uploading final step state")

		// send API call to update the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#StepService.Update
		_, _, err := c.Step.Update(b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), s)
		if err != nil {
			l.Errorf("unable to upload final step state: %v", err)
		}
	}
}
