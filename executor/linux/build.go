// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/worker/internal/build"
	context2 "github.com/go-vela/worker/internal/context"
	"github.com/go-vela/worker/internal/image"
	"github.com/go-vela/worker/internal/outputs"
	"github.com/go-vela/worker/internal/step"
)

// CreateBuild configures the build for execution.
func (c *client) CreateBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, c.Vela, c.err, c.Logger) }()

	// update the build fields
	c.build.SetStatus(constants.StatusRunning)
	c.build.SetStarted(time.Now().UTC().Unix())
	c.build.SetHost(c.Hostname)
	c.build.SetDistribution(c.Driver())
	c.build.SetRuntime(c.Runtime.Driver())

	c.Logger.Info("uploading build state")
	// send API call to update the build
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#BuildService.Update
	c.build, _, c.err = c.Vela.Build.Update(c.build)
	if c.err != nil {
		return fmt.Errorf("unable to upload build state: %w", c.err)
	}

	// setup the runtime build
	c.err = c.Runtime.SetupBuild(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to setup build %s: %w", c.pipeline.ID, c.err)
	}

	// load the init step from the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadInit
	c.init, c.err = step.LoadInit(c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to load init step from pipeline: %w", c.err)
	}

	c.Logger.Infof("creating %s step", c.init.Name)
	// create the step
	c.err = c.CreateStep(ctx, c.init)
	if c.err != nil {
		return fmt.Errorf("unable to create %s step: %w", c.init.Name, c.err)
	}

	c.Logger.Infof("planning %s step", c.init.Name)
	// plan the step
	c.err = c.PlanStep(ctx, c.init)
	if c.err != nil {
		return fmt.Errorf("unable to plan %s step: %w", c.init.Name, c.err)
	}

	return c.err
}

// PlanBuild prepares the build for execution.
func (c *client) PlanBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, c.Vela, c.err, c.Logger) }()

	// load the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_init, err := step.Load(c.init, &c.steps)
	if err != nil {
		return err
	}

	// load the logs for the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(c.init, &c.stepLogs)
	if err != nil {
		return err
	}

	// put worker information into init logs
	_log.AppendData([]byte(fmt.Sprintf("> Worker Information:\n Host: %s\n Version: %s\n Runtime: %s\n", c.Hostname, c.Version, c.Runtime.Driver())))

	// defer taking a snapshot of the init step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#SnapshotInit
	defer func() {
		if c.err != nil {
			_init.SetStatus(constants.StatusFailure)
		}

		step.SnapshotInit(c.init, c.build, c.Vela, c.Logger, _init, _log)
	}()

	c.Logger.Info("creating network")
	// create the runtime network for the pipeline
	c.err = c.Runtime.CreateNetwork(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to create network: %w", c.err)
	}

	// update the init log with progress
	_log.AppendData([]byte("> Inspecting runtime network...\n"))

	// inspect the runtime network for the pipeline
	network, err := c.Runtime.InspectNetwork(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect network: %w", err)
	}

	// update the init log with network information
	_log.AppendData(network)

	c.Logger.Info("creating volume")
	// create the runtime volume for the pipeline
	c.err = c.Runtime.CreateVolume(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to create volume: %w", c.err)
	}

	// update the init log with progress
	_log.AppendData([]byte("> Inspecting runtime volume...\n"))

	// inspect the runtime volume for the pipeline
	volume, err := c.Runtime.InspectVolume(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect volume: %w", err)
	}

	// update the init log with volume information
	_log.AppendData(volume)

	// update the init log with progress
	_log.AppendData([]byte("> Preparing secrets...\n"))

	// iterate through each secret provided in the pipeline
	for _, secret := range c.pipeline.Secrets {
		// ignore pulling secrets coming from plugins
		if !secret.Origin.Empty() {
			continue
		}

		// only pull in secrets that are set to be pulled in at the start
		if strings.EqualFold(secret.Pull, constants.SecretPullStep) {
			_log.AppendData([]byte(fmt.Sprintf("> Skipping pull: secret <%s> lazy loaded\n", secret.Name)))

			continue
		}

		c.Logger.Infof("pulling secret: %s", secret.Name)

		s, err := c.secret.pull(secret)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to pull secrets: %w", err)
		}

		_log.AppendData([]byte(
			fmt.Sprintf("$ vela view secret --secret.engine %s --secret.type %s --org %s --repo %s --name %s \n",
				secret.Engine, secret.Type, s.GetOrg(), s.GetRepo(), s.GetName())))

		sRaw, err := json.MarshalIndent(s.Sanitize(), "", " ")
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to decode secret: %w", err)
		}

		_log.AppendData(append(sRaw, "\n"...))

		// add secret to the appropriate map
		if s.GetAllowSubstitution() {
			c.Secrets[secret.Name] = s
		} else {
			c.NoSubSecrets[secret.Name] = s
		}
	}

	// escape newlines in secrets loaded on build_start
	escapeNewlineSecrets(c.Secrets)

	return nil
}

// AssembleBuild prepares the containers within a build for execution.
//
//nolint:funlen // consider abstracting parts here but for now this is fine
func (c *client) AssembleBuild(ctx context.Context) error {
	// defer taking a snapshot of the build
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Snapshot
	defer func() { build.Snapshot(c.build, c.Vela, c.err, c.Logger) }()

	// load the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Load
	_init, err := step.Load(c.init, &c.steps)
	if err != nil {
		return err
	}

	// load the logs for the init step from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#LoadLogs
	_log, err := step.LoadLogs(c.init, &c.stepLogs)
	if err != nil {
		return err
	}

	// defer an upload of the init step
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Upload
	defer func() {
		if c.err != nil {
			_init.SetStatus(constants.StatusFailure)
		}

		step.Upload(c.init, c.build, c.Vela, c.Logger, _init)
	}()

	defer func() {
		c.Logger.Infof("uploading %s step logs", c.init.Name)
		// send API call to update the logs for the step
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.UpdateStep
		_, err = c.Vela.Log.UpdateStep(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), c.init.Number, _log)
		if err != nil {
			c.Logger.Errorf("unable to upload %s logs: %v", c.init.Name, err)
		}
	}()

	// update the init log with progress
	_log.AppendData([]byte("> Preparing service images...\n"))

	// create the services for the pipeline
	for _, s := range c.pipeline.Services {
		// TODO: remove this; but we need it for tests
		s.Detach = true

		c.Logger.Infof("creating %s service", s.Name)

		_log.AppendData([]byte(fmt.Sprintf("> Preparing service image %s...\n", s.Image)))

		// create the service
		c.err = c.CreateService(ctx, s)
		if c.err != nil {
			return fmt.Errorf("unable to create %s service: %w", s.Name, c.err)
		}

		c.Logger.Infof("inspecting %s service", s.Name)
		// inspect the service image
		image, err := c.Runtime.InspectImage(ctx, s)
		if err != nil {
			c.err = err

			return fmt.Errorf("unable to inspect %s service: %w", s.Name, err)
		}

		// update the init log with service image info
		_log.AppendData(image)
	}

	// update the init log with progress
	_log.AppendData([]byte("> Preparing stage images...\n"))

	// create the stages for the pipeline
	for _, s := range c.pipeline.Stages {
		//
		if s.Name == constants.InitName {
			continue
		}

		c.Logger.Infof("creating %s stage", s.Name)
		// create the stage
		c.err = c.CreateStage(ctx, s)
		if c.err != nil {
			return fmt.Errorf("unable to create %s stage: %w", s.Name, c.err)
		}
	}

	// update the init log with progress
	_log.AppendData([]byte("> Preparing step images...\n"))

	// create the steps for the pipeline
	for _, s := range c.pipeline.Steps {
		if s.Name == constants.InitName {
			continue
		}

		c.Logger.Infof("creating %s step", s.Name)

		_log.AppendData([]byte(fmt.Sprintf("> Preparing step image %s...\n", s.Image)))

		// create the step
		c.err = c.CreateStep(ctx, s)
		if c.err != nil {
			return fmt.Errorf("unable to create %s step: %w", s.Name, c.err)
		}

		c.Logger.Infof("inspecting %s step", s.Name)
		// inspect the step image
		image, err := c.Runtime.InspectImage(ctx, s)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s step: %w", s.Name, c.err)
		}

		// update the init log with step image info
		_log.AppendData(image)
	}

	// update the init log with progress
	_log.AppendData([]byte("> Preparing secret images...\n"))

	// create the secrets for the pipeline
	for _, s := range c.pipeline.Secrets {
		// skip over non-plugin secrets
		if s.Origin.Empty() {
			continue
		}

		// verify secret image is allowed to run
		if c.enforceTrustedRepos {
			priv, err := image.IsPrivilegedImage(s.Origin.Image, c.privilegedImages)
			if err != nil {
				return err
			}

			if priv && !c.build.GetRepo().GetTrusted() {
				return fmt.Errorf("attempting to use privileged image (%s) as untrusted repo", s.Origin.Image)
			}
		}

		c.Logger.Infof("creating %s secret", s.Origin.Name)

		// fetch request token if id_request used in origin config
		var requestToken string

		if len(s.Origin.IDRequest) > 0 {
			opts := &vela.RequestTokenOptions{
				Image:    s.Origin.Image,
				Request:  s.Origin.IDRequest,
				Commands: len(s.Origin.Commands) > 0 || len(s.Origin.Entrypoint) > 0,
			}

			tkn, _, err := c.Vela.Build.GetIDRequestToken(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), opts)
			if err != nil {
				return err
			}

			requestToken = tkn.GetToken()
		}

		// create the service
		c.err = c.secret.create(ctx, s.Origin, requestToken)
		if c.err != nil {
			return fmt.Errorf("unable to create %s secret: %w", s.Origin.Name, c.err)
		}

		c.Logger.Infof("inspecting %s secret", s.Origin.Name)
		// inspect the service image
		image, err := c.Runtime.InspectImage(ctx, s.Origin)
		if err != nil {
			c.err = err
			return fmt.Errorf("unable to inspect %s secret: %w", s.Origin.Name, err)
		}

		// update the init log with secret image info
		_log.AppendData(image)
	}

	// create outputs container with a timeout equal to the repo timeout
	c.err = c.outputs.create(ctx, c.OutputCtn, int64(60*c.build.GetRepo().GetTimeout()))
	if c.err != nil {
		return fmt.Errorf("unable to create outputs container: %w", c.err)
	}

	// inspect the runtime build (eg a kubernetes pod) for the pipeline
	buildOutput, err := c.Runtime.InspectBuild(ctx, c.pipeline)
	if err != nil {
		c.err = err
		return fmt.Errorf("unable to inspect build: %w", err)
	}

	if len(buildOutput) > 0 {
		// update the init log with progress
		// (an empty value allows the runtime to opt out of providing this)
		_log.AppendData(buildOutput)
	}

	// assemble runtime build just before any containers execute
	c.err = c.Runtime.AssembleBuild(ctx, c.pipeline)
	if c.err != nil {
		return fmt.Errorf("unable to assemble runtime build %s: %w", c.pipeline.ID, c.err)
	}

	// update the init log with progress
	_log.AppendData([]byte("> Executing secret images...\n"))

	return c.err
}

// ExecBuild runs a pipeline for a build.
//
//nolint:funlen // there is a lot going on here and will probably always be long
func (c *client) ExecBuild(ctx context.Context) error {
	defer func() {
		// Exec* calls are responsible for sending StreamRequest messages.
		// close the channel at the end of ExecBuild to signal that
		// nothing else will send more StreamRequest messages.
		close(c.streamRequests)

		// defer an upload of the build
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/build#Upload
		build.Upload(c.build, c.Vela, c.err, c.Logger)
	}()

	// output maps for dynamic environment variables captured from volume
	var opEnv, maskEnv map[string]string

	// test report object for storing the test report information
	var tr *api.TestReport

	// Flag to track if we've already created the test report record
	testReportCreated := false

	// fire up output container to run with the build
	c.Logger.Infof("creating outputs container %s", c.OutputCtn.ID)

	// execute outputs container
	c.err = c.outputs.exec(ctx, c.OutputCtn)
	if c.err != nil {
		return fmt.Errorf("unable to exec outputs container: %w", c.err)
	}

	c.Logger.Info("executing secret images")
	// execute the secret
	c.err = c.secret.exec(ctx, &c.pipeline.Secrets)
	if c.err != nil {
		return fmt.Errorf("unable to execute secret: %w", c.err)
	}

	// execute the services for the pipeline
	for _, _service := range c.pipeline.Services {
		c.Logger.Infof("planning %s service", _service.Name)
		// plan the service
		c.err = c.PlanService(ctx, _service)
		if c.err != nil {
			return fmt.Errorf("unable to plan service: %w", c.err)
		}

		c.Logger.Infof("executing %s service", _service.Name)
		// execute the service
		c.err = c.ExecService(ctx, _service)
		if c.err != nil {
			return fmt.Errorf("unable to execute service: %w", c.err)
		}
	}

	// execute the steps for the pipeline
	for _, _step := range c.pipeline.Steps {
		if _step.Name == constants.InitName {
			continue
		}

		// poll outputs
		opEnv, maskEnv, c.err = c.outputs.poll(ctx, c.OutputCtn)
		if c.err != nil {
			return fmt.Errorf("unable to exec outputs container: %w", c.err)
		}

		opEnv = outputs.Sanitize(_step, opEnv)
		maskEnv = outputs.Sanitize(_step, maskEnv)

		// merge env from outputs
		//
		//nolint:errcheck // only errors with empty environment input, which does not matter here
		_step.MergeEnv(opEnv)

		// merge env from masked outputs
		//
		//nolint:errcheck // only errors with empty environment input, which does not matter here
		_step.MergeEnv(maskEnv)

		// check if the step should be skipped
		//
		// https://pkg.go.dev/github.com/go-vela/worker/internal/step#Skip
		skip, err := step.Skip(_step, c.build, c.build.GetStatus(), c.Storage)
		if err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		if skip {
			continue
		}

		// Check if this step has test_report and storage is disabled
		//if !_step.TestReport.Empty() && c.Storage == nil {
		//	c.Logger.Infof("skipping %s step: storage is disabled but test_report is defined", _step.Name)
		//
		//	//// Load step model
		//	//stepData, err := step.Load(_step, &c.steps)
		//	//if err != nil {
		//	//	return fmt.Errorf("unable to load step: %w", err)
		//	//}
		//	//
		//	//// Load or create logs for this step
		//	////stepLog, err := step.LoadLogs(_step, &c.stepLogs)
		//	////if err != nil {
		//	////	return fmt.Errorf("unable to load step logs: %w", err)
		//	////}
		//	//
		//	//// Ensure timestamps
		//	//now := time.Now().UTC().Unix()
		//	//if stepData.GetStarted() == 0 {
		//	//	stepData.SetStarted(now)
		//	//}
		//	//
		//	//stepData.SetStatus(constants.StatusError)
		//	//stepData.SetExitCode(0)
		//	//stepData.SetFinished(now)
		//
		//	// send API call to update the step
		//	//
		//	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#StepService.Update
		//	//_tsstep, _, err := c.Vela.Step.Update(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), stepData)
		//	//if err != nil {
		//	//	return err
		//	//}
		//	//
		//	//// send API call to capture the step log
		//	////
		//	//// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#LogService.GetStep
		//	//_log, _, err := c.Vela.Log.GetStep(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), _tsstep.GetNumber())
		//	//if err != nil {
		//	//	return err
		//	//}
		//	//_log.AppendData([]byte("Storage is disabled, contact Vela Admins\n"))
		//	//
		//	//// add a step log to a map
		//	//c.stepLogs.Store(_step.ID, _log)
		//	//stepLog.AppendData([]byte("Storage is disabled, contact Vela Admins\n"))
		//	//stepLog.SetData([]byte("Storage is disabled, contact Vela Admins\n"))
		//	//// Upload logs so UI can display the message
		//	//if _, err := c.Vela.Log.
		//	//	UpdateStep(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), *stepData.Number, stepLog); err != nil {
		//	//	c.Logger.Errorf("unable to upload skipped step logs: %v", err)
		//	//}
		//	// Upload step status
		//	//step.Upload(_step, c.build, c.Vela, c.Logger, stepData)
		//
		//	continue
		//}

		// add netrc to secrets for masking in logs
		sec := &pipeline.StepSecret{
			Target: "VELA_NETRC_PASSWORD",
		}
		_step.Secrets = append(_step.Secrets, sec)

		// load any lazy secrets into the container environment
		c.err = loadLazySecrets(c, _step)
		if c.err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		c.Logger.Infof("planning %s step", _step.Name)
		// plan the step
		c.err = c.PlanStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

		// add masked outputs to secret map so they can be masked in logs
		for key := range maskEnv {
			sec := &pipeline.StepSecret{
				Target: key,
			}
			_step.Secrets = append(_step.Secrets, sec)
		}

		// logic for polling files only if the test-report step is present
		// iterate through the steps in the build

		// TODO: API to return if storage is enabled
		//if c.Storage == nil && _step.TestReport.Empty() || c.Storage == nil && !_step.TestReport.Empty() {
		//	c.Logger.Infof("storage disabled, skipping test report for %s step", _step.Name)
		//	// skip if no storage client
		//	// but test report is defined in step
		//	continue
		//} else if !_step.TestReport.Empty() && c.Storage != nil {
		c.Logger.Debug("creating test report record in database")
		// send API call to update the test report
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#TestReportService.Add
		// TODO: .Add should be .Update
		// TODO: handle somewhere if multiple test report keys exist in pipeline
		if !testReportCreated {
			tr, c.err = c.CreateTestReport()
			if c.err != nil {
				return fmt.Errorf("unable to create test report: %w", c.err)
			}

			testReportCreated = true
		}

		if len(_step.TestReport.Results) != 0 {
			err := c.outputs.pollFiles(ctx, c.OutputCtn, _step.TestReport.Results, c.build, tr)
			if err != nil {
				c.Logger.Errorf("unable to poll files for results: %v", err)
			}
		}

		if len(_step.TestReport.Attachments) != 0 {
			err := c.outputs.pollFiles(ctx, c.OutputCtn, _step.TestReport.Attachments, c.build, tr)
			if err != nil {
				c.Logger.Errorf("unable to poll files for attachments: %v", err)
			}
		}
		//}

		// perform any substitution on dynamic variables
		err = _step.Substitute()
		if err != nil {
			return err
		}

		// inject no-substitution secrets for container
		err = injectSecrets(_step, c.NoSubSecrets)
		if err != nil {
			return err
		}

		c.Logger.Infof("executing %s step", _step.Name)
		// execute the step
		c.err = c.ExecStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to execute step: %w", c.err)
		}
	}

	// create an error group with the context for each stage
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#WithContext
	stages, stageCtx := errgroup.WithContext(ctx)

	// create a map to track the progress of each stage
	stageMap := new(sync.Map)

	// iterate through each stage in the pipeline
	for _, _stage := range c.pipeline.Stages {
		if _stage.Name == constants.InitName {
			continue
		}

		// https://golang.org/doc/faq#closures_and_goroutines
		stage := _stage

		// create a new channel for each stage in the map
		stageMap.Store(stage.Name, make(chan error))

		// spawn errgroup routine for the stage
		//
		// https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Go
		stages.Go(func() error {
			c.Logger.Infof("planning %s stage", stage.Name)
			// plan the stage
			c.err = c.PlanStage(stageCtx, stage, stageMap)
			if c.err != nil {
				return fmt.Errorf("unable to plan stage: %w", c.err)
			}

			c.Logger.Infof("executing %s stage", stage.Name)
			// execute the stage
			c.err = c.ExecStage(stageCtx, stage, stageMap)
			if c.err != nil {
				return fmt.Errorf("unable to execute stage: %w", c.err)
			}

			return nil
		})
	}

	c.Logger.Debug("waiting for stages completion")
	// wait for the stages to complete or return an error
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Wait
	c.err = stages.Wait()
	if c.err != nil {
		return fmt.Errorf("unable to wait for stages: %w", c.err)
	}

	return c.err
}

// StreamBuild receives a StreamRequest and then
// runs StreamService or StreamStep in a goroutine.
func (c *client) StreamBuild(ctx context.Context) error {
	// cancel streaming after a timeout once the build has finished
	delayedCtx, cancelStreaming := context2.
		WithDelayedCancelPropagation(ctx, c.logStreamingTimeout, "streaming", c.Logger)
	defer cancelStreaming()

	// create an error group with the parent context
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup#WithContext
	streams, streamCtx := errgroup.WithContext(delayedCtx)

	defer func() {
		c.Logger.Trace("waiting for stream functions to return")

		err := streams.Wait()
		if err != nil {
			c.Logger.Errorf("error in a stream request, %v", err)
		}

		cancelStreaming()
		// wait for context to be done before reporting that everything has returned.
		<-delayedCtx.Done()
		// there might be one more log message from WithDelayedCancelPropagation
		// but there's not a good way to wait for that goroutine to finish.

		c.Logger.Info("all stream functions have returned")
	}()

	// allow the runtime to do log/event streaming setup at build-level
	streams.Go(func() error {
		// If needed, the runtime should handle synchronizing with
		// AssembleBuild which runs concurrently with StreamBuild.
		return c.Runtime.StreamBuild(streamCtx, c.pipeline)
	})

	for {
		select {
		case req, ok := <-c.streamRequests:
			if !ok {
				// ExecBuild is done requesting streams
				c.Logger.Debug("not accepting any more stream requests as channel is closed")
				return nil
			}

			streams.Go(func() error {
				// update engine logger with step metadata
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
				logger := c.Logger.WithField(req.Key, req.Container.Name)

				logger.Debugf("streaming %s container %s", req.Key, req.Container.ID)

				err := req.Stream(streamCtx, req.Container)
				if err != nil {
					logger.Error(err)
				}

				return nil
			})
		case <-delayedCtx.Done():
			c.Logger.Debug("not accepting any more stream requests as streaming context is canceled")
			// build done or canceled
			return nil
		}
	}
}

// loadLazySecrets is a helper function that injects secrets
// into the container right before execution, rather than
// during build planning. It is only available for the Docker runtime.
//

func loadLazySecrets(c *client, _step *pipeline.Container) error {
	_log := new(api.Log)

	lazySecrets := make(map[string]*api.Secret)
	lazyNoSubSecrets := make(map[string]*api.Secret)

	// this requires a small preface and brief description on
	// how normal secrets make it into a container:
	//
	// 1. pull secrets
	// 2. add them to the internal secrets map @ c.Secrets
	// 3. call escapeNewlineSecrets() on c.Secrets
	// 4. inject them into the container via injectSecrets()
	// 5. call container.Substitute()
	//
	// 1-3 happens in PlanBuild. 4 and 5 happens in
	// CreateStep and CreateService and for secrets added
	// via plugin.
	//
	// it's important to call out that container.Substitute()
	// can inadvertently(?) tweak the value of secrets,
	// particularly multiline secrets and/or secrets with
	// escaped newlines (for example). even worse, calling it
	// multiple times on the same container can tweak
	// them further. this is due to the json marshal/unmarshal
	// dance that happens during the substitution process.
	//
	// we can't move .Substitute() here because other aspects
	// of the build process depend on variables being
	// substituted earlier.
	//
	// so, to ensure lazy loaded secrets get the same
	// (mis)treatment and value (!) as regular secrets,
	// we will do the following here:
	//
	//  1. create a temporary map for lazy loaded secrets
	//  2. pull the lazy loaded secrets
	//  3. add them to temporary map
	//  4. call escapeNewlineSecrets() on temp map
	//  5. IF there are no lazy secrets, we stop here
	//  6. create a temporary copy of the step/container
	//  7. remove all existing environment variables except
	//     those needed for secret injection from the temp
	//     copy of the step/container
	//  8. inject the lazy loaded secrets into the
	//     temp step/container
	//  9. call .Substitute on the temp step/container
	//  10. move the lazy loaded secrets over to the
	//     actual step/container
	//
	// this will ensure the lazy loaded secrets return
	// the same value as they would as regular secrets
	// and also keep this process isolated to lazy secrets
	// create a temporary map akin to c.Secrets
	// ---- END ----

	// iterate through step secrets
	for _, s := range _step.Secrets {
		// iterate through each secret provided in the pipeline
		for _, secret := range c.pipeline.Secrets {
			// only lazy load non-plugin, step_start secrets
			if !secret.Origin.Empty() || !strings.EqualFold(s.Source, secret.Name) || strings.EqualFold(secret.Pull, constants.SecretPullBuild) {
				continue
			}

			// lazy loading not supported with Kubernetes, log info and continue
			if strings.EqualFold(constants.DriverKubernetes, c.Runtime.Driver()) {
				_log.AppendData([]byte(
					fmt.Sprintf("unable to pull secret %s: lazy loading secrets not available with Kubernetes runtime\n", s.Source)))

				_, err := c.Vela.Log.UpdateStep(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), _step.Number, _log)
				if err != nil {
					return err
				}

				continue
			}

			c.Logger.Infof("pulling secret %s", secret.Name)

			s, err := c.secret.pull(secret)
			if err != nil {
				c.err = err
				return fmt.Errorf("unable to pull secrets: %w", err)
			}

			_log.AppendData([]byte(
				fmt.Sprintf("$ vela view secret --secret.engine %s --secret.type %s --org %s --repo %s --name %s \n",
					secret.Engine, secret.Type, s.GetOrg(), s.GetRepo(), s.GetName())))

			sRaw, err := json.MarshalIndent(s.Sanitize(), "", " ")
			if err != nil {
				c.err = err
				return fmt.Errorf("unable to decode secret: %w", err)
			}

			_log.AppendData(append(sRaw, "\n"...))

			_, err = c.Vela.Log.UpdateStep(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), _step.Number, _log)
			if err != nil {
				return err
			}

			// add secret to the appropriate temp map
			if s.GetAllowSubstitution() {
				lazySecrets[secret.Name] = s
			} else {
				lazyNoSubSecrets[secret.Name] = s
			}
		}
	}

	// if we had lazy secrets, get them into the container
	if len(lazySecrets) > 0 {
		// create a copy of the current step/container
		tmpStep := new(pipeline.Container)
		*tmpStep = *_step

		c.Logger.Debug("clearing environment in temp step/container")
		// empty the environment
		tmpStep.Environment = map[string]string{}
		// but keep VELA_BUILD_EVENT as it's used in injectSecrets
		if _, ok := _step.Environment["VELA_BUILD_EVENT"]; ok {
			tmpStep.Environment["VELA_BUILD_EVENT"] = _step.Environment["VELA_BUILD_EVENT"]
		}

		c.Logger.Debug("escaping newlines in lazy loaded secrets")
		// escape newlines for secrets loaded on step_start
		escapeNewlineSecrets(lazySecrets)

		c.Logger.Debug("injecting lazy loaded secrets")
		// inject secrets for container
		err := injectSecrets(tmpStep, lazySecrets)
		if err != nil {
			return err
		}

		c.Logger.Debug("substituting container configuration after lazy loaded secret injection")
		// substitute container configuration
		err = tmpStep.Substitute()
		if err != nil {
			return err
		}

		c.Logger.Debug("injecting no-sub lazy loaded secrets")
		// inject secrets for container
		err = injectSecrets(tmpStep, lazyNoSubSecrets)
		if err != nil {
			return err
		}

		c.Logger.Debug("merge lazy loaded secrets into container")
		// merge lazy load secrets into original container
		err = _step.MergeEnv(tmpStep.Environment)
		if err != nil {
			return fmt.Errorf("failed to merge environment")
		}

		// clear out temporary var
		tmpStep = nil
	}

	return nil
}

// DestroyBuild cleans up the build after execution.
func (c *client) DestroyBuild(ctx context.Context) error {
	var err error

	defer func() {
		c.Logger.Info("deleting runtime build")
		// remove the runtime build for the pipeline
		err = c.Runtime.RemoveBuild(ctx, c.pipeline)
		if err != nil {
			c.Logger.Errorf("unable to remove runtime build: %v", err)
		}
	}()

	// destroy the steps for the pipeline
	for _, _step := range c.pipeline.Steps {
		if _step.Name == constants.InitName {
			continue
		}

		c.Logger.Infof("destroying %s step", _step.Name)
		// destroy the step
		err = c.DestroyStep(ctx, _step)
		if err != nil {
			c.Logger.Errorf("unable to destroy step: %v", err)
		}
	}

	// destroy the stages for the pipeline
	for _, _stage := range c.pipeline.Stages {
		if _stage.Name == constants.InitName {
			continue
		}

		c.Logger.Infof("destroying %s stage", _stage.Name)
		// destroy the stage
		err = c.DestroyStage(ctx, _stage)
		if err != nil {
			c.Logger.Errorf("unable to destroy stage: %v", err)
		}
	}

	// destroy the services for the pipeline
	for _, _service := range c.pipeline.Services {
		c.Logger.Infof("destroying %s service", _service.Name)
		// destroy the service
		err = c.DestroyService(ctx, _service)
		if err != nil {
			c.Logger.Errorf("unable to destroy service: %v", err)
		}
	}

	// destroy the secrets for the pipeline
	for _, _secret := range c.pipeline.Secrets {
		// skip over non-plugin secrets
		if _secret.Origin.Empty() {
			continue
		}

		c.Logger.Infof("destroying %s secret", _secret.Name)
		// destroy the secret
		err = c.secret.destroy(ctx, _secret.Origin)
		if err != nil {
			c.Logger.Errorf("unable to destroy secret: %v", err)
		}
	}

	// destroy output container
	err = c.outputs.destroy(ctx, c.OutputCtn)
	if err != nil {
		c.Logger.Errorf("unable to destroy output container: %v", err)
	}

	c.Logger.Info("deleting volume")
	// remove the runtime volume for the pipeline
	err = c.Runtime.RemoveVolume(ctx, c.pipeline)
	if err != nil {
		c.Logger.Errorf("unable to remove volume: %v", err)
	}

	c.Logger.Info("deleting network")
	// remove the runtime network for the pipeline
	err = c.Runtime.RemoveNetwork(ctx, c.pipeline)
	if err != nil {
		c.Logger.Errorf("unable to remove network: %v", err)
	}

	return err
}
