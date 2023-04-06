package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kordape/ottct-main-service/internal/handler"
)

const (
	charset  = "UTF-8"
	htmlBody = `
	<h2>Fake news spotted!</h2>
	<p> A tweet by %s contains fake news: </p>
	<p> %s </p>
	`
)

type SendFakeNewsEmailFn func(ctx context.Context, user handler.User, entityId, tweet string) error

func SendFakeNewsEmailFnBuilder(sesClient *sesv2.Client, sender string) SendFakeNewsEmailFn {
	return func(ctx context.Context, user handler.User, entity string, tweet string) error {
		input := &sesv2.SendEmailInput{
			Destination: &types.Destination{
				ToAddresses: []string{
					user.Email,
				},
			},
			Content: &types.EmailContent{
				Simple: &types.Message{
					Body: &types.Body{
						Html: &types.Content{
							Charset: aws.String(charset),
							Data:    aws.String(fmt.Sprintf(htmlBody, entity, tweet))},
					},
					Subject: &types.Content{
						Charset: aws.String(charset),
						Data:    aws.String("Fake News found"),
					},
				},
			},
			FromEmailAddress: aws.String(sender),
		}

		// Attempt to send the email.
		_, err := sesClient.SendEmail(ctx, input)
		if err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}

		return nil
	}
}
