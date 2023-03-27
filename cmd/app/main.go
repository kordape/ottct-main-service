package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	sqsservice "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kordape/ottct-main-service/config"
	"github.com/kordape/ottct-main-service/internal/app"
	"github.com/kordape/ottct-main-service/internal/event"
	"github.com/kordape/ottct-main-service/internal/worker"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

const (
	defaultTickInterval = 5
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
			URL:               "http://localstack:4566",
			SigningRegion:     os.Getenv("AWS_REGION"),
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

	sqsClient := sqs.NewClient(sqsservice.NewFromConfig(awsCfg), os.Getenv("QUEUE_URL"))

	w := worker.NewWorker(
		log,
		defaultTickInterval,
		event.ReceiveFakeNewsEventFnBuilder(sqsClient, log),
		event.DeleteEventFnBuilder(sqsClient, log),
		event.SendNotificationFnBuilder(),
	)

	// TODO context?
	w.Run()

	// Run
	app.Run(cfg)
}
