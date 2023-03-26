package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kordape/ottct-main-service/internal/event"
)

// const (
// 	defaultTickInterval = 10 * time.Second
// )

type Worker struct {
	period    time.Duration // seconds
	quit      chan bool
	receiveFn event.ReceiveFakeNewsEventFn
}

func NewWorker(period int, receiveFn event.ReceiveFakeNewsEventFn) *Worker {
	return &Worker{
		period:    time.Duration(period),
		quit:      make(chan bool),
		receiveFn: receiveFn,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(w.period * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				w.receiveFn(context.Background())

			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.quit)
}

type SQSMessage struct {
	Message       string
	ReceiptHandle string
}

type ReceiveMessageFromQueue func(context context.Context) (*SQSMessage, error)

func ReceiveMessageFromQueueFn(sqsClient *sqs.SQS, queue string) ReceiveMessageFromQueue {
	return func(context context.Context) (*SQSMessage, error) {
		input := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queue),
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(5),
			// AttributeNames: []*string{
			// 	aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			// },
		}

		output, err := sqsClient.ReceiveMessageWithContext(context, input)
		if err != nil {
			return nil, fmt.Errorf("cannot receive sqs message: %w", err)
		}

		if len(output.Messages) == 0 {
			return nil, nil
		}

		return &SQSMessage{
			Message:       output.Messages[0].String(),
			ReceiptHandle: *output.Messages[0].ReceiptHandle,
		}, nil
	}
}

type DeleteMessageFromQueue func(context context.Context, receiptHandle string) error

func DeleteMessageFromQueueFn(sqsClient *sqs.SQS, queue string) DeleteMessageFromQueue {
	return func(context context.Context, receiptHandle string) error {

		input := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queue),
			ReceiptHandle: aws.String(receiptHandle),
		}

		_, err := sqsClient.DeleteMessageWithContext(context, input)
		if err != nil {
			return fmt.Errorf("cannot delete sqs message: %w", err)
		}

		return nil
	}
}
