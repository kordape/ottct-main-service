// Package http implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/kordape/ottct-main-service/internal/controller/http/v1"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/token"
)

func NewRouter(handler *gin.Engine, l logger.Interface, userManager *handler.AuthManager, tokenManager *token.Manager, entityManager *handler.EntityManager, subscriptionsManager *handler.SubscriptionManager) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		v1.NewRoutes(h, l, userManager, tokenManager, entityManager, subscriptionsManager)
	}
}
