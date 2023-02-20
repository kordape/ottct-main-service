// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kordape/ottct-main-service/config"
	"github.com/kordape/ottct-main-service/internal/controller/http"
	"github.com/kordape/ottct-main-service/internal/database/postgres"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
	"github.com/kordape/ottct-main-service/pkg/logger"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	userManager, err := handler.NewAuthManager(db, log, validator.New(), cfg.SecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, log, userManager)
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
