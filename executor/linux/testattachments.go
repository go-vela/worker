// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	api "github.com/go-vela/server/api/types"
)

// CreateTestAttachment creates a test attachment record in the database
// after a file has been successfully uploaded to storage.
func (c *client) CreateTestAttachment(ctx context.Context, fileName, presignURL string, size int64, tr *api.TestReport) error {
	// extract file extension and type information
	fileExt := filepath.Ext(fileName)

	// create object path matching the storage upload format
	objectPath := fmt.Sprintf("%s/%s/%s/%s",
		c.build.GetRepo().GetOrg(),
		c.build.GetRepo().GetName(),
		strconv.FormatInt(c.build.GetNumber(), 10),
		fileName)

	// create timestamp for record creation
	createdAt := time.Now().Unix()

	// build test attachment record
	testAttachment := &api.TestAttachment{
		TestReportID: tr.ID, // will be populated by the API based on build context
		FileName:     &fileName,
		ObjectPath:   &objectPath,
		FileSize:     &size,
		FileType:     &fileExt,
		PresignedURL: &presignURL,
		CreatedAt:    &createdAt,
	}

	// update test attachment in database
	ta, resp, err := c.Vela.TestAttachment.Update(
		ctx,
		c.build.GetRepo().GetOrg(),
		c.build.GetRepo().GetName(),
		c.build.GetNumber(),
		testAttachment,
	)
	if err != nil {
		return fmt.Errorf("failed to create test attachment record: build=%d, status=%d, error=%w",
			c.build.GetNumber(), resp.StatusCode, err)
	}

	c.Logger.Debugf("created test attachment record: id=%d, file=%s", ta.GetID(), fileName)

	return nil
}
