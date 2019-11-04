// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-vela/types"
)

// Pull pops an item from the specified channel off the queue.
func (c *client) Pull(channel string) (*types.Item, error) {
	// blocking list pop item from queue
	result, err := c.Queue.BLPop(0, channel).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to pop item from queue: %w", err)
	}

	item := new(types.Item)
	// unmarshal result into queue item
	err = json.Unmarshal([]byte(result[1]), item)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from queue: %w", err)
	}

	return item, nil
}
