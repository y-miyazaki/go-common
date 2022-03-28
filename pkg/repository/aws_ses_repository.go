package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"
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
	e                              *logrus.Entry
	s                              *ses.SES
	configurationSetName           *string
	isOutputLogPersonalInformation bool
}

// NewAWSSESRepository returns AWSSESRepository instance.
func NewAWSSESRepository(
	e *logrus.Entry,
	s *ses.SES,
	configurationSetName *string,
	isOutputLogPersonalInformation bool,

) *AWSSESRepository {
	return &AWSSESRepository{
		e:                              e,
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

func (r *AWSSESRepository) log(
	to string,
	subject string,
	responseObject *ses.SendEmailOutput,
	responseError error,
) {
	e := r.e
	// Check output personal information flag.
	if r.isOutputLogPersonalInformation {
		e = e.
			WithField("to", to).
			WithField("subject", subject)
	}
	if responseObject.MessageId != nil {
		e = e.WithField("messageId", *responseObject.MessageId)
	}
	if responseError == nil {
		e.Info("Successfully sent an SES email")
	} else {
		e.WithError(responseError).Error("Error while sending an SES email")
	}
}
