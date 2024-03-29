package infrastructure

import (
	"net/http"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/transport"
	"go.uber.org/zap"
)

// NewAWSSES returns ses instance.
func NewAWSSES(
	o *session.Options,
	c *aws.Config,
) *ses.SES {
	s := session.Must(session.NewSessionWithOptions(*o))
	return ses.New(s, c)
}

// GetAWSSESConfig get config.
func GetAWSSESConfig(l *logger.Logger, id, secret, token, region, endpoint string) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-ses"),
				transport.HTTPLoggerTypeExternal,
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

// GetAWSSESConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSSESConfigNoCredentials(l *logger.Logger, region, endpoint string) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-ses"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return &aws.Config{
		Region:     aws.String(region),
		HTTPClient: httpClient,
	}
}

// GetAWSSESConfigZap get config.
func GetAWSSESConfigZap(l *logger.ZapLogger, id, secret, token, region, endpoint string) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-ses")),
				transport.HTTPLoggerTypeExternal,
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

// GetAWSSESConfigNoCredentialsZap get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetSESConfigNoCredentialsZap(l *logger.ZapLogger, region, endpoint string) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-ses")),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return &aws.Config{
		Region:     aws.String(region),
		HTTPClient: httpClient,
	}
}
