// Package infrastructure provides AWS S3 integration helpers using AWS SDK v2.
package infrastructure

import (
	"context"
	"fmt"
	"net/http"

	"go-common/pkg/logger"
	"go-common/pkg/transport"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

const (
	serviceFieldKey = "service"
	awsS3Service    = "aws-s3"
	awsSESService   = "aws-ses"
)

// NewAWSS3 returns S3 client using AWS SDK v2.
func NewAWSS3(c *aws.Config, optFns ...func(*s3.Options)) *s3.Client { // nolint:gocritic
	return s3.NewFromConfig(*c, optFns...)
}

// NewAWSS3Downloader returns S3 Downloader using AWS SDK v2.
func NewAWSS3Downloader(c manager.DownloadAPIClient, options ...func(*manager.Downloader)) *manager.Downloader {
	return manager.NewDownloader(c, options...)
}

// NewAWSS3Uploader returns S3 Uploader using AWS SDK v2.
func NewAWSS3Uploader(c manager.UploadAPIClient, options ...func(*manager.Uploader)) *manager.Uploader {
	return manager.NewUploader(c, options...)
}

// GetAWSS3Config get S3 config using AWS SDK v2.
func GetAWSS3Config(l *logger.Logger, key, secret, sessionToken, region, endpoint string, _ bool) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField(serviceFieldKey, awsS3Service),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, sessionToken)),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load: %w", err)
	}
	return cfg, nil
}

// GetAWSS3ConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSS3ConfigNoCredentials(l *logger.Logger, region, endpoint string, _ bool) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField(serviceFieldKey, awsS3Service),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load no credentials: %w", err)
	}
	return cfg, nil
}

// GetAWSS3ConfigZap get config.
func GetAWSS3ConfigZap(l *logger.ZapLogger, key, secret, sessionToken, region, endpoint string, _ bool) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String(serviceFieldKey, awsS3Service)),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, sessionToken)),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load zap: %w", err)
	}
	return cfg, nil
}

// GetAWSS3ConfigNoCredentialsZap get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSS3ConfigNoCredentialsZap(l *logger.ZapLogger, region, endpoint string, _ bool) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String(serviceFieldKey, awsS3Service)),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load no creds zap: %w", err)
	}
	return cfg, nil
}
