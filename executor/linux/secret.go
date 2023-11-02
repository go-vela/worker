// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/internal/step"

	"github.com/sirupsen/logrus"
)

// secretSvc handles communication with secret processes during a build.
type secretSvc svc

var (
	// ErrUnrecognizedSecretType defines the error type when the
	// SecretType provided to the client is unsupported.
	ErrUnrecognizedSecretType = errors.New("unrecognized secret type")

	// ErrUnableToRetrieve defines the error type when the
	// secret is not able to be retrieved from the server.
	ErrUnableToRetrieve = errors.New("unable to retrieve secret")
)

// create configures the secret plugin for execution.
func (s *secretSvc) create(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with secret metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := s.client.Logger.WithField("secret", ctn.Name)

	ctn.Environment["VELA_DISTRIBUTION"] = s.client.build.GetDistribution()
	ctn.Environment["BUILD_HOST"] = s.client.build.GetHost()
	ctn.Environment["VELA_HOST"] = s.client.build.GetHost()
	ctn.Environment["VELA_RUNTIME"] = s.client.build.GetRuntime()
	ctn.Environment["VELA_VERSION"] = s.client.Version

	logger.Debug("setting up container")
	// setup the runtime container
	err := s.client.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("injecting secrets")
	// inject secrets for container
	err = injectSecrets(ctn, s.client.Secrets)
	if err != nil {
		return err
	}

	logger.Debug("substituting container configuration")
	// substitute container configuration
	err = ctn.Substitute()
	if err != nil {
		return fmt.Errorf("unable to substitute container configuration")
	}

	return nil
}

// destroy cleans up secret plugin after execution.
func (s *secretSvc) destroy(ctx context.Context, ctn *pipeline.Container) error {
	// update engine logger with secret metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := s.client.Logger.WithField("secret", ctn.Name)

	logger.Debug("inspecting container")
	// inspect the runtime container
	err := s.client.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("removing container")
	// remove the runtime container
	err = s.client.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// exec runs a secret plugins for a pipeline.
func (s *secretSvc) exec(ctx context.Context, p *pipeline.SecretSlice) error {
	// stream all the logs to the init step
	_init, err := step.Load(s.client.init, &s.client.steps)
	if err != nil {
		return err
	}

	defer func() {
		_init.SetFinished(time.Now().UTC().Unix())

		s.client.Logger.Infof("uploading %s step state", _init.GetName())
		// send API call to update the build
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.Update
		_, _, err = s.client.Vela.Step.Update(s.client.repo.GetOrg(), s.client.repo.GetName(), s.client.build.GetNumber(), _init)
		if err != nil {
			s.client.Logger.Errorf("unable to upload init state: %v", err)
		}
	}()

	// execute the secrets for the pipeline
	for _, _secret := range *p {
		// skip over non-plugin secrets
		if _secret.Origin.Empty() {
			continue
		}

		// update engine logger with secret metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
		logger := s.client.Logger.WithField("secret", _secret.Origin.Name)

		logger.Debug("running container")
		// run the runtime container
		err := s.client.Runtime.RunContainer(ctx, _secret.Origin, s.client.pipeline)
		if err != nil {
			return err
		}

		// trigger StreamStep goroutine with logging context
		s.client.streamRequests <- message.StreamRequest{
			Key:       "secret",
			Stream:    s.stream,
			Container: _secret.Origin,
		}

		logger.Debug("waiting for container")
		// wait for the runtime container
		err = s.client.Runtime.WaitContainer(ctx, _secret.Origin)
		if err != nil {
			return err
		}

		logger.Debug("inspecting container")
		// inspect the runtime container
		err = s.client.Runtime.InspectContainer(ctx, _secret.Origin)
		if err != nil {
			return err
		}

		// check the step exit code
		if _secret.Origin.ExitCode != 0 {
			// check if we ignore step failures
			if !_secret.Origin.Ruleset.Continue {
				// set build status to failure
				s.client.build.SetStatus(constants.StatusFailure)
			}

			// update the step fields
			_init.SetExitCode(_secret.Origin.ExitCode)
			_init.SetStatus(constants.StatusFailure)

			return fmt.Errorf("%s container exited with non-zero code", _secret.Origin.Name)
		}

		// send API call to update the build
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.Update
		_, _, err = s.client.Vela.Step.Update(s.client.repo.GetOrg(), s.client.repo.GetName(), s.client.build.GetNumber(), _init)
		if err != nil {
			s.client.Logger.Errorf("unable to upload init state: %v", err)
		}
	}

	return nil
}

// pull defines a function that pulls the secrets from the server for a given pipeline.
func (s *secretSvc) pull(secret *pipeline.Secret) (*library.Secret, error) {
	_secret := new(library.Secret)

	switch secret.Type {
	// handle repo secrets
	case constants.SecretOrg:
		org, key, err := secret.ParseOrg(s.client.repo.GetOrg())
		if err != nil {
			return nil, err
		}

		// send API call to capture the org secret
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SecretService.Get
		_secret, _, err = s.client.Vela.Secret.Get(secret.Engine, secret.Type, org, "*", key)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrUnableToRetrieve, err)
		}

		secret.Value = _secret.GetValue()

	// handle repo secrets
	case constants.SecretRepo:
		org, repo, key, err := secret.ParseRepo(s.client.repo.GetOrg(), s.client.repo.GetName())
		if err != nil {
			return nil, err
		}

		// send API call to capture the repo secret
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SecretService.Get
		_secret, _, err = s.client.Vela.Secret.Get(secret.Engine, secret.Type, org, repo, key)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrUnableToRetrieve, err)
		}

		secret.Value = _secret.GetValue()

	// handle shared secrets
	case constants.SecretShared:
		org, team, key, err := secret.ParseShared()
		if err != nil {
			return nil, err
		}

		// send API call to capture the repo secret
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SecretService.Get
		_secret, _, err = s.client.Vela.Secret.Get(secret.Engine, secret.Type, org, team, key)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrUnableToRetrieve, err)
		}

		secret.Value = _secret.GetValue()

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnrecognizedSecretType, secret.Type)
	}

	return _secret, nil
}

// stream tails the output for a secret plugin.
func (s *secretSvc) stream(ctx context.Context, ctn *pipeline.Container) error {
	// stream all the logs to the init step
	_log, err := step.LoadLogs(s.client.init, &s.client.stepLogs)
	if err != nil {
		return err
	}

	// update engine logger with secret metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithField
	logger := s.client.Logger.WithField("secret", ctn.Name)

	// create new buffer for uploading logs
	logs := new(bytes.Buffer)

	defer func() {
		// NOTE: Whenever the stream ends we want to ensure
		// that this function makes the call to update
		// the step logs
		logger.Trace(logs.String())

		// update the existing log with the last bytes
		//
		// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
		_log.AppendData(logs.Bytes())

		logger.Debug("uploading logs")
		// send API call to update the logs for the service
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#LogService.UpdateService
		_, err = s.client.Vela.Log.UpdateStep(s.client.repo.GetOrg(), s.client.repo.GetName(), s.client.build.GetNumber(), ctn.Number, _log)
		if err != nil {
			logger.Errorf("unable to upload container logs: %v", err)
		}
	}()

	logger.Debug("tailing container")
	// tail the runtime container
	rc, err := s.client.Runtime.TailContainer(ctx, ctn)
	if err != nil {
		return err
	}
	defer rc.Close()

	// create new scanner from the container output
	scanner := bufio.NewScanner(rc)

	// scan entire container output
	for scanner.Scan() {
		// write all the logs from the scanner
		logs.Write(append(scanner.Bytes(), []byte("\n")...))

		// if we have at least 1000 bytes in our buffer
		if logs.Len() > 1000 {
			logger.Trace(logs.String())

			// update the existing log with the new bytes
			//
			// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.AppendData
			_log.AppendData(logs.Bytes())

			logger.Debug("appending logs")
			// send API call to append the logs for the init step
			//
			// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#LogService.UpdateStep
			_, err = s.client.Vela.Log.UpdateStep(s.client.repo.GetOrg(), s.client.repo.GetName(), s.client.build.GetNumber(), s.client.init.Number, _log)
			if err != nil {
				return err
			}

			// flush the buffer of logs
			logs.Reset()
		}
	}

	logger.Info("finished streaming logs")

	return scanner.Err()
}

// TODO: Evaluate pulling this into a "bool" types function for injecting
//
// helper function to check secret whitelist before setting value.
func injectSecrets(ctn *pipeline.Container, m map[string]*library.Secret) error {
	// inject secrets for step
	for _, _secret := range ctn.Secrets {
		logrus.Tracef("looking up secret %s from pipeline secrets", _secret.Source)
		// lookup container secret in map
		s, ok := m[_secret.Source]
		if !ok {
			continue
		}

		logrus.Tracef("matching secret %s to container %s", _secret.Source, ctn.Name)
		// ensure the secret matches with the container
		if s.Match(ctn) {
			ctn.Environment[strings.ToUpper(_secret.Target)] = s.GetValue()
		}
	}

	return nil
}

// escapeNewlineSecrets is a helper function to double-escape escaped newlines,
// double-escaped newlines are resolved to newlines during env substitution.
func escapeNewlineSecrets(m map[string]*library.Secret) {
	for i, secret := range m {
		// only double-escape secrets that have been manually escaped
		if !strings.Contains(secret.GetValue(), "\\\\n") {
			s := strings.Replace(secret.GetValue(), "\\n", "\\\n", -1)
			m[i].Value = &s
		}
	}
}
