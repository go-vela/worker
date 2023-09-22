// SPDX-License-Identifier: Apache-2.0

package kubernetes

import "github.com/go-vela/types/constants"

// Driver outputs the configured runtime driver.
func (c *client) Driver() string {
	return constants.DriverKubernetes
}
