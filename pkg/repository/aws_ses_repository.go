package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// AWSSESRepositoryInterface interface.
type AWSSESRepositoryInterface interface {
	// Use via SendEmailService
	SendTextEmail(from, to, subject, content string) error
	SendHtmlEmail(from, to, subject, content string) error
	SendEmail(from, to, subject, content string) error
}

// AWSSESRepository struct.
type AWSSESRepository struct {
	logger                         *logger.Logger
	s                              *ses.SES
	configurationSetName           *string
	isOutputLogPersonalInformation bool
}

// NewAWSSESRepository returns AWSSESRepository instance.
func NewAWSSESRepository(l *logger.Logger, s *ses.SES, configurationSetName *string, isOutputLogPersonalInformation bool) *AWSSESRepository {
	return &AWSSESRepository{
		logger:                         l,
		s:                              s,
		configurationSetName:           configurationSetName,
		isOutputLogPersonalInformation: isOutputLogPersonalInformation,
	}
}

// SendTextEmail sends text email.
func (r *AWSSESRepository) SendTextEmail(from, to, subject, content string) error {
	response, err := r.s.SendEmail(&ses.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
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

	r.log(to, subject, response, err)

	return err
}

// SendHTMLEmail sends HTML email.
func (r *AWSSESRepository) SendHTMLEmail(from, to, subject, content string) error {
	response, err := r.s.SendEmail(&ses.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
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

	r.log(to, subject, response, err)

	return err
}

// SendEmail sends email.
func (r *AWSSESRepository) SendEmail(from, to, subject, contentText, contentHTML string) error {
	response, err := r.s.SendEmail(&ses.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
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
	r.log(to, subject, response, err)
	return err
}

// SendBulkTemplatedEmail sends bulk emails.
// Note: One or more Destination objects. All of the recipients in a Destination receive the same version of the email.
//       You can specify up to 50 Destination objects within a Destinations array.
func (r *AWSSESRepository) SendBulkTemplatedEmail(from, template, defaultTemplateData string, destinations []*ses.BulkEmailDestination) error {
	response, err := r.s.SendBulkTemplatedEmail(&ses.SendBulkTemplatedEmailInput{
		DefaultTemplateData: aws.String(defaultTemplateData),
		Destinations:        destinations,
		Source:              aws.String(from),
		Template:            aws.String(template),
	})
	r.logBulkTemplated(template, defaultTemplateData, response, err)
	return err
}

func (r *AWSSESRepository) log(to, subject string, responseObject *ses.SendEmailOutput, responseError error) {
	log := r.logger

	// Check output personal information flag.
	if r.isOutputLogPersonalInformation {
		log = r.logger.
			WithField("to", to).
			WithField("subject", subject)
	}
	if responseObject.MessageId != nil {
		log = log.WithField("messageId", *responseObject.MessageId)
	}
	if responseError == nil {
		log.Debug("Successfully sent an SES email")
	} else {
		log.WithError(responseError).Error("Error while sending an SES email")
	}
}

func (r *AWSSESRepository) logBulkTemplated(template, defaultTemplateData string, responseObject *ses.SendBulkTemplatedEmailOutput, responseError error) {
	log := r.logger

	log = log.WithField("template", template).WithField("defaultTemplateData", defaultTemplateData)
	if responseObject.Status != nil {
		log = log.WithField("status", responseObject.Status)
	}
	if responseError == nil {
		log.Debug("Successfully sent an SES email")
	} else {
		log.WithError(responseError).Error("Error while sending an SES email")
	}
}
