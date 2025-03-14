// SPDX-License-Identifier: Apache-2.0

package build

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// Upload tracks the final state of the build
// and attempts to upload it to the server.
func Upload(b *api.Build, c *vela.Client, e error, l *logrus.Entry) {
	// handle the build based off the status provided
	switch b.GetStatus() {
	// build is in a canceled state
	case constants.StatusCanceled:
		fallthrough
	// build is in a error state
	case constants.StatusError:
		fallthrough
	// build is in a failure state
	case constants.StatusFailure:
		// if the build is in a canceled, error
		// or failure state we DO NOT want to
		// update the state to be success
		break
	// build is in a pending state
	case constants.StatusPending:
		// if the build is in a pending state
		// then something must have gone
		// drastically wrong because this
		// SHOULD NOT happen
		b.SetStatus(constants.StatusKilled)
	default:
		// update the build with a success state
		b.SetStatus(constants.StatusSuccess)
	}

	// check if the build is not in a canceled status
	if !strings.EqualFold(b.GetStatus(), constants.StatusCanceled) {
		// check if the error provided is empty
		if e != nil {
			// update the build with error based values
			b.SetError(e.Error())
			b.SetStatus(constants.StatusError)
		}
	}

	// update the build with the finished timestamp
	b.SetFinished(time.Now().UTC().Unix())

	// check if the logger provided is empty
	if l == nil {
		// create new logger
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
		l = logrus.NewEntry(logrus.StandardLogger())
	}

	// check if the Vela client provided is empty
	if c != nil {
		l.Debug("uploading final build state")

		// send API call to update the build
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#BuildService.Update
		_, _, err := c.Build.Update(b)
		if err != nil {
			l.Errorf("unable to upload final build state: %v", err)
		}
	}
}
