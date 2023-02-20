package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/token"
)

type routes struct {
	l logger.Interface
}

func NewRoutes(handler *gin.RouterGroup, l logger.Interface, userManager handler.AuthManager, tokenManager *token.Manager) {
	r := &routes{l}

	authMiddleware := AuthMiddleware(tokenManager)

	echo := handler.Group("/echo")
	{
		echo.GET("/", r.echoHandler)
	}

	signup := handler.Group("/signup")
	{
		signup.POST("/", r.newSignUpHandler(userManager))
	}

	auth := handler.Group("/auth")
	{
		auth.POST("/", r.newAuthHandler(userManager))
	}

	secureEcho := handler.Group("/secureecho", authMiddleware)
	{
		secureEcho.GET("/", r.echoHandler)
	}
}
