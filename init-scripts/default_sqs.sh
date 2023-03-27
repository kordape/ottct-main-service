#!/usr/bin/env bash


awslocal sqs create-queue --queue-name default-fake-news

awslocal sqs send-message --queue-url http://localstack:4566/00000000000/default-fake-news --message-body '{"tweetContent":"test","entityId":"entity","tweetTimestamp":"2009-11-10T23:00:00Z"}'