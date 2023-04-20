// Package http implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	v0 "github.com/kordape/ottct-main-service/internal/controller/http/v0"
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
	hv0 := handler.Group("/v0")
	{
		v0.NewRoutes(
			hv0,
			l,
			userManager,
			tokenManager,
			entityManager,
			subscriptionsManager,
			twitterManager,
		)
	}

	hv1 := handler.Group("/v1")
	{
		v1.NewRoutes(
			hv1,
			l,
			userManager,
			tokenManager,
			entityManager,
			subscriptionsManager,
			twitterManager,
		)
	}
}
