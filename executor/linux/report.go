// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"
)

// CreateTestReport creates a test report record in the database for the current build.
func (c *client) CreateTestReport(ctx context.Context) (*api.TestReport, error) {
	// create empty test report for the build
	testReport := &api.TestReport{}

	// update test report in database
	tr, resp, err := c.Vela.TestReport.Update(
		ctx,
		c.build.GetRepo().GetOrg(),
		c.build.GetRepo().GetName(),
		c.build.GetNumber(),
		testReport,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create test report record: build=%d, status=%d, error=%w",
			c.build.GetNumber(), resp.StatusCode, err)
	}

	c.Logger.Debugf("created test report record: id=%d, build=%d", tr.GetID(), c.build.GetNumber())

	return tr, nil
}
