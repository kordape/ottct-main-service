package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/internal/ses"
	"github.com/kordape/ottct-main-service/pkg/sqs"
	"github.com/sirupsen/logrus"
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

func (w *Worker) Run(log *logrus.Entry, subscriptionsManager *handler.SubscriptionManager) {
	ticker := time.NewTicker(w.period * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Debug("Worker tick")
				ctx := context.Background()
				messages, err := w.fakeNewsQueue.ReceiveMessages(ctx, sqs.WithVisibilityTimeout(20), sqs.WithMaxNumberOfMessages(5))
				if err != nil {
					log.WithError(err).Error("Error receiving message")
					continue
				}
				if len(messages) == 0 {
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
					log.WithError(err).Error("Error unmarshalling messages into proper structs")
				}
				eventsLen := len(events)

				log.WithField("events_len", eventsLen).Debug("Received fake news event(s)")

				for i, e := range events {
					log.WithFields(logrus.Fields{
						"index": i,
						"event": e,
					}).Debug("Started handling message")

					users, err := subscriptionsManager.GetSubscriptionsByEntity(e.EntityID, log)
					if err != nil {
						log.WithError(err).Error("Error getting subscribed users")
						continue
					}

					for _, user := range users {
						log.WithField("user", user.Id).Debug("Attempting to send email to user")

						// TODO how should we handle if sending email to one user fails?

						err = w.sendEmailFn(ctx, user, e.EntityID, e.TweetContent)
						if err != nil {
							log.WithError(err).Error("Error sending email to subscribed user")
							continue

						}

						log.Debug("Email sent to user")
					}

					err = w.fakeNewsQueue.DeleteMessage(ctx, e.ReceiptHandle)
					if err != nil {
						log.WithError(err).Error("error deleting message from queue")
						continue
					}

					log.Debug("Successfully deleted message")
				}

			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}
