package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kordape/ottct-main-service/internal/event"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

// const (
// 	defaultTickInterval = 10 * time.Second
// )

type Worker struct {
	log                logger.Interface
	period             time.Duration // seconds
	quit               chan bool
	receiveMessageFn   event.ReceiveFakeNewsEventFn
	deleteMessageFn    event.DeleteEventFn
	sendNotificationFn event.SendNotificationFn
}

func NewWorker(log logger.Interface, period int, receiveFn event.ReceiveFakeNewsEventFn, deleteFn event.DeleteEventFn, sendNotificationFn event.SendNotificationFn) *Worker {
	return &Worker{
		log:                log,
		period:             time.Duration(period),
		quit:               make(chan bool),
		receiveMessageFn:   receiveFn,
		deleteMessageFn:    deleteFn,
		sendNotificationFn: sendNotificationFn,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(w.period * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				fakeNews, err := w.receiveMessageFn(ctx)
				if err != nil {
					w.log.Error(fmt.Sprintf("error receiving message: %s", err))
					continue
				}
				if fakeNews == nil {
					w.log.Debug("no new messages available")
					continue
				}

				w.log.Debug(fmt.Sprintf("received message: %s", fakeNews))

				// TODO get all users subscribed to entity

				err = w.sendNotificationFn(ctx, event.SendNotificationEvent{
					// TODO populate values
				})
				if err != nil {
					w.log.Error(fmt.Sprintf("error receiving message: %s", err))
				}

			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// func (w *Worker) Stop() {
// 	close(w.quit)
// }

func (w *Worker) sendNotification(fakeNews event.FakeNews) error {
	// TODO: implement
	// send email to sqs
	w.log.Debug("processing fake news event")
	return nil
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
