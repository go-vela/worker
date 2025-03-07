// SPDX-License-Identifier: Apache-2.0

package linux

import "github.com/go-vela/server/constants"

// Driver outputs the configured executor driver.
func (c *client) Driver() string {
	return constants.DriverLinux
}
