package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	sqsservice "github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/kordape/ottct-main-service/config"
	"github.com/kordape/ottct-main-service/internal/app"
	"github.com/kordape/ottct-main-service/internal/sns"
	"github.com/kordape/ottct-main-service/internal/sqs"
	"github.com/kordape/ottct-main-service/internal/worker"
	"github.com/kordape/ottct-main-service/pkg/logger"
	sqspkg "github.com/kordape/ottct-main-service/pkg/sqs"
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

	sqsClient := sqspkg.NewClient(sqsservice.NewFromConfig(awsCfg), fmt.Sprintf("%s/000000000000/%s", cfg.AWS.EndpointUrl, cfg.AWS.FakeNewsQueueName))

	w := worker.NewWorker(
		log,
		cfg.App.PollerInterval,
		sqs.ReceiveFakeNewsEventsFnBuilder(sqsClient, log),
		sqs.DeleteMessageFnBuilder(sqsClient, log),
		sns.SendNotificationEventFnBuilder(),
	)

	// Run sqs poller worker (as a background process)
	w.Run()

	// Run app
	app.Run(cfg)
}
