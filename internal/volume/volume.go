// SPDX-License-Identifier: Apache-2.0

package volume

import (
	"fmt"
	"strings"
)

// Volume represents the volume definition used
// to create volumes for a container.
type Volume struct {
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	AccessMode  string `json:"access_mode,omitempty"`
}

// Parse digests the provided volume into a fully
// qualified volume reference. If an error
// occurs, it will return a nil volume.
func Parse(_volume string) *Volume {
	// parse the image provided into a fully qualified canonical reference
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/internal/image?tab=doc#ParseWithError
	v, err := ParseWithError(_volume)
	if err != nil {
		return nil
	}

	return v
}

// ParseWithError digests the provided volume into a
// fully qualified volume reference. If an error
// occurs, it will return a nil volume and the
// produced error.
func ParseWithError(_volume string) (*Volume, error) {
	// split each slice element into source, destination and access mode
	parts := strings.Split(_volume, ":")

	switch len(parts) {
	case 1:
		// return the read-only volume with the same source and destination
		return &Volume{
			Source:      parts[0],
			Destination: parts[0],
			AccessMode:  "ro",
		}, nil
	case 2:
		// return the read-only volume with different source and destination
		return &Volume{
			Source:      parts[0],
			Destination: parts[1],
			AccessMode:  "ro",
		}, nil
	case 3:
		// return the full volume with source, destination and access mode
		return &Volume{
			Source:      parts[0],
			Destination: parts[1],
			AccessMode:  parts[2],
		}, nil
	default:
		return nil, fmt.Errorf("volume %s requires at least 1, but no more than 2, `:`", _volume)
	}
}
