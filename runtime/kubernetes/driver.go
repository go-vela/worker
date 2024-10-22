// SPDX-License-Identifier: Apache-2.0

package kubernetes

import "github.com/go-vela/server/constants"

// Driver outputs the configured runtime driver.
func (c *client) Driver() string {
	return constants.DriverKubernetes
}
