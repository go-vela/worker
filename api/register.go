// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/go-vela/worker/router/middleware/token"
)

// swagger:operation POST /register system Register
//
// Fill registration token channel in worker to continue operation
//
// ---
// produces:
// - application/json
// parameters:
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully passed token to worker
//     schema:
//       type: string
//   '401':
//     description: No token was passed
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to pass token to worker
//     schema:
//       "$ref": "#/definitions/Error"

// Register will pass the token given in the request header to the register token
// channel of the worker. This will unblock operation if the worker has not been
// registered and the provided registration token is valid.
func Register(c *gin.Context) {
	// extract the worker hostname that was packed into gin context
	w, ok := c.Get("worker-hostname")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no worker hostname in the context")
		return
	}

	// extract the register token channel that was packed into gin context
	v, ok := c.Get("register-token")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no register token channel in the context")
		return
	}

	// make sure we configured the channel properly
	rChan, ok := v.(chan string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "register token channel in the context is the wrong type")
		return
	}

	// if token is present in the channel, deny registration
	// this will likely never happen as the channel is offloaded immediately
	if len(rChan) > 0 {
		c.JSON(http.StatusOK, "worker already registered")
		return
	}

	// retrieve auth token from header
	token, err := token.Retrieve(c.Request)
	if err != nil {
		// an error occurs when no token was passed
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// extract the subject from the token
	sub, err := getSubjectFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// make sure we configured the hostname properly
	hostname, ok := w.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "worker hostname in the context is the wrong type")
		return
	}

	// if the subject doesn't match the worker hostname return an error
	if sub != hostname {
		c.JSON(http.StatusUnauthorized, "worker hostname is invalid")
		return
	}

	// write registration token to auth token channel
	rChan <- token

	c.JSON(http.StatusOK, "successfully passed token to worker")
}

// getSubjectFromToken is a helper function to extract
// the subject from the token claims.
func getSubjectFromToken(token string) (string, error) {
	// create a new JWT parser
	j := jwt.NewParser()

	// parse the payload
	t, _, err := j.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("unable to parse token")
	}

	sub, err := t.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("unable to get subject from token")
	}

	// make sure there was a subject defined
	if len(sub) == 0 {
		return "", fmt.Errorf("no subject defined in token")
	}

	return sub, nil
}
