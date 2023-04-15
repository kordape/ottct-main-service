package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
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
		ctx := context.Background()

		for {
			select {
			case <-ticker.C:
				log.Debug("Worker tick")
				err := w.work(ctx, log, subscriptionsManager)
				if err != nil {
					log.Errorf("Error while running worker: %s", err)
				}
			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) work(ctx context.Context, log *logrus.Entry, subscriptionsManager *handler.SubscriptionManager) (err error) {
	defer recoverWithLog(&err)

	messages, err := w.fakeNewsQueue.ReceiveMessages(ctx, sqs.WithVisibilityTimeout(20), sqs.WithMaxNumberOfMessages(5))
	if err != nil {
		return fmt.Errorf("error receiving message from queue: %w", err)
	}
	if len(messages) == 0 {
		return nil
	}

	// convert messages to events
	events, err := func() ([]sqs.FakeNewsEvent, error) {
		events := make([]sqs.FakeNewsEvent, len(messages))
		for i, msg := range messages {
			var event sqs.FakeNewsEvent
			err = json.Unmarshal([]byte(msg.Body), &event)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling messages into FakeNewsEvents: %s", err)
			}

			event.ReceiptHandle = msg.ReceiptHandle

			events[i] = event
		}

		return events, nil
	}()
	if err != nil {
		return fmt.Errorf("error unmarshalling messages into proper structs: %w", err)
	}

	eventsLen := len(events)

	log.WithField("events_len", eventsLen).Debug("Received fake news event(s)")

	// handle events
	for i, e := range events {
		log.WithFields(logrus.Fields{
			"index": i,
			"event": e,
		}).Debug("Started handling message")

		users, err := subscriptionsManager.GetSubscriptionsByEntity(e.EntityID, log)
		if err != nil {
			return fmt.Errorf("error getting subscribed users: %w", err)
		}

		failedUsers := make([]uint, 0)

		for _, user := range users {
			log.WithField("user", user.Id).Debug("Attempting to send email to user")

			// TODO how should we handle if sending email to one user fails?
			// Now we resend to the ones who already received the notification in case the message deletion fails.

			err = w.sendEmailFn(ctx, user, e.EntityID, e.TweetContent)
			if err != nil {
				log.WithField("user", user.Id).Debug("Failed to send email to user")
				failedUsers = append(failedUsers, user.Id)
			}

			log.WithField("user", user.Id).Debug("Email sent to user")
		}

		err = w.fakeNewsQueue.DeleteMessage(ctx, e.ReceiptHandle)
		if err != nil {
			return fmt.Errorf("error deleting message from queue: %w", err)

		}

		log.Debug("Successfully deleted message")

		if len(failedUsers) > 0 {
			log.WithField("failed_to_send", failedUsers).Error("failed to send emails to all users")
		}

	}

	return nil
}

func recoverWithLog(err *error) {
	if p := recover(); p != nil {
		var tmp error
		if e, ok := p.(error); ok {
			tmp = fmt.Errorf("panic: %w", e)
		} else {
			tmp = fmt.Errorf("panic: %s", p)
		}
		logrus.WithField("Stack", string(debug.Stack())).Errorf("%s", tmp)
		if err != nil {
			*err = tmp
		}
	}
}
