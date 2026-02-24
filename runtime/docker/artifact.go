// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/containerd/errdefs"

	"github.com/go-vela/server/compiler/types/pipeline"
)

// PollFileNames searches for files matching the provided patterns within a container.
func (c *client) PollFileNames(ctx context.Context, ctn *pipeline.Container, _step *pipeline.Container) ([]string, error) {
	c.Logger.Tracef("gathering files from container %s", ctn.ID)

	if ctn.Image == "" {
		return nil, nil
	}

	var results []string

	seen := make(map[string]bool)
	paths := _step.Artifacts.Paths
	workspacePrefix := _step.Environment["VELA_WORKSPACE"] + "/"

	for _, pattern := range paths {
		searchDir := extractSearchDir(workspacePrefix + pattern)

		reader, _, err := c.Docker.CopyFromContainer(ctx, ctn.ID, searchDir)
		if err != nil {
			if errdefs.IsNotFound(err) {
				c.Logger.Debugf("search directory %q not found in container: %v", searchDir, err)
				continue
			}

			return nil, fmt.Errorf("copy from container for search dir %q: %w", searchDir, err)
		}

		tr := tar.NewReader(reader)

		for {
			header, err := tr.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				reader.Close()

				return nil, fmt.Errorf("read tar entry for search dir %q: %w", searchDir, err)
			}

			if header.Typeflag != tar.TypeReg {
				continue
			}

			filePath := filepath.Clean(header.Name)

			matched, err := filepath.Match(pattern, filePath)
			if err == nil && matched && !seen[filePath] {
				results = append(results, workspacePrefix+filePath)
				seen[filePath] = true
			}
		}

		reader.Close()
	}

	if len(results) == 0 {
		return results, fmt.Errorf("no matching files found for patterns: %v", paths)
	}

	return results, nil
}

// extractSearchDir extracts a search directory from a glob pattern.
func extractSearchDir(pattern string) string {
	// if no wildcard, determine directory
	idx := strings.IndexAny(pattern, "*?[")
	if idx == -1 {
		return filepath.Dir(pattern)
	}

	// determine directory before wildcard
	dir := filepath.Dir(pattern[:idx])
	if dir == "" || dir == "." {
		return "/"
	}

	return dir
}

// PollFileContent retrieves the content and size of a file inside a container.
func (c *client) PollFileContent(ctx context.Context, ctn *pipeline.Container, path string) (io.Reader, int64, error) {
	c.Logger.Tracef("gathering test results and attachments from container %s", ctn.ID)

	if len(ctn.Image) == 0 || len(path) == 0 {
		return nil, 0, nil
	}

	// copy file from outputs container
	reader, _, err := c.Docker.CopyFromContainer(ctx, ctn.ID, path)
	if err != nil {
		c.Logger.Debugf("PollFileContent CopyFromContainer failed for %q: %v", path, err)
		// early non-error exit if not found
		if errdefs.IsNotFound(err) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	defer reader.Close()

	// docker returns a tar archive for the path
	tr := tar.NewReader(reader)

	header, err := tr.Next()
	if err != nil {
		// if the tar has no entries or is finished unexpectedly
		if errors.Is(err, io.EOF) {
			c.Logger.Debugf("PollFileContent: no tar entries for %q", path)

			return nil, 0, nil
		}

		c.Logger.Debugf("PollFileContent tr.Next failed for %q: %v", path, err)

		return nil, 0, err
	}

	// Ensure the tar entry is a regular file (not dir, symlink, etc.)
	if header.Typeflag != tar.TypeReg {
		c.Logger.Debugf("PollFileContent unexpected tar entry type %v for %q", header.Typeflag, path)

		return nil, 0, fmt.Errorf("unexpected tar entry type %v for %q", header.Typeflag, path)
	}

	// Read file contents. Use io.ReadAll to avoid dealing with CopyN EOF nuances.
	fileBytes, err := io.ReadAll(tr)
	if err != nil {
		c.Logger.Debugf("PollFileContent ReadAll failed for %q: %v", path, err)

		return nil, 0, err
	}

	if len(fileBytes) == 0 {
		c.Logger.Errorf("PollFileContent returned no data for path: %s", path)

		return nil, 0, fmt.Errorf("no data returned from container for %q", path)
	}

	// Return a reader and length (use int64 for size)
	return bytes.NewReader(fileBytes), int64(len(fileBytes)), nil
}
