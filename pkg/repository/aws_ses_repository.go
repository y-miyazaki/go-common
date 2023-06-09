package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

// AWSSESRepositoryInterface interface.
type AWSSESRepositoryInterface interface {
	// Use via SendEmailService
	SendTextEmail(from string, to []string, replyTo []string, subject, content string) (*ses.SendEmailOutput, error)
	SendHTMLEmail(from string, to []string, replyTo []string, subject, content string) (*ses.SendEmailOutput, error)
	SendEmail(from string, to []string, replyTo []string, subject, contentText, contentHTML string) (*ses.SendEmailOutput, error)
	SendBulkTemplatedEmail(from string, replyTo []string, template, defaultTemplateData string, destinations []*ses.BulkEmailDestination) (*ses.SendBulkTemplatedEmailOutput, error)
}

// AWSSESRepository struct.
type AWSSESRepository struct {
	s                    *ses.SES
	configurationSetName *string
}

// NewAWSSESRepository returns AWSSESRepository instance.
func NewAWSSESRepository(s *ses.SES, configurationSetName *string) *AWSSESRepository {
	return &AWSSESRepository{
		s:                    s,
		configurationSetName: configurationSetName,
	}
}

// SendTextEmail sends text email.
func (r *AWSSESRepository) SendTextEmail(from string, to []string, replyTo []string, subject, content string) (*ses.SendEmailOutput, error) {
	response, err := r.s.SendEmail(&ses.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(to),
		},
		ReplyToAddresses: aws.StringSlice(replyTo),
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(content),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	})

	return response, err
}

// SendHTMLEmail sends HTML email.
func (r *AWSSESRepository) SendHTMLEmail(from string, to []string, replyTo []string, subject, content string) (*ses.SendEmailOutput, error) {
	response, err := r.s.SendEmail(&ses.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(to),
		},
		ReplyToAddresses: aws.StringSlice(replyTo),
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(content),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	})

	return response, err
}

// SendEmail sends email.
func (r *AWSSESRepository) SendEmail(from string, to []string, replyTo []string, subject, contentText, contentHTML string) (*ses.SendEmailOutput, error) {
	response, err := r.s.SendEmail(&ses.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(to),
		},
		ReplyToAddresses: aws.StringSlice(replyTo),
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(contentText),
				},
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(contentHTML),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	})
	return response, err
}

// SendBulkTemplatedEmail sends bulk emails.
// Note: One or more Destination objects. All of the recipients in a Destination receive the same version of the email.
// You can specify up to 50 Destination objects within a Destinations array.
func (r *AWSSESRepository) SendBulkTemplatedEmail(from string, replyTo []string, template, defaultTemplateData string, destinations []*ses.BulkEmailDestination) (*ses.SendBulkTemplatedEmailOutput, error) {
	response, err := r.s.SendBulkTemplatedEmail(&ses.SendBulkTemplatedEmailInput{
		Source:              aws.String(from),
		Destinations:        destinations,
		ReplyToAddresses:    aws.StringSlice(replyTo),
		DefaultTemplateData: aws.String(defaultTemplateData),
		Template:            aws.String(template),
	})
	return response, err
}
