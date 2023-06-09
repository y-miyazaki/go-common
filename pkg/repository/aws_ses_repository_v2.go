package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// AWSV2SESRepositoryInterface interface.
type AWSV2SESRepositoryInterface interface {
	// Use via SendEmailService
	SendTextEmail(from string, to []string, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error)
	SendHTMLEmail(from string, to []string, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error)
	SendEmail(from string, to []string, replyTo []string, subject, contentText, contentHTML string) (*sesv2.SendEmailOutput, error)
	SendBulkEmail(from string, replyTo []string, template, defaultTemplateData string, BulkEmailEntries []types.BulkEmailEntry) (*sesv2.SendBulkEmailOutput, error)
}

// AWSV2SESRepository struct.
type AWSV2SESRepository struct {
	c                    *sesv2.Client
	configurationSetName *string
}

// NewAWSV2SESRepository returns AWSV2SESRepository instance.
func NewAWSV2SESRepository(c *sesv2.Client, configurationSetName *string) *AWSV2SESRepository {
	return &AWSV2SESRepository{
		c:                    c,
		configurationSetName: configurationSetName,
	}
}

// SendTextEmail sends text email.
func (r *AWSV2SESRepository) SendTextEmail(from string, to []string, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	return r.c.SendEmail(context.Background(), &sesv2.SendEmailInput{
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
}

// SendHTMLEmail sends HTML email.
func (r *AWSV2SESRepository) SendHTMLEmail(from string, to []string, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	return r.c.SendEmail(context.Background(), &sesv2.SendEmailInput{
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
}

// SendEmail sends email.
func (r *AWSV2SESRepository) SendEmail(from string, to []string, replyTo []string, subject, contentText, contentHTML string) (*sesv2.SendEmailOutput, error) {
	return r.c.SendEmail(context.Background(), &sesv2.SendEmailInput{
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
}

// SendBulkEmail sends bulk emails.
// Note: One or more Destination objects. All of the recipients in a Destination receive the same version of the email.
// You can specify up to 50 Destination objects within a Destinations array.
func (r *AWSV2SESRepository) SendBulkEmail(from string, replyTo []string, template, defaultTemplateData string, bulkEmailEntries []types.BulkEmailEntry) (*sesv2.SendBulkEmailOutput, error) {
	return r.c.SendBulkEmail(context.Background(), &sesv2.SendBulkEmailInput{
		FromEmailAddress: aws.String(from),
		ReplyToAddresses: replyTo,
		DefaultContent: &types.BulkEmailContent{
			Template: &types.Template{
				TemplateData: aws.String(defaultTemplateData),
			},
		},
		BulkEmailEntries: bulkEmailEntries,
	})
}
