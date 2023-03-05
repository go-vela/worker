package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/router/middleware/token"
	"github.com/sirupsen/logrus"
	"net/http"
)

// swagger:operation POST /api/v1/register-worker system register-worker
//
// # Perform a re-registration of the worker
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
//
// responses:
//
//	'501':
//	  description: Endpoint is not yet implemented
//	  schema:
//	    type: string
//
// RegisterWorker represents the API handler to register the worker
// by providing a registration token
func RegisterWorker(c *gin.Context) {

	// extract the auth token channel that was packed into gin context
	v, ok := c.Get("auth-token")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no auth token channel in the context")
		return
	}
	s, ok := c.Get("success")
	if !ok {
		logrus.Infof("s type is %T", s)
		c.JSON(http.StatusInternalServerError, "no success channel in the context")
		return
	}
	r, ok := c.Get("registered")
	if !ok {
		logrus.Infof("r type is %T", r)
		c.JSON(http.StatusInternalServerError, "no registered channel in the context")
		return
	}
	// make sure we configured the channel properly
	authChannel, ok := v.(chan string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "auth token channel in the context is the wrong type")
		return
	}
	// make sure we configured it properly
	successLoopChannel, ok := s.(chan bool)
	if !ok {
		c.JSON(http.StatusInternalServerError, "success channel in the context is the wrong type")
		return
	}
	// make sure we configured it properly
	registeredLoopChannel, ok := r.(chan bool)
	if !ok {
		c.JSON(http.StatusInternalServerError, "registered channel in the context is the wrong type")
		return
	}
	if len(registeredLoopChannel) > 0 {
		c.JSON(http.StatusOK, "worker is already registered")
		return
	}
	tkn, _ := token.Retrieve(c.Request)
	logrus.Infof("token %s", tkn)
	// write registration token to auth token channel
	authChannel <- tkn
	for v := range successLoopChannel {
		fmt.Println("received token from operate: ", v)
		if v == true {
			c.JSON(http.StatusOK, "worker has been registered")
			return
		}
		if v == false {
			c.JSON(http.StatusBadRequest, "Unable to register worker")
			return
		}

	}

}
