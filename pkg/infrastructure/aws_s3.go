package infrastructure

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/transport"
)

// NewS3Session returns Session.
func NewS3Session(o *session.Options) *session.Session {
	return session.Must(session.NewSessionWithOptions(*o))
}

// NewS3 returns S3.
func NewS3(o *session.Options, c *aws.Config) *s3.S3 {
	s := session.Must(session.NewSessionWithOptions(*o))
	return s3.New(s, c)
}

// NewDownloader returns Downloader.
func NewDownloader(s *session.Session) *s3manager.Downloader {
	return s3manager.NewDownloader(s)
}

// NewUploader returns Uploader.
func NewUploader(s *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(s)
}

// GetS3DefaultOptions retrieves the options that enable SharedConfigState.
// Basically, Options assumes that an IAM role has been assigned in the execution environment.
func GetS3DefaultOptions() *session.Options {
	return &session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
}

// GetS3Config get config.
func GetS3Config(l *logger.Logger, id, secret, token, region, endpoint string, isMinio bool) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-s3"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			id,
			secret,
			token),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(isMinio),
		S3ForcePathStyle: aws.Bool(isMinio),
		HTTPClient:       httpClient,
	}
}

// GetS3ConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetS3ConfigNoCredentials(l *logger.Logger, region, endpoint string, isMinio bool) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-s3"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(isMinio),
		S3ForcePathStyle: aws.Bool(isMinio),
		HTTPClient:       httpClient,
	}
}
