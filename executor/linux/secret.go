// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// PullSecret defines a function that pulls the secrets for a given pipeline.
func (c *client) PullSecret(ctx context.Context) error {
	var err error
	p := c.pipeline

	secrets := make(map[string]*library.Secret)
	sec := new(library.Secret)

	// iterate through each secret provided in the pipeline
	for _, s := range p.Secrets {
		// if the secret isn't a native or vault type
		if !strings.EqualFold(s.Engine, constants.DriverNative) &&
			!strings.EqualFold(s.Engine, constants.DriverVault) {
			return fmt.Errorf("unrecognized secret engine: %s", s.Engine)
		}

		switch s.Type {
		// handle org secrets
		case constants.SecretOrg:
			c.logger.Debug("pulling org secret")
			// get org secret
			sec, err = c.getOrg(s)
			if err != nil {
				return err
			}
		// handle repo secrets
		case constants.SecretRepo:
			c.logger.Debug("pulling repo secret")
			// get repo secret
			sec, err = c.getRepo(s)
			if err != nil {
				return err
			}
		// handle shared secrets
		case constants.SecretShared:
			c.logger.Debug("pulling shared secret")
			// get shared secret
			sec, err = c.getShared(s)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognized secret type: %s", s.Type)
		}

		// add secret to the map
		secrets[s.Name] = sec
	}

	// overwrite the engine secret map
	c.Secrets = secrets
	return nil
}

// getOrg is a helper function to parse and capture
// the org secret from the provided secret engine.
func (c *client) getOrg(s *pipeline.Secret) (*library.Secret, error) {
	c.logger.Tracef("pulling %s %s secret %s", s.Engine, s.Type, s.Name)

	// variables necessary for secret
	org := c.repo.GetOrg()
	repo := "*"
	path := s.Key

	// check if the full path was provided
	if strings.Contains(path, "/") {
		// split the full path into parts
		parts := strings.SplitN(path, "/", 2)

		// secret is invalid
		if len(parts) != 2 {
			return nil, fmt.Errorf("path %s for %s secret %s is invalid", s.Key, s.Type, s.Name)
		}

		// check if the org provided matches what we expect
		if strings.EqualFold(parts[0], org) {
			// update the variables
			org = parts[0]
			path = parts[1]
		}
	}

	// send API call to capture the org secret
	secret, _, err := c.Vela.Secret.Get(s.Engine, s.Type, org, repo, path)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve %s secret %s: %w", s.Type, s.Key, err)
	}

	// overwrite the secret value
	s.Value = secret.GetValue()
	return secret, nil
}

// getRepo is a helper function to parse and capture
// the repo secret from the provided secret engine.
func (c *client) getRepo(s *pipeline.Secret) (*library.Secret, error) {
	c.logger.Tracef("pulling %s %s secret %s", s.Engine, s.Type, s.Name)

	// variables necessary for secret
	org := c.repo.GetOrg()
	repo := c.repo.GetName()
	path := s.Key

	// check if the full path was provided
	if strings.Contains(path, "/") {
		// split the full path into parts
		parts := strings.SplitN(path, "/", 3)

		// secret is invalid
		if len(parts) != 3 {
			return nil, fmt.Errorf("path %s for %s secret %s is invalid", s.Key, s.Type, s.Name)
		}

		// check if the org provided matches what we expect
		if strings.EqualFold(parts[0], org) {
			// update the org variable
			org = parts[0]

			// check if the repo provided matches what we expect
			if strings.EqualFold(parts[1], repo) {
				// update the variables
				repo = parts[1]
				path = parts[2]
			}
		}
	}

	// send API call to capture the repo secret
	secret, _, err := c.Vela.Secret.Get(s.Engine, s.Type, org, repo, path)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve %s secret %s: %w", s.Type, s.Key, err)
	}

	// overwrite the secret value
	s.Value = secret.GetValue()
	return secret, nil
}

// getShared is a helper function to parse and capture
// the shared secret from the provided secret engine.
func (c *client) getShared(s *pipeline.Secret) (*library.Secret, error) {
	c.logger.Tracef("pulling %s %s secret %s", s.Engine, s.Type, s.Name)

	// variables necessary for secret
	org := c.repo.GetOrg()
	var team string
	path := s.Key

	// check if the full path was provided
	if strings.Contains(path, "/") {
		// split the full path into parts
		parts := strings.SplitN(path, "/", 3)

		// secret is invalid
		if len(parts) != 3 {
			return nil, fmt.Errorf("Path %s for %s secret %s is invalid", s.Key, s.Type, s.Name)
		}

		// check if the org provided matches what we expect
		if strings.EqualFold(parts[0], org) {
			// update the org variable
			org = parts[0]

			// check if the team provided matches what we expect
			if strings.EqualFold(parts[1], team) {
				// update the variables
				team = parts[1]
				path = parts[2]
			}
		}
	}

	// send API call to capture the shared secret
	secret, _, err := c.Vela.Secret.Get(s.Engine, s.Type, org, team, path)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve %s secret %s: %w", s.Type, s.Key, err)
	}

	// overwrite the secret value
	s.Value = secret.GetValue()
	return secret, nil
}

// helper function to check secret whitelist before setting value
// TODO: Evaluate pulling this into a "bool" types function for injecting
func injectSecrets(s *pipeline.Container, m map[string]*library.Secret) error {

	// inject secrets for step
	for _, secret := range s.Secrets {
		eACL, iACL := false, false
		events, images := []string{}, []string{}

		if value, ok := m[secret.Source]; ok {
			events = value.GetEvents()
			images = value.GetImages()
		}

		// check event whitelist for if statements
		if !s.Ruleset.If.Empty() {
			for _, e := range events {
				if s.Ruleset.If.Match(&pipeline.RuleData{Event: e}, "and") {
					eACL = true
					break
				}
			}
		}

		// check event whitelist for unless statements
		if !s.Ruleset.Unless.Empty() {
			for _, e := range events {
				if s.Ruleset.Unless.Match(&pipeline.RuleData{Event: e}, "and") {
					eACL = true
					break
				}
			}
		}

		// check images whitelist
		for _, i := range images {
			if strings.HasPrefix(s.Image, i) && (len(i) != 0) {
				iACL = true
				break
			}
		}

		// inject secrets into environment
		switch {
		case iACL && (len(events) == 0):
			s.Environment[strings.ToUpper(secret.Target)] = *m[secret.Source].Value
		case eACL && (len(images) == 0):
			s.Environment[strings.ToUpper(secret.Target)] = *m[secret.Source].Value
		case eACL && iACL:
			s.Environment[strings.ToUpper(secret.Target)] = *m[secret.Source].Value
		}
	}

	return nil
}
