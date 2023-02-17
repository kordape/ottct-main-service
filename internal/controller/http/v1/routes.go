package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

type routes struct {
	l logger.Interface
}

func NewRoutes(handler *gin.RouterGroup, l logger.Interface, userManager handler.UserManager) {
	r := &routes{l}

	echo := handler.Group("/echo")
	{
		echo.GET("/", r.echoHandler)
	}

	users := handler.Group("/users")
	{
		users.POST("/", r.newPostUsersHandler(userManager))
	}
}
