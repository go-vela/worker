// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package image

import (
	"github.com/docker/distribution/reference"
)

// Parse digests the provided image into a fully
// qualified canonical reference. If an error
// occurs, it will return the provided image.
func Parse(_image string) string {
	// parse the image provided into a fully qualified canonical reference
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/internal/image?tab=doc#ParseWithError
	_canonical, err := ParseWithError(_image)
	if err != nil {
		return _image
	}

	return _canonical
}

// ParseWithError digests the provided image into a
// fully qualified canonical reference. If an error
// occurs, it will return the last digested form of
// the image.
func ParseWithError(_image string) (string, error) {
	// parse the image provided into a
	// named, fully qualified reference
	//
	// https://pkg.go.dev/github.com/docker/distribution/reference?tab=doc#ParseAnyReference
	_reference, err := reference.ParseAnyReference(_image)
	if err != nil {
		return _image, err
	}

	// ensure we have the canonical form of the named reference
	//
	// https://pkg.go.dev/github.com/docker/distribution/reference?tab=doc#ParseNamed
	_canonical, err := reference.ParseNamed(_reference.String())
	if err != nil {
		return _reference.String(), err
	}

	// ensure the canonical reference has a tag
	//
	// https://pkg.go.dev/github.com/docker/distribution/reference?tab=doc#TagNameOnly
	return reference.TagNameOnly(_canonical).String(), nil
}

// IsPrivilegedImage digests the provided image with a
// privileged pattern to see if the image meets the criteria
// needed to allow a Docker Socket mount.
func IsPrivilegedImage(image, privileged string) (bool, error) {
	// parse the image provided into a
	// named, fully qualified reference
	//
	// https://pkg.go.dev/github.com/docker/distribution/reference?tab=doc#ParseAnyReference
	_refImg, err := reference.ParseAnyReference(image)
	if err != nil {
		return false, err
	}

	// ensure we have the canonical form of the named reference
	//
	// https://pkg.go.dev/github.com/docker/distribution/reference?tab=doc#ParseNamed
	_canonical, err := reference.ParseNamed(_refImg.String())
	if err != nil {
		return false, err
	}

	// add default tag "latest" when tag does not exist
	_refImg = reference.TagNameOnly(_canonical)

	// check if the image matches the privileged pattern
	//
	// https://pkg.go.dev/github.com/docker/distribution/reference#FamiliarMatch
	match, err := reference.FamiliarMatch(privileged, _refImg)
	if err != nil {
		return false, err
	}

	return match, nil
}
