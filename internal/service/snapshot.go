// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Snapshot creates a moment in time record of the
// service and attempts to upload it to the server.
func Snapshot(ctx context.Context, ctn *pipeline.Container, b *api.Build, c *vela.Client, l *logrus.Entry, s *api.Service) {
	// check if the build is not in a canceled status
	if !strings.EqualFold(s.GetStatus(), constants.StatusCanceled) {
		// check if the container is running in headless mode
		if !ctn.Detach {
			// update the service fields to indicate a success
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

			// update the service fields to indicate a failure
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
		l.Debug("uploading service snapshot")

		// send API call to update the service
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#SvcService.Update
		_, _, err := c.Svc.Update(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), s)
		if err != nil {
			l.Errorf("unable to upload service snapshot: %v", err)
		}
	}
}
