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
	ses                            *ses.SES
	configurationSetName           *string
	isOutputLogPersonalInformation bool
}

// NewAWSSERepository returns AWSSESRepository instance.
func NewAWSSERepository(
	e *logrus.Entry,
	ses *ses.SES,
	configurationSetName *string,
	isOutputLogPersonalInformation bool,

) *AWSSESRepository {
	return &AWSSESRepository{
		e:                              e,
		ses:                            ses,
		configurationSetName:           configurationSetName,
		isOutputLogPersonalInformation: isOutputLogPersonalInformation,
	}
}

func (r *AWSSESRepository) SendTextEmail(from, to, subject, content string) error {
	response, err := r.ses.SendEmail(&ses.SendEmailInput{
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

func (r *AWSSESRepository) SendHtmlEmail(from, to, subject, content string) error {
	response, err := r.ses.SendEmail(&ses.SendEmailInput{
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

func (r *AWSSESRepository) SendEmail(from, to, subject, contentText, contentHtml string) error {
	response, err := r.ses.SendEmail(&ses.SendEmailInput{
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
					Data:    aws.String(contentHtml),
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
	if r.isOutputLogPersonalInformation {
		e := e.
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
