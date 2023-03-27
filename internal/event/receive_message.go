package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

type FakeNews struct {
	EntityId  string
	Timestamp time.Time
	Content   string
}

type ReceiveFakeNewsEventFn func(ctx context.Context) (*FakeNews, error)

func ReceiveFakeNewsEventFnBuilder(client sqs.Client, log logger.Interface) ReceiveFakeNewsEventFn {
	return func(ctx context.Context) (*FakeNews, error) {
		msg, err := client.ReceiveMessage(ctx, sqs.WithVisibilityTimeout(10))
		if err != nil {
			return nil, fmt.Errorf("error receiving message: %s", err)
		}
		if msg == nil {
			return nil, nil
		}

		var fakeNews FakeNews
		err = json.Unmarshal([]byte(msg.Body), &fakeNews)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling message body: %s", err)
		}

		return &fakeNews, nil
	}
}

type DeleteEventFn func(ctx context.Context, receiptHandle string) error

func DeleteEventFnBuilder(client sqs.Client, log logger.Interface) DeleteEventFn {
	return func(ctx context.Context, receiptHandle string) error {
		err := client.DeleteMessage(ctx, receiptHandle)
		if err != nil {
			log.Error(fmt.Sprintf("error deleting message: %s", err))
			return fmt.Errorf("error deleting message: %s", err)
		}

		return nil
	}
}
