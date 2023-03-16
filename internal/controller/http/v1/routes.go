package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/kordape/ottct-main-service/internal/database/postgres"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/token"
)

type routes struct {
	l  logger.Interface
	db *postgres.DB
}

func NewRoutes(handler *gin.RouterGroup, l logger.Interface, userManager *handler.AuthManager, tokenManager *token.Manager, entityManager *handler.EntityManager, subscriptionsManager *handler.SubscriptionManager) {
	r := &routes{l: l}

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

	entities := handler.Group("/entities", authMiddleware)
	{
		entities.GET("/", r.getEntitiesHandler(entityManager))
	}

	subscriptions := handler.Group("/subscribe", authMiddleware)
	{
		subscriptions.GET("/", r.getSubscriptionsHandler(subscriptionsManager, tokenManager))
		subscriptions.POST("/:entityid", r.updateSubscriptionsHandler(entityManager, subscriptionsManager, tokenManager))
	}
}
