package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

type ReceiveFakeNewsEventsFn func(ctx context.Context) ([]sqs.FakeNewsEvent, error)

func ReceiveFakeNewsEventsFnBuilder(client sqs.Client, log logger.Interface) ReceiveFakeNewsEventsFn {
	return func(ctx context.Context) ([]sqs.FakeNewsEvent, error) {
		msgs, err := client.ReceiveMessages(ctx, sqs.WithVisibilityTimeout(20), sqs.WithMaxNumberOfMessages(5))
		if err != nil {
			return nil, fmt.Errorf("error receiving message: %s", err)
		}

		events := make([]sqs.FakeNewsEvent, len(msgs))
		for i, msg := range msgs {
			var event sqs.FakeNewsEvent
			err = json.Unmarshal([]byte(msg.Body), &event)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling message body into FakeNewsEvent: %s", err)
			}

			event.ReceiptHandle = msg.ReceiptHandle

			events[i] = event
		}

		return events, nil
	}
}
