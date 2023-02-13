package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *routes) echoHandler(c *gin.Context) {
	r.l.Debug("Request received")
	c.JSON(http.StatusOK, "Echo response")
}
