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
	dockerContainerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// isAllowedExt returns true if ext (".xml", ".png", etc.) is in your allow-list.
func isAllowedExt(ext string) bool {
	ext = strings.ToLower(ext)
	for _, a := range constants.AllAllowedExtensions {
		if ext == a {
			return true
		}
	}

	return false
}

// execContainerLines runs `sh -c cmd` in the named container and
// returns its stdout split by newline (error if anything on stderr).
func (c *client) execContainerLines(ctx context.Context, containerID, cmd string) ([]string, error) {
	execConfig := dockerContainerTypes.ExecOptions{
		Tty:          true,
		Cmd:          []string{"sh", "-c", cmd},
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := c.Docker.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return nil, fmt.Errorf("create exec: %w", err)
	}

	attach, err := c.Docker.ContainerExecAttach(ctx, resp.ID, dockerContainerTypes.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("attach exec: %w", err)
	}

	defer attach.Close()

	var outBuf, errBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&outBuf, &errBuf, attach.Reader); err != nil {
		return nil, fmt.Errorf("copy exec output: %w", err)
	}

	if errBuf.Len() > 0 {
		return nil, fmt.Errorf("exec error: %s", errBuf.String())
	}

	lines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")

	return lines, nil
}

// PollFileNames searches for files matching the provided patterns within a container.
func (c *client) PollFileNames(ctx context.Context, ctn *pipeline.Container, paths []string) ([]string, error) {
	c.Logger.Tracef("gathering files from container %s", ctn.ID)

	if ctn.Image == "" {
		return nil, nil
	}

	var results []string

	for _, pattern := range paths {
		// use find command to locate files matching the pattern
		cmd := fmt.Sprintf("find / -type f -path '*%s' -print", pattern)
		c.Logger.Debugf("searching for files with pattern: %s", pattern)

		lines, err := c.execContainerLines(ctx, ctn.ID, cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to search for pattern %q: %w", pattern, err)
		}

		c.Logger.Tracef("found %d candidates for pattern %s", len(lines), pattern)

		// process each found file
		for _, line := range lines {
			filePath := filepath.Clean(strings.TrimSpace(line))
			if filePath == "" {
				continue
			}

			// check if file extension is allowed
			ext := strings.ToLower(filepath.Ext(filePath))
			if !isAllowedExt(ext) {
				c.Logger.Debugf("skipping file %s (extension %s not allowed)", filePath, ext)
				continue
			}

			c.Logger.Debugf("accepted file: %s", filePath)
			results = append(results, filePath)
		}
	}

	if len(results) == 0 {
		return results, fmt.Errorf("no matching files found for patterns: %v", paths)
	}

	c.Logger.Infof("found %d files matching patterns", len(results))

	return results, nil
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
