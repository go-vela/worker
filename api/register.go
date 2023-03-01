package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/router/middleware/token"
)

func Register(c *gin.Context) {
	// extract the deadloop channel that was packed into gin context
	v, ok := c.Get("auth-token")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no deadloop channel in the context")
		return
	}

	// make sure we configured it properly
	authChannel, ok := v.(chan string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "deadloop channel in the context is the wrong type")
		return
	}

	if len(authChannel) > 0 {
		c.JSON(http.StatusOK, "worker already registered")
		return
	}

	// this is a fake token, we would fetch this from the JSON body
	token, err := token.Retrieve(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no deadloop channel in the context")
		return
	}

	// send the token
	authChannel <- token

	// somehow we need to make sure the registration worked
	// maybe a second channel for registration results?
	c.JSON(http.StatusOK, "successfully registered the worker")
}
