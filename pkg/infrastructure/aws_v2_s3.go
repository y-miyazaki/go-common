package infrastructure

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/transport"
	"go.uber.org/zap"
)

// NewAWSS3 returns S3.
func NewAWSV2S3(c *aws.Config, optFns ...func(*s3.Options)) *s3.Client { // nolint:gocritic
	return s3.NewFromConfig(*c, optFns...)
}

// NewAWSV2S3Downloader returns Downloader.
func NewAWSV2S3Downloader(c manager.DownloadAPIClient, options ...func(*manager.Downloader)) *manager.Downloader {
	return manager.NewDownloader(c, options...)
}

// NewAWSV2S3Uploader returns Uploader.
func NewAWSV2S3Uploader(c manager.UploadAPIClient, options ...func(*manager.Uploader)) *manager.Uploader {
	return manager.NewUploader(c, options...)
}

// GetAWSV2S3Config get config.
func GetAWSV2S3Config(l *logger.Logger, key, secret, sessionToken, region, endpoint string, isMinio bool) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-s3"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, sessionToken)),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if endpoint != "" {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           endpoint,
						SigningRegion: region,
					}, nil
				}
				return aws.Endpoint{}, nil
			})),
		config.WithHTTPClient(httpClient),
	)
}

// GetAWSV2S3ConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSV2S3ConfigNoCredentials(l *logger.Logger, region, endpoint string, isMinio bool) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-s3"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			})),
		config.WithHTTPClient(httpClient),
	)
}

// GetAWSV2S3ConfigZap get config.
func GetAWSV2S3ConfigZap(l *logger.ZapLogger, key, secret, sessionToken, region, endpoint string, isMinio bool) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-s3")),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, sessionToken)),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			})),
		config.WithHTTPClient(httpClient),
	)
}

// GetAWSV2S3ConfigNoCredentialsZap get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSV2S3ConfigNoCredentialsZap(l *logger.ZapLogger, region, endpoint string, isMinio bool) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-s3")),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			})),
		config.WithHTTPClient(httpClient),
	)
}
