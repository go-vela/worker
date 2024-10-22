// SPDX-License-Identifier: Apache-2.0

package local

import "github.com/go-vela/server/constants"

// Driver outputs the configured executor driver.
func (c *client) Driver() string {
	return constants.DriverLocal
}
