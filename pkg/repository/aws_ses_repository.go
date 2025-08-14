package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// AWSSESRepositoryInterface interface.
// nolint:iface,revive,unused
type AWSSESRepositoryInterface interface {
	SendTextEmail(from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error)
	SendHTMLEmail(from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error)
	SendEmail(from string, to, replyTo []string, subject, contentText, contentHTML string) (*sesv2.SendEmailOutput, error)
	SendBulkEmail(from string, replyTo []string, defaultTemplateData string, bulkEmailEntries []types.BulkEmailEntry) (*sesv2.SendBulkEmailOutput, error)
}

// AWSSESRepository struct.
type AWSSESRepository struct {
	c                    *sesv2.Client
	configurationSetName *string
}

// NewAWSSESRepository returns AWSSESRepository instance.
func NewAWSSESRepository(c *sesv2.Client, configurationSetName *string) *AWSSESRepository {
	return &AWSSESRepository{
		c:                    c,
		configurationSetName: configurationSetName,
	}
}

// SendTextEmail sends text email.
func (r *AWSSESRepository) SendTextEmail(from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(context.Background(), &sesv2.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		Destination: &types.Destination{
			ToAddresses: to,
		},
		ReplyToAddresses: replyTo,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(content),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendTextEmail: %w", err)
	}
	return res, nil
}

// SendHTMLEmail sends HTML email.
func (r *AWSSESRepository) SendHTMLEmail(from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(context.Background(), &sesv2.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		Destination: &types.Destination{
			ToAddresses: to,
		},
		ReplyToAddresses: replyTo,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(content),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendHTMLEmail: %w", err)
	}
	return res, nil
}

// SendEmail sends email.
func (r *AWSSESRepository) SendEmail(from string, to, replyTo []string, subject, contentText, contentHTML string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(context.Background(), &sesv2.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		Destination: &types.Destination{
			ToAddresses: to,
		},
		ReplyToAddresses: replyTo,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(contentText),
					},
					Html: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(contentHTML),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendEmail: %w", err)
	}
	return res, nil
}

// SendBulkEmail sends bulk emails.
// Note: One or more Destination objects. All of the recipients in a Destination receive the same version of the email.
// You can specify up to 50 Destination objects within a Destinations array.
func (r *AWSSESRepository) SendBulkEmail(from string, replyTo []string, defaultTemplateData string, bulkEmailEntries []types.BulkEmailEntry) (*sesv2.SendBulkEmailOutput, error) {
	res, err := r.c.SendBulkEmail(context.Background(), &sesv2.SendBulkEmailInput{
		FromEmailAddress: aws.String(from),
		ReplyToAddresses: replyTo,
		DefaultContent: &types.BulkEmailContent{
			Template: &types.Template{
				TemplateData: aws.String(defaultTemplateData),
			},
		},
		BulkEmailEntries: bulkEmailEntries,
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendBulkEmail: %w", err)
	}
	return res, nil
}
