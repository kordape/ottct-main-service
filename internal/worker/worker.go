package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/internal/ses"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

type Worker struct {
	period        time.Duration // seconds
	quit          chan bool
	fakeNewsQueue sqs.Client
	sendEmailFn   ses.SendFakeNewsEmailFn
}

func NewWorker(period int, fakeNewsQueue sqs.Client, sendEmailFn ses.SendFakeNewsEmailFn) *Worker {
	return &Worker{
		period:        time.Duration(period),
		quit:          make(chan bool),
		fakeNewsQueue: fakeNewsQueue,
		sendEmailFn:   sendEmailFn,
	}
}

func (w *Worker) Run(log logger.Interface, subscriptionsManager *handler.SubscriptionManager) {
	ticker := time.NewTicker(w.period * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				messages, err := w.fakeNewsQueue.ReceiveMessages(ctx, sqs.WithVisibilityTimeout(20), sqs.WithMaxNumberOfMessages(5))
				if err != nil {
					log.Error(fmt.Sprintf("Error receiving message: %s", err))
					continue
				}
				if len(messages) == 0 {
					log.Debug("No new messages available...")
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
					log.Error(fmt.Sprintf("Error unmarshalling messages into proper structs: %s", err))
				}
				eventsLen := len(events)

				log.Debug(fmt.Sprintf("Received %d fake news event(s).", eventsLen))

				for i, e := range events {
					log.Debug("Started handling message %d out of %d.", i+1, eventsLen)

					users, err := subscriptionsManager.GetSubscriptionsByEntity(e.EntityID)
					if err != nil {
						log.Error(fmt.Sprintf("Error sending getting subscribed users: %s", err))
						continue
					}

					
					for _, user := range users {
						log.Debug(fmt.Sprintf("Attempting to send email to user with id: %d.", user.Id))

						// TODO how should we handle if sending email to one user fails?

						err = w.sendEmailFn(ctx, user, e.EntityID, e.TweetContent)
						if err != nil {
							log.Error(fmt.Sprintf("Error sending email to subscribed user with id %d: %s", user.Id, err))
							continue

						}

						log.Debug(fmt.Sprintf("Email sent to user with id: %d.", user.Id))
					}

					err = w.fakeNewsQueue.DeleteMessage(ctx, e.ReceiptHandle)
					if err != nil {
						log.Error(fmt.Sprintf("error deleting message from queue: %s", err))
						continue
					}

					log.Debug(fmt.Sprintf("Successfully deleted message %d out of %d from queue.", i+1, eventsLen))
				}

			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}
