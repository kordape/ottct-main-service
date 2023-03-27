package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/kordape/ottct-main-service/internal/sns"
	"github.com/kordape/ottct-main-service/internal/sqs"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

// const (
// 	defaultTickInterval = 10 * time.Second
// )

type Worker struct {
	log                logger.Interface
	period             time.Duration // seconds
	quit               chan bool
	receiveMessageFn   sqs.ReceiveFakeNewsEventFn
	deleteMessageFn    sqs.DeleteMessageFn
	sendNotificationFn sns.SendNotificationEventFn
}

func NewWorker(log logger.Interface, period int, receiveFn sqs.ReceiveFakeNewsEventFn, deleteFn sqs.DeleteMessageFn, sendNotificationFn sns.SendNotificationEventFn) *Worker {
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
				fakeNewsEvent, receiptHandle, err := w.receiveMessageFn(ctx)
				if err != nil {
					w.log.Error(fmt.Sprintf("error receiving message: %s", err))
					continue
				}
				if fakeNewsEvent == nil {
					w.log.Debug("no new messages available")
					continue
				}

				w.log.Debug("Received fake news event.")

				// TODO get all users subscribed to entity

				err = w.sendNotificationFn(ctx, sns.SendNotificationEvent{
					// TODO populate values
				})
				if err != nil {
					w.log.Error(fmt.Sprintf("error receiving message: %s", err))
					continue
				}
				w.log.Debug("Notifications sent.")

				err = w.deleteMessageFn(ctx, receiptHandle)
				if err != nil {
					w.log.Error(fmt.Sprintf("error deleting message from queue: %s", err))
					continue
				}
				w.log.Debug("Message successfully deleted from queue.")

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

func (w *Worker) sendNotification(context.Context, sns.SendNotificationEvent) error {
	// TODO: implement
	// send email to sqs
	return nil
}
