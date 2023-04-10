package main

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sqsservice "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kordape/ottct-main-service/config"
	"github.com/kordape/ottct-main-service/internal/app"
	"github.com/kordape/ottct-main-service/internal/database/postgres"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/internal/ses"
	"github.com/kordape/ottct-main-service/internal/worker"
	sqspkg "github.com/kordape/ottct-main-service/pkg/sqs"
	"github.com/kordape/ottct-main-service/pkg/token"
	"github.com/kordape/ottct-poller-service/pkg/predictor"
	"github.com/kordape/ottct-poller-service/pkg/twitter"
)

func main() {

	log := logrus.StandardLogger()
	logrus.SetReportCaller(true)
	logrus.SetFormatter(
		&logrus.TextFormatter{
			ForceColors: true,
		},
	)

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(level)
	}

	awsCfg, err := initAWSConfig(cfg.AWS.Region, cfg.AWS.EndpointURL)
	if err != nil {
		log.Fatal(err)
	}

	sqsClient := sqspkg.NewClient(sqsservice.NewFromConfig(awsCfg), cfg.AWS.FakeNewsQueueUrl)

	dbClient, err := gorm.Open(pg.Open(cfg.DB.URL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.New(dbClient, logrus.NewEntry(log))
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

	userManager, err := handler.NewAuthManager(db, validator.New(), tokenManager)
	if err != nil {
		log.Fatal(err)
	}

	entityManager := handler.NewEntityManager(db)

	subscriptionsManager, err := handler.NewSubscriptionManager(db, validator.New())
	if err != nil {
		log.Fatal(err)
	}

	w := worker.NewWorker(
		cfg.App.PollerInterval,
		sqsClient,
		ses.SendFakeNewsEmailFnBuilder(sesv2.NewFromConfig(awsCfg), cfg.AWS.VerifiedSender),
	)

	// Run sqs poller worker (as a background process)
	w.Run(logrus.NewEntry(log), subscriptionsManager)

	twitterManager, err := handler.NewTwitterManager(
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
	app.Run(cfg, logrus.NewEntry(log), userManager, tokenManager, entityManager, subscriptionsManager, twitterManager)
}

func initAWSConfig(region, endpoint string) (aws.Config, error) {
	if len(endpoint) > 0 {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == sqsservice.ServiceID || service == sesv2.ServiceID {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}
			// Returning EndpointNotFoundError will allow the service to fallback
			// to it's default resolution.
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		return awsconfig.LoadDefaultConfig(
			context.Background(),
			awsconfig.WithEndpointResolverWithOptions(customResolver),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("id", "fake-secret", "fake-token")),
			awsconfig.WithRegion(region),
		)
	}

	return awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(region))
}
