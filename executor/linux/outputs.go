// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	envparse "github.com/hashicorp/go-envparse"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
)

// outputSvc handles communication with the outputs container during the build.
type outputSvc svc

// create configures the outputs container for execution.
func (o *outputSvc) create(ctx context.Context, ctn *pipeline.Container, timeout int64) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// set up outputs logger
	logger := o.client.Logger.WithField("outputs", "outputs")

	// Encode script content to Base64
	script := base64.StdEncoding.EncodeToString(
		fmt.Appendf(nil, "mkdir /vela/outputs\nsleep %d\n", timeout),
	)

	// set the entrypoint for the ctn
	ctn.Entrypoint = []string{"/bin/sh", "-c"}

	// set the commands for the ctn
	ctn.Commands = []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"}

	// set the environment variables for the ctn
	ctn.Environment["HOME"] = "/root"
	ctn.Environment["SHELL"] = "/bin/sh"
	ctn.Environment["VELA_BUILD_SCRIPT"] = script

	logger.Debug("setting up outputs container")
	// setup the runtime container
	err := o.client.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// destroy cleans up outputs container after execution.
func (o *outputSvc) destroy(ctx context.Context, ctn *pipeline.Container) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// update engine logger with outputs metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("outputs", ctn.Name)

	logger.Debug("inspecting outputs container")
	// inspect the runtime container
	err := o.client.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("removing outputs container")
	// remove the runtime container
	err = o.client.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// exec runs the outputs sidecar container for a pipeline.
func (o *outputSvc) exec(ctx context.Context, _outputs *pipeline.Container) error {
	// exit if outputs container has not been configured
	if len(_outputs.Image) == 0 {
		return nil
	}

	logrus.Debug("running outputs container")
	// run the runtime container
	err := o.client.Runtime.RunContainer(ctx, _outputs, o.client.pipeline)
	if err != nil {
		return err
	}

	logrus.Debug("inspecting outputs container")
	// inspect the runtime container
	err = o.client.Runtime.InspectContainer(ctx, _outputs)
	if err != nil {
		return err
	}

	return nil
}

// poll tails the output for sidecar container.
func (o *outputSvc) poll(ctx context.Context, ctn *pipeline.Container) (map[string]string, map[string]string, error) {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil, nil, nil
	}

	// update engine logger with outputs metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("outputs", ctn.Name)

	logger.Debug("tailing container")

	var (
		filePaths = []string{
			"/vela/outputs/.env",
			"/vela/outputs/masked.env",
			"/vela/outputs/base64.env",
			"/vela/outputs/masked.base64.env",
		}

		outputMap = make(map[string]string)
		maskMap   = make(map[string]string)
	)

	for _, p := range filePaths {
		outputBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, p)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to poll outputs container %s: %w", ctn.Name, err)
		}

		reader := bytes.NewReader(outputBytes)

		switch p {
		case "/vela/outputs/.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse output map: %v", err)
			}

			// add to output map
			maps.Copy(outputMap, parsed)

		case "/vela/outputs/masked.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse masked output map: %v", err)
			}

			// add to mask map
			maps.Copy(maskMap, parsed)

		case "/vela/outputs/base64.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse base64 output map: %v", err)
			}

			for k, v := range parsed {
				// decode the base64 value
				decodedValue, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					logger.Debugf("unable to decode base64 value for key %s: %v", k, err)
					continue
				}

				parsed[k] = string(decodedValue)
			}

			// add to output map
			maps.Copy(outputMap, parsed)

		case "/vela/outputs/masked.base64.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse masked base64 output map: %v", err)
			}

			for k, v := range parsed {
				// decode the base64 value
				decodedValue, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					logger.Debugf("unable to decode base64 value for key %s: %v", k, err)
					continue
				}

				parsed[k] = string(decodedValue)
			}

			// add to mask map
			maps.Copy(maskMap, parsed)
		}
	}

	return outputMap, maskMap, nil
}

// pollFiles polls the output for files from the sidecar container.
func (o *outputSvc) pollFiles(ctx context.Context, ctn *pipeline.Container, _step *pipeline.Container, b *api.Build) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return fmt.Errorf("no outputs container configured")
	}

	// update engine logger with outputs metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("artifact-outputs", ctn.Name)

	// load the step log so artifact messages are streamed to the UI.
	// Wait briefly to allow the StreamStep goroutine to finish its final
	// log upload so that we do not race with it. Then re-fetch the log
	// from the server to get the authoritative copy that includes all
	// container output — this prevents artifact messages from overwriting
	// the step's own logs.
	time.Sleep(2 * time.Second)

	_log, _, err := o.client.Vela.Log.GetStep(ctx,
		b.GetRepo().GetOrg(), b.GetRepo().GetName(),
		b.GetNumber(), _step.Number)
	if err != nil {
		logger.Warnf("unable to fetch step log for artifact streaming: %v", err)
	}

	// store back into the map so future references are consistent
	if _log != nil {
		o.client.stepLogs.Store(_step.ID, _log)
	}

	// streamLog appends a message to the step log and pushes it to the server
	// so that artifact progress is visible in the UI on the attached step.
	streamLog := func(msg string) {
		logger.Info(msg)

		if _log == nil {
			return
		}

		_log.AppendData([]byte(fmt.Sprintf("[artifact] %s\n", msg)))

		_, err := o.client.Vela.Log.UpdateStep(ctx,
			b.GetRepo().GetOrg(), b.GetRepo().GetName(),
			b.GetNumber(), _step.Number, _log)
		if err != nil {
			logger.Errorf("unable to update step log: %v", err)
		}
	}

	streamLog(fmt.Sprintf("starting artifact upload for step %s (build: %d, repo: %s/%s)",
		_step.Name, b.GetNumber(), b.GetRepo().GetOrg(), b.GetRepo().GetName()))
	streamLog(fmt.Sprintf("configured artifact paths: %v", _step.Artifacts.Paths))

	// grab file paths from the container
	filesPath, err := o.client.Runtime.PollFileNames(ctx, ctn, _step)
	if err != nil {
		streamLog(fmt.Sprintf("failed to discover artifact files: %v", err))
		return fmt.Errorf("unable to poll file names: %w", err)
	}

	streamLog(fmt.Sprintf("discovered %d artifact file(s) matching configured paths", len(filesPath)))

	if len(filesPath) == 0 {
		streamLog(fmt.Sprintf("no files found matching artifact paths: %v — ensure your step produces files at the expected locations", _step.Artifacts.Paths))
		return fmt.Errorf("no files found for file list: %v", _step.Artifacts.Paths)
	}

	logger.Debugf("matched files: %v", filesPath)

	// create http client for uploading files to storage
	putClient := http.DefaultClient
	putClient.Timeout = time.Second * 30

	// track upload statistics
	var (
		uploaded int
		skipped  int
		failed   int
	)

	// process each file found
	for _, filePath := range filesPath {
		fileName := filepath.Base(filePath)
		logger.Debugf("processing file: %s (path: %s)", fileName, filePath)

		// skip hidden files and files within hidden directories
		if isHidden(filePath) {
			logger.Debugf("skipping hidden file or directory: %s", filePath)
			skipped++

			continue
		}

		url, _, err := o.client.Vela.Build.GetPresignedPutURL(ctx, fileName, b.GetRepo().GetOrg(), b.GetRepo().GetName(),
			b.GetNumber())
		if err != nil {
			streamLog(fmt.Sprintf("artifact %q could not be uploaded — the server did not provide an upload URL. "+
				"This may indicate that artifact storage is not configured or the server encountered an error. "+
				"Please contact your Vela administrator if this persists. (error: %v)", fileName, err))
			failed++

			continue
		}

		// get file content from container
		reader, size, err := o.client.Runtime.PollFileContent(ctx, ctn, filePath)
		if err != nil {
			streamLog(fmt.Sprintf("unable to read artifact file %q from container: %v", filePath, err))
			failed++

			continue
		}

		logger.Infof("artifact file %q size: %d bytes", fileName, size)

		// TODO: surface this skip to the user
		if o.client.fileSizeLimit > 0 && size > o.client.fileSizeLimit {
			streamLog(fmt.Sprintf("skipping artifact %q — file size (%d bytes) exceeds the per-file limit (%d bytes)",
				fileName, size, o.client.fileSizeLimit))
			skipped++

			continue
		}

		if o.client.buildFileSizeLimit > 0 && size+o.client.Uploaded > o.client.buildFileSizeLimit {
			streamLog(fmt.Sprintf("skipping artifact %q — uploading this file would exceed the per-build size limit (%d bytes). "+
				"Total uploaded so far: %d bytes", fileName, o.client.buildFileSizeLimit, o.client.Uploaded))
			skipped++

			continue
		}

		// create storage object path
		objectName := fmt.Sprintf("%s/%s/%s/%s",
			b.GetRepo().GetOrg(),
			b.GetRepo().GetName(),
			strconv.FormatInt(b.GetNumber(), 10),
			fileName)

		streamLog(fmt.Sprintf("uploading artifact %q to storage (object: %s, size: %d bytes)", fileName, objectName, size))

		err = uploadObject(ctx, putClient, reader, size, fileName, url.URL)
		if err != nil {
			streamLog(fmt.Sprintf("failed to upload artifact %q: %v", fileName, err))
			failed++

			continue
		}

		o.client.Uploaded += size
		uploaded++

		streamLog(fmt.Sprintf("successfully uploaded artifact %q (%d bytes)", fileName, size))
	}

	streamLog(fmt.Sprintf("artifact upload complete — uploaded: %d, skipped: %d, failed: %d (total files: %d)",
		uploaded, skipped, failed, len(filesPath)))

	return nil
}

// isHidden reports whether any component of the given path (file or directory)
// starts with a ".", which indicates a hidden file or directory.
func isHidden(path string) bool {
	for part := range strings.SplitSeq(filepath.ToSlash(path), "/") {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}

	return false
}

// uploadObject uploads an object to a bucket in MinIO.ts.
func uploadObject(ctx context.Context, putClient *http.Client, reader io.Reader, size int64, filename, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reader)
	if err != nil {
		return fmt.Errorf("could not create PUT request: %w", err)
	}

	// Set the Content-Type header based on the file extension
	ext := filepath.Ext(filename)

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req.Header.Set("Content-Type", contentType)

	// Set the Content-Length header
	req.ContentLength = size

	// Perform the HTTP request to upload the object
	//
	//nolint:bodyclose // body closes on line 310
	resp, err := putClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not upload data to bucket: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Warnf("could not close response body: %v", err)
		}
	}(resp.Body)

	// Check for a successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
