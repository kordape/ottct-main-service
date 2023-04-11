// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kordape/ottct-main-service/config"
	"github.com/kordape/ottct-main-service/internal/controller/http"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
	"github.com/kordape/ottct-main-service/pkg/token"
)

// Run creates objects via constructors.
func Run(
	cfg *config.Config,
	log *logrus.Entry,
	userManager *handler.AuthManager,
	tokenManager *token.Manager,
	entityManager *handler.EntityManager,
	subscriptionsManager *handler.SubscriptionManager,
	twitterManager *handler.TwitterManager,
) {
	// HTTP Server
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	http.NewRouter(handler, log, userManager, tokenManager, entityManager, subscriptionsManager, twitterManager)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
