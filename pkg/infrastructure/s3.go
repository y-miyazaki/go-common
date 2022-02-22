package infrastructure

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/transport"
)

// NewS3 return session.
func NewS3(o *session.Options, c *aws.Config) *s3.S3 {
	s := session.Must(session.NewSessionWithOptions(*o))
	return s3.New(s, c)
}

// GetDefaultOptions retrieves the options that enable SharedConfigState.
// Basically, Options assumes that an IAM role has been assigned in the execution environment.
func GetDefaultOptions() *session.Options {
	return &session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
}

// GetConfig get config.
func GetConfig(e *logrus.Entry, id, secret, token, region, endpoint string, isMinio bool) *aws.Config {
	return &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			id,
			secret,
			token),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(isMinio),
		S3ForcePathStyle: aws.Bool(isMinio),
		HTTPClient: &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				e,
				transport.TransportHTTPLoggerTypeExternal,
			),
		},
	}
}

// GetConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetConfigNoCredentials(e *logrus.Entry, region, endpoint string, isMinio bool) *aws.Config {
	return &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(isMinio),
		S3ForcePathStyle: aws.Bool(isMinio),
		HTTPClient: &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				e,
				transport.TransportHTTPLoggerTypeExternal,
			),
		},
	}
}
