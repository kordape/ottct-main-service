package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kordape/tweety/pkg/logger"
)

type routes struct {
	l logger.Interface
}

func newRoutes(handler *gin.RouterGroup, l logger.Interface) {
	r := &routes{l}

	h := handler.Group("/users")
	{
		h.GET("/echo", r.echoHandler)
	}
}

func (r *routes) echoHandler(c *gin.Context) {
	r.l.Debug("Request received")
	c.JSON(http.StatusOK, "Echo response")
}
