package sqs

import (
	"context"
	"fmt"

	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/sqs"
)

type DeleteMessageFn func(ctx context.Context, receiptHandle string) error

func DeleteMessageFnBuilder(client sqs.Client, log logger.Interface) DeleteMessageFn {
	return func(ctx context.Context, receiptHandle string) error {
		err := client.DeleteMessage(ctx, receiptHandle)
		if err != nil {
			log.Error(fmt.Sprintf("error deleting message: %s", err))
			return fmt.Errorf("error deleting message: %s", err)
		}

		return nil
	}
}
