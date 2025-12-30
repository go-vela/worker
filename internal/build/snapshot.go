// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// Snapshot creates a moment in time record of the build
// and attempts to upload it to the server.
func Snapshot(ctx context.Context, b *api.Build, c *vela.Client, e error, l *logrus.Entry) {
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
		// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
		l = logrus.NewEntry(logrus.StandardLogger())
	}

	// check if the Vela client provided is empty
	if c != nil {
		l.Debug("uploading build snapshot")

		// send API call to update the build
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#BuildService.Update
		_, _, err := c.Build.Update(ctx, b)
		if err != nil {
			l.Errorf("unable to upload build snapshot: %v", err)
		}
	}
}
