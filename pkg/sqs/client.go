package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// ReceiveOption is a functional option that can augment or modify a sqs.ReceiveMessageInput request.
type ReceiveOption func(*sqs.ReceiveMessageInput)

// WithVisibilityTimeout returns a ReceiveOption which setup the visibility timeout.
func WithVisibilityTimeout(s int32) ReceiveOption {
	return func(input *sqs.ReceiveMessageInput) {
		input.VisibilityTimeout = s
	}
}

// WithWaitTimeSeconds returns a ReceiveOption which setup the wait time.
func WithWaitTimeSeconds(s int32) ReceiveOption {
	return func(input *sqs.ReceiveMessageInput) {
		input.WaitTimeSeconds = s
	}
}

// Client represents a client that communicates with Amazon SQS about the request.
type Client interface {
	ReceiveMessage(ctx context.Context, options ...ReceiveOption) (*Message, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

type client struct {
	SQS *sqs.Client
	URL string
}

type Message struct {
	ReceiptHandle string
	Body          string
}

// NewClient returns a new SQS client.
func NewClient(sqsClient *sqs.Client, queueURL string) Client {
	return &client{
		SQS: sqsClient,
		URL: queueURL,
	}
}

func (c client) ReceiveMessage(ctx context.Context, options ...ReceiveOption) (*Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.URL),
		MaxNumberOfMessages: *aws.Int32(1),
	}

	for _, option := range options {
		if option != nil {
			option(input)
		}
	}

	output, err := c.SQS.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive the message from the queue: %w", err)
	}

	if len(output.Messages) == 0 {
		return nil, nil
	}

	return &Message{
		ReceiptHandle: *output.Messages[0].ReceiptHandle,
		Body:          *output.Messages[0].Body,
	}, nil
}

func (c client) DeleteMessage(ctx context.Context, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.URL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := c.SQS.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete the message from the queue: %w", err)
	}

	return nil
}
