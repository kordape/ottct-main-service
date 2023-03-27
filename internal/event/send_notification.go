package event

import "context"

type SendNotificationEvent struct {
	UserEmail string
	Entity    string
	Content   string
}

type SendNotificationFn func(ctx context.Context, event SendNotificationEvent) error

func SendNotificationFnBuilder( /*TODO pass ses or sns client (and topic name)*/ ) SendNotificationFn {
	return func(ctx context.Context, event SendNotificationEvent) error {
		// TODO implement
		return nil
	}
}
