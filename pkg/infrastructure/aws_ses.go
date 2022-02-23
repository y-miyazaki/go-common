package infrastructure

import (
	"net/http"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	logrus "github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/transport"
)

// NewSES returns ses instance.
func NewSES(
	o *session.Options,
	c *aws.Config,
) *ses.SES {
	s := session.Must(session.NewSessionWithOptions(*o))
	return ses.New(s, c)
}

// GetSESConfig get config.
func GetSESConfig(e *logrus.Entry, id, secret, token, region, endpoint string) *aws.Config {
	return &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			id,
			secret,
			token),
		Region: aws.String(region),
		HTTPClient: &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				e.WithField("service", "aws-ses"),
				transport.TransportHTTPLoggerTypeExternal,
			),
		},
	}
}

// GetSESConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetSESConfigNoCredentials(e *logrus.Entry, region, endpoint string) *aws.Config {
	return &aws.Config{
		Region: aws.String(region),
		HTTPClient: &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				e.WithField("service", "aws-ses"),
				transport.TransportHTTPLoggerTypeExternal,
			),
		},
	}
}
