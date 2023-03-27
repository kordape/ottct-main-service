package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

type ReceiveFakeNewsEventFn func(ctx context.Context) (*sqs.FakeNewsEvent, string, error)

func ReceiveFakeNewsEventFnBuilder(client sqs.Client, log logger.Interface) ReceiveFakeNewsEventFn {
	return func(ctx context.Context) (*sqs.FakeNewsEvent, string, error) {
		msg, err := client.ReceiveMessage(ctx, sqs.WithVisibilityTimeout(10))
		if err != nil {
			return nil, "", fmt.Errorf("error receiving message: %s", err)
		}
		if msg == nil {
			return nil, "", nil
		}

		var fakeNews sqs.FakeNewsEvent
		err = json.Unmarshal([]byte(msg.Body), &fakeNews)
		if err != nil {
			return nil, "", fmt.Errorf("error unmarshalling message body into FakeNewsEvent: %s", err)
		}

		return &fakeNews, msg.ReceiptHandle, nil
	}
}
