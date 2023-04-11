// Package http implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	v1 "github.com/kordape/ottct-main-service/internal/controller/http/v1"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/token"
)

func NewRouter(
	handler *gin.Engine,
	l *logrus.Entry,
	userManager *handler.AuthManager,
	tokenManager *token.Manager,
	entityManager *handler.EntityManager,
	subscriptionsManager *handler.SubscriptionManager,
	twitterManager *handler.TwitterManager,
) {
	// Options
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		v1.NewRoutes(
			h,
			l,
			userManager,
			tokenManager,
			entityManager,
			subscriptionsManager,
			twitterManager,
		)
	}
}
