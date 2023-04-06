package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kordape/ottct-main-service/internal/sns"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

type Worker struct {
	log                logger.Interface
	period             time.Duration // seconds
	quit               chan bool
	fakeNewsQueue      sqs.Client
	sendNotificationFn sns.SendNotificationEventFn
}

func NewWorker(log logger.Interface, period int, fakeNewsQueue sqs.Client, sendNotificationFn sns.SendNotificationEventFn) *Worker {
	return &Worker{
		log:           log,
		period:        time.Duration(period),
		quit:          make(chan bool),
		sendNotificationFn: sendNotificationFn,
		fakeNewsQueue: fakeNewsQueue,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(w.period * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				messages, err := w.fakeNewsQueue.ReceiveMessages(ctx, sqs.WithVisibilityTimeout(20), sqs.WithMaxNumberOfMessages(5))
				if err != nil {
					w.log.Error(fmt.Sprintf("Error receiving message: %s", err))
					continue
				}
				if len(messages) == 0 {
					w.log.Debug("No new messages available...")
					continue
				}

				events, err := func() ([]sqs.FakeNewsEvent, error) {
					events := make([]sqs.FakeNewsEvent, len(messages))
					for i, msg := range messages {
						var event sqs.FakeNewsEvent
						err = json.Unmarshal([]byte(msg.Body), &event)
						if err != nil {
							return nil, fmt.Errorf("Error unmarshalling messages into FakeNewsEvents: %s", err)
						}

						event.ReceiptHandle = msg.ReceiptHandle

						events[i] = event
					}

					return events, nil
				}()
				if err != nil {
					w.log.Error(fmt.Sprintf("Error unmarshalling messages into proper structs: %s", err))
					continue
				}

				w.log.Debug("Received fake news events.")

				// TODO get all users subscribed to entity

				for _, e := range events {
					err = w.sendNotificationFn(ctx, sns.SendNotificationEvent{
						// TODO populate values
					})
					if err != nil {
						w.log.Error(fmt.Sprintf("Error sending notification: %s", err))
						continue
					}
					w.log.Debug("Notifications about the fake news event sent!")

					err = w.fakeNewsQueue.DeleteMessage(ctx, e.ReceiptHandle)
					if err != nil {
						w.log.Error(fmt.Sprintf("Error deleting message from queue: %s", err))
						continue
					}

					w.log.Debug("Message successfully deleted from queue!")
				}

			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) sendNotification(context.Context, sns.SendNotificationEvent) error {
	// TODO: implement
	// send email to subscribed user
	return nil
}
