package sns

import "context"

type SendNotificationEvent struct {
	UserEmail string
	Entity    string
	Content   string
}

type SendNotificationEventFn func(ctx context.Context, event SendNotificationEvent) error

func SendNotificationEventFnBuilder( /*TODO pass ses or sns client (and topic name)*/ ) SendNotificationEventFn {
	return func(ctx context.Context, event SendNotificationEvent) error {
		// TODO implement
		return nil
	}
}
