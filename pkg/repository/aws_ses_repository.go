package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// AWSSESClientInterface defines the interface for SES client operations
type AWSSESClientInterface interface {
	// SendEmail sends an email using Amazon SES
	SendEmail(_ context.Context, _ *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
	// SendBulkEmail sends bulk emails using Amazon SES
	SendBulkEmail(_ context.Context, _ *sesv2.SendBulkEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendBulkEmailOutput, error)
}

// AWSSESRepository struct.
type AWSSESRepository struct {
	c                    AWSSESClientInterface
	configurationSetName *string
}

// NewAWSSESRepository returns AWSSESRepository instance.
func NewAWSSESRepository(c *sesv2.Client, configurationSetName *string) *AWSSESRepository {
	return &AWSSESRepository{
		c:                    c,
		configurationSetName: configurationSetName,
	}
}

// NewAWSSESRepositoryWithInterface returns AWSSESRepository instance with interface (for testing).
func NewAWSSESRepositoryWithInterface(c AWSSESClientInterface, configurationSetName *string) *AWSSESRepository {
	return &AWSSESRepository{
		c:                    c,
		configurationSetName: configurationSetName,
	}
}

// SendTextEmail sends text email.
func (r *AWSSESRepository) SendTextEmail(ctx context.Context, from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(ctx, &sesv2.SendEmailInput{
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
func (r *AWSSESRepository) SendHTMLEmail(ctx context.Context, from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(ctx, &sesv2.SendEmailInput{
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
func (r *AWSSESRepository) SendEmail(ctx context.Context, from string, to, replyTo []string, subject, contentText, contentHTML string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(ctx, &sesv2.SendEmailInput{
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
func (r *AWSSESRepository) SendBulkEmail(ctx context.Context, from string, replyTo []string, defaultTemplateData string, bulkEmailEntries []types.BulkEmailEntry) (*sesv2.SendBulkEmailOutput, error) {
	res, err := r.c.SendBulkEmail(ctx, &sesv2.SendBulkEmailInput{
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
