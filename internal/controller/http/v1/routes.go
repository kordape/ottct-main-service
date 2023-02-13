package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/kordape/tweety/pkg/logger"
)

type routes struct {
	l logger.Interface
}

func NewRoutes(handler *gin.RouterGroup, l logger.Interface) {
	r := &routes{l}

	h := handler.Group("/users")
	{
		h.GET("/echo", r.echoHandler)
	}
}
