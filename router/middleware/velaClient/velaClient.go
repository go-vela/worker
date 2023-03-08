package velaClient

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-vela/sdk-go/vela"
)

// Retrieve gets the build in the given context.
func Retrieve(c *gin.Context) (*vela.Client, error) {
	vc := c.Request.Context().Value("client").(*vela.Client)
	if vc == nil {
		return nil, fmt.Errorf("vc is nil %v", vc)
	}
	return vc, nil
}
