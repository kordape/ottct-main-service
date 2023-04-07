package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sqsservice "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-playground/validator/v10"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kordape/ottct-main-service/config"
	"github.com/kordape/ottct-main-service/internal/app"
	"github.com/kordape/ottct-main-service/internal/database/postgres"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/internal/ses"
	"github.com/kordape/ottct-main-service/internal/worker"
	"github.com/kordape/ottct-main-service/pkg/logger"
	sqspkg "github.com/kordape/ottct-main-service/pkg/sqs"
	"github.com/kordape/ottct-main-service/pkg/token"
	"github.com/kordape/ottct-poller-service/pkg/predictor"
	"github.com/kordape/ottct-poller-service/pkg/twitter"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log := logger.New(cfg.Log.Level)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               cfg.AWS.EndpointUrl,
			SigningRegion:     cfg.AWS.Region,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		panic("configuration error: " + err.Error())
	}

	sqsClient := sqspkg.NewClient(sqsservice.NewFromConfig(awsCfg), cfg.AWS.FakeNewsQueueUrl)

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

	subscriptionsManager, err := handler.NewSubscriptionManager(db, log, validator.New())
	if err != nil {
		log.Fatal(err)
	}

	w := worker.NewWorker(
		cfg.App.PollerInterval,
		sqsClient,
		ses.SendFakeNewsEmailFnBuilder(sesv2.NewFromConfig(awsCfg), cfg.AWS.VerifiedSender),
	)

	// Run sqs poller worker (as a background process)
	w.Run(log, subscriptionsManager)

	twitterManager, err := handler.NewTwitterManager(
		log,
		validator.New(),
		twitter.New(
			&http.Client{
				Timeout: 5 * time.Second,
			},
			cfg.TwitterBearerKey,
		),
		predictor.New(
			&http.Client{
				Timeout: 5 * time.Second,
			},
			cfg.PredictorURL,
		),
	)

	// Run app
	app.Run(cfg, log, userManager, tokenManager, entityManager, subscriptionsManager, twitterManager)
}
