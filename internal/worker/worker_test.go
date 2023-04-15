package worker

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/internal/mocks"
	"github.com/kordape/ottct-main-service/internal/ses"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

var hook *test.Hook

type mockQueueClient struct {
	wantMsg             bool
	wantReceiveMsgError bool
	wantDeleteMsgError  bool
}

func (c mockQueueClient) ReceiveMessages(ctx context.Context, options ...sqs.ReceiveOption) ([]sqs.Message, error) {
	if c.wantReceiveMsgError {
		return nil, fmt.Errorf("error")
	}

	if c.wantMsg {
		return []sqs.Message{
			{
				ReceiptHandle: "testHandle",
				Body:          `{"tweetContent":"testing","entityId":"1","tweetTimestamp":"2009-11-10T23:00:00Z"}`,
			},
		}, nil
	}

	return nil, nil

}

func (c mockQueueClient) DeleteMessage(ctx context.Context, receiptHandle string) error {
	if c.wantDeleteMsgError {
		return fmt.Errorf("error")
	}

	return nil
}

func mockSendEmailFn(wantError bool) ses.SendFakeNewsEmailFn {
	return func(ctx context.Context, user handler.User, entityId, tweet string) error {
		if wantError {
			return fmt.Errorf("error")
		}

		return nil
	}
}

func TestWorker_work(t *testing.T) {
	storage := mocks.NewMockSubscriptionStorage(gomock.NewController(t))
	storage.EXPECT().GetSubscriptionsByEntity("1").Return([]handler.User{{Id: 1}}, nil).MinTimes(1)

	sm, _ := handler.NewSubscriptionManager(storage, &validator.Validate{})

	type args struct {
		log                  *logrus.Entry
		subscriptionsManager *handler.SubscriptionManager
	}

	tests := []struct {
		name    string
		worker  *Worker
		args    args
		wantErr bool
	}{
		{
			name: "Successfully read msg and sent notification",
			worker: NewWorker(
				1,
				mockQueueClient{
					wantMsg:             true,
					wantReceiveMsgError: false,
					wantDeleteMsgError:  false,
				},
				mockSendEmailFn(false),
			),
			args: args{
				subscriptionsManager: sm,
				log:                  logrus.NewEntry(logrus.StandardLogger()),
			},
			wantErr: false,
		},
		{
			name: "No new messages",
			worker: NewWorker(
				1,
				mockQueueClient{
					wantMsg:             false,
					wantReceiveMsgError: false,
					wantDeleteMsgError:  false,
				},
				mockSendEmailFn(false),
			),
			args: args{
				subscriptionsManager: sm,
				log:                  logrus.NewEntry(logrus.StandardLogger()),
			},
			wantErr: false,
		},
		{
			name: "Receive msg error",
			worker: NewWorker(
				1,
				mockQueueClient{
					wantMsg:             false,
					wantReceiveMsgError: true,
					wantDeleteMsgError:  false,
				},
				mockSendEmailFn(false),
			),
			args: args{
				subscriptionsManager: sm,
				log:                  logrus.NewEntry(logrus.StandardLogger()),
			},
			wantErr: true,
		},
		{
			name: "Delete message error",
			worker: NewWorker(
				1,
				mockQueueClient{
					wantMsg:             true,
					wantReceiveMsgError: false,
					wantDeleteMsgError:  true,
				},
				mockSendEmailFn(false),
			),
			args: args{
				subscriptionsManager: sm,
				log:                  logrus.NewEntry(logrus.StandardLogger()),
			},
			wantErr: true,
		},
		{
			name: "Panic handled when no logger is provided",
			worker: NewWorker(
				1,
				mockQueueClient{
					wantMsg:             true,
					wantReceiveMsgError: false,
					wantDeleteMsgError:  true,
				},
				mockSendEmailFn(false),
			),
			args: args{
				subscriptionsManager: sm,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.worker.work(context.TODO(), tt.args.log, tt.args.subscriptionsManager); (err != nil) != tt.wantErr {
				t.Errorf("Worker.work() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ExampleWorker_work_failedToSendNotificationsToAllUsers() {
	hook.Reset()
	storage := mocks.NewMockSubscriptionStorage(gomock.NewController(nil))
	storage.EXPECT().GetSubscriptionsByEntity("1").Return([]handler.User{{Id: 1}}, nil).MinTimes(1)

	sm, _ := handler.NewSubscriptionManager(storage, &validator.Validate{})

	worker := NewWorker(
		1,
		mockQueueClient{
			wantMsg:             true,
			wantReceiveMsgError: false,
			wantDeleteMsgError:  false,
		},
		mockSendEmailFn(true),
	)

	worker.work(context.TODO(), logrus.NewEntry(logrus.StandardLogger()), sm)

	printEntries(hook)
	// Output:
	// ERROR: failed to send emails to all users

}

func printEntries(hook *test.Hook) {
	for _, e := range hook.AllEntries() {
		fmt.Printf("%s: %s\n", strings.ToUpper(e.Level.String()), e.Message)
	}
}

func init() {
	_, temp := test.NewNullLogger()
	logrus.AddHook(temp)
	hook = temp
}
