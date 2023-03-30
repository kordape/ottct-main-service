// Package app configures and runs application.
package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kordape/ottct-main-service/config"
	httplayer "github.com/kordape/ottct-main-service/internal/controller/http"
	"github.com/kordape/ottct-main-service/internal/database/postgres"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/token"
	"github.com/kordape/ottct-poller-service/pkg/predictor"
	"github.com/kordape/ottct-poller-service/pkg/twitter"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	dbClient, err := gorm.Open(pg.Open(cfg.DB.URL), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.New(dbClient, log)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	tokenManager, err := token.NewManager(cfg.SecretKey, "ottct")
	if err != nil {
		log.Fatal(err)
	}

	userManager, err := handler.NewAuthManager(db, log, validator.New(), tokenManager)
	if err != nil {
		log.Fatal(err)
	}

	entityManager := handler.NewEntityManager(db, log)

	subscriptionsManager := handler.NewSubscriptionManager(db, log)

	twitterManager, err := handler.NewTwitterManager(
		log,
		validator.New(),
		twitter.New(
			&http.Client{
				Timeout: 10 * time.Second,
			},
			cfg.TwitterBearerKey,
		),
		predictor.New(
			&http.Client{
				Timeout: 10 * time.Second,
			},
			cfg.PredictorURL,
		),
	)

	// HTTP Server
	handler := gin.New()
	httplayer.NewRouter(handler, log, userManager, tokenManager, entityManager, subscriptionsManager, twitterManager)
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
