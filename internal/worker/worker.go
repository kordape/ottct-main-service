package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/kordape/ottct-main-service/internal/sns"
	"github.com/kordape/ottct-main-service/internal/sqs"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

type Worker struct {
	log                logger.Interface
	period             time.Duration // seconds
	quit               chan bool
	receiveMessagesFn  sqs.ReceiveFakeNewsEventsFn
	deleteMessageFn    sqs.DeleteMessageFn
	sendNotificationFn sns.SendNotificationEventFn
}

func NewWorker(log logger.Interface, period int, receiveFn sqs.ReceiveFakeNewsEventsFn, deleteFn sqs.DeleteMessageFn, sendNotificationFn sns.SendNotificationEventFn) *Worker {
	return &Worker{
		log:                log,
		period:             time.Duration(period),
		quit:               make(chan bool),
		receiveMessagesFn:  receiveFn,
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
				events, err := w.receiveMessagesFn(ctx)
				if err != nil {
					w.log.Error(fmt.Sprintf("error receiving message: %s", err))
					continue
				}
				if len(events) == 0 {
					w.log.Debug("no new messages available")
					continue
				}

				w.log.Debug("Received fake news events.")

				// TODO get all users subscribed to entity

				for _, e := range events {
					err = w.sendNotificationFn(ctx, sns.SendNotificationEvent{
						// TODO populate values
					})
					if err != nil {
						w.log.Error(fmt.Sprintf("error sending notification: %s", err))
						continue
					}
					w.log.Debug("Notifications sent.")

					err = w.deleteMessageFn(ctx, e.ReceiptHandle)
					if err != nil {
						w.log.Error(fmt.Sprintf("error deleting message from queue: %s", err))
						continue
					}

					w.log.Debug("Message successfully deleted from queue.")
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
