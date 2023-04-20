package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kordape/ottct-main-service/internal/database/postgres"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
	"github.com/kordape/ottct-main-service/pkg/token"
)

type routes struct {
	db *postgres.DB
}

func NewRoutes(
	handler *gin.RouterGroup,
	l *logrus.Entry,
	userManager *handler.AuthManager,
	tokenManager *token.Manager,
	entityManager *handler.EntityManager,
	subscriptionsManager *handler.SubscriptionManager,
	twitterManager *handler.TwitterManager,
) {
	r := &routes{}

	log := l.WithField("domain", "api-v0")
	handler.Use(httpserver.Logging(log))

	echo := handler.Group("/echo")
	{
		echo.GET("/", r.echoHandler)
	}

	entities := handler.Group("/entities")
	{
		entities.GET("/", r.getEntitiesHandler(entityManager))
	}

	tweets := handler.Group("/tweets")
	{
		tweets.GET("/", r.newGetAnalyticsHandler(twitterManager))
	}

	subscribe := handler.Group("/subscribe")
	{
		subscribe.POST("/", r.newPostSubscribeHandler(entityManager, userManager, subscriptionsManager))
	}
}
