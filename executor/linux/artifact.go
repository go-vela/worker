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

// CreateArtifact creates an artifact record in the database
// after a file has been successfully uploaded to storage.
func (c *client) CreateArtifact(ctx context.Context, fileName, presignURL string, size int64) error {
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

	// build artifact record
	artifact := &api.Artifact{
		BuildID:      c.build.ID,
		FileName:     &fileName,
		ObjectPath:   &objectPath,
		FileSize:     &size,
		FileType:     &fileExt,
		PresignedURL: &presignURL,
		CreatedAt:    &createdAt,
	}

	// create artifact record in database
	a, resp, err := c.Vela.Artifact.Update(
		ctx,
		c.build.GetRepo().GetOrg(),
		c.build.GetRepo().GetName(),
		c.build.GetNumber(),
		artifact,
	)
	if err != nil {
		return fmt.Errorf("failed to create artifact record: build=%d, status=%d, error=%w",
			c.build.GetNumber(), resp.StatusCode, err)
	}

	c.Logger.Debugf("created artifact record: id=%d, file=%s", a.GetID(), fileName)

	return nil
}
