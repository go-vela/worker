// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"

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

// PollFileNames grabs files name from provided path
// within a container and uploads them to s3
func (c *client) PollFileNames(ctx context.Context, ctn *pipeline.Container, paths []string) ([]string, error) {
	c.Logger.Tracef("gathering files from container %s", ctn.ID)
	if ctn.Image == "" {
		return nil, nil
	}

	results := make([]string, 0)
	for _, p := range paths {
		// use find on the container to locate candidates
		cmd := fmt.Sprintf("find / -type f -path '*%s' -print", p)
		c.Logger.Infof("running: %s", cmd)

		lines, err := c.execContainerLines(ctx, ctn.ID, cmd)
		if err != nil {
			return nil, fmt.Errorf("searching %q: %w", p, err)
		}
		c.Logger.Tracef("candidates: %v", lines)

		for _, raw := range lines {
			// 1) strip whitespace (including CR/LF) and normalize the path
			fp := filepath.Clean(strings.TrimSpace(raw))
			if fp == "" {
				continue
			}

			// 2) quick ext check
			ext := strings.ToLower(filepath.Ext(fp))
			if !isAllowedExt(ext) {
				c.Logger.Infof("skipping %s (ext %s not allowed)", fp, ext)
				continue
			}

			// 3) accept
			c.Logger.Infof("accepted file: %s", fp)
			results = append(results, fp)
		}
	}

	// Add logging for empty results in PollFileNames
	if len(results) == 0 {
		c.Logger.Errorf("PollFileNames found no matching files for paths: %v", paths)
		return results, fmt.Errorf("no matching files found for any paths %v", paths)
	}
	return results, nil

}

// PollFileContent retrieves the content and size of a file inside a container.
func (c *client) PollFileContent(ctx context.Context, ctn *pipeline.Container, path string) (io.Reader, int64, error) {
	c.Logger.Tracef("gathering test results and attachments from container %s", ctn.ID)

	if len(ctn.Image) == 0 {
		// return an empty reader instead of nil
		return bytes.NewReader(nil), 0, fmt.Errorf("empty container image")
	}
	cmd := []string{"sh", "-c", fmt.Sprintf("base64 %s", path)}
	execConfig := dockerContainerTypes.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: false,
		Tty:          false,
	}

	c.Logger.Infof("executing command for content: %v", execConfig.Cmd)
	execID, err := c.Docker.ContainerExecCreate(ctx, ctn.ID, execConfig)
	if err != nil {
		c.Logger.Debugf("PollFileContent exec-create failed for %q: %v", path, err)
		return nil, 0, fmt.Errorf("failed to create exec instance: %w", err)
	}
	resp, err := c.Docker.ContainerExecAttach(ctx, execID.ID, dockerContainerTypes.ExecAttachOptions{})
	if err != nil {
		c.Logger.Debugf("PollFileContent exec-attach failed for %q: %v", path, err)
		return nil, 0, fmt.Errorf("failed to attach to exec instance: %w", err)
	}

	defer func() {
		if resp.Conn != nil {
			resp.Close()
		}
	}()

	outputStdout := new(bytes.Buffer)
	outputStderr := new(bytes.Buffer)

	if resp.Reader != nil {
		_, err := stdcopy.StdCopy(outputStdout, outputStderr, resp.Reader)
		if err != nil {
			c.Logger.Errorf("unable to copy logs for container: %v", err)
		}
	}

	if outputStderr.Len() > 0 {
		return nil, 0, fmt.Errorf("error: %s", outputStderr.String())
	}
	data := outputStdout.Bytes()

	// Add logging for empty data in PollFileContent
	if len(data) == 0 {
		c.Logger.Errorf("PollFileContent returned no data for path: %s", path)
		return nil, 0, fmt.Errorf("no data returned from base64 command")
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		c.Logger.Errorf("unable to decode base64 data: %v", err)
		return nil, 0, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return bytes.NewReader(decoded), int64(len(decoded)), nil
}
