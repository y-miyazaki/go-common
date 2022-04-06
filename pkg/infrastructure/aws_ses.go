package infrastructure

import (
	"net/http"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/y-miyazaki/go-common/pkg/logger"
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
func GetSESConfig(logger *logger.Logger, id, secret, token, region, endpoint string) *aws.Config {
	var httpClient *http.Client
	if logger != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				logger.WithField("service", "aws-ses"),
				transport.TransportHTTPLoggerTypeExternal,
			),
		}
	}
	return &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			id,
			secret,
			token),
		Region:     aws.String(region),
		HTTPClient: httpClient,
	}
}

// GetSESConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetSESConfigNoCredentials(logger *logger.Logger, region, endpoint string) *aws.Config {
	var httpClient *http.Client
	if logger != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				logger.WithField("service", "aws-ses"),
				transport.TransportHTTPLoggerTypeExternal,
			),
		}
	}
	return &aws.Config{
		Region:     aws.String(region),
		HTTPClient: httpClient,
	}
}
