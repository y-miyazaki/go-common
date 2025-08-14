package infrastructure

import (
	"context"
	"net/http"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/transport"
	"go.uber.org/zap"
)

// NewAWSS3Session returns Session.
func NewAWSS3Session(o *session.Options) *session.Session {
	return session.Must(session.NewSessionWithOptions(*o))
}

// NewAWSS3 returns S3.
func NewAWSS3(o *session.Options, c *aws.Config) *s3.S3 {
	s := session.Must(session.NewSessionWithOptions(*o))
	return s3.New(s, c)
}

// NewAWSS3Downloader returns Downloader.
func NewAWSS3Downloader(s *session.Session) *s3manager.Downloader {
	return s3manager.NewDownloader(s)
}

// NewAWSS3Uploader returns Uploader.
func NewAWSS3Uploader(s *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(s)
}

// GetAWSS3Config get config.
func GetAWSS3Config(l *logger.Logger, id, secret, token, region, endpoint string, isMinio bool) *aws.Config {
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

// GetAWSS3ConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSS3ConfigNoCredentials(l *logger.Logger, region, endpoint string, isMinio bool) *aws.Config {
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

// GetAWSS3ConfigZap get config.
func GetAWSS3ConfigZap(l *logger.ZapLogger, id, secret, token, region, endpoint string, isMinio bool) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-s3")),
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

// GetAWSS3ConfigNoCredentialsZap get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSS3ConfigNoCredentialsZap(l *logger.ZapLogger, region, endpoint string, isMinio bool) *aws.Config {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-s3")),
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

// GetAWSV2S3Config returns AWS SDK v2 S3 configuration.
func GetAWSV2S3Config(l *logger.Logger, id, secret, token, region, endpoint string, isMinio bool) (aws2.Config, error) {
	var loadOptions []func(*config.LoadOptions) error
	if region != "" {
		loadOptions = append(loadOptions, config.WithRegion(region))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), loadOptions...)
	if err != nil {
		return cfg, err
	}
	// Note: credentials, custom endpoint, and HTTP client for v2 can be configured here when needed.
	return cfg, nil
}
