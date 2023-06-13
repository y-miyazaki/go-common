package infrastructure

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/transport"
	"go.uber.org/zap"
)

// NewAWSV2 returns ses instance.
func NewAWSV2SES(c aws.Config, optFns ...func(*sesv2.Options)) *sesv2.Client { // nolint:gocritic
	return sesv2.NewFromConfig(c, optFns...)
}

// GetSESConfig get config.
func GetAWSV2SESConfig(l *logger.Logger, key, secret, sessionToken, region, endpoint string) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-ses"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, sessionToken)),
		config.WithHTTPClient(httpClient),
	)
}

// GetSESConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSV2SESConfigNoCredentials(l *logger.Logger, region, endpoint string) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField("service", "aws-ses"),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
	)
}

// GetSESConfigZap get config.
func GetAWSV2SESConfigZap(l *logger.ZapLogger, key, secret, sessionToken, region, endpoint string) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-ses")),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, sessionToken)),
		config.WithHTTPClient(httpClient),
	)
}

// GetAWSV2SESConfigNoCredentialsZap get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSV2SESConfigNoCredentialsZap(l *logger.ZapLogger, region, endpoint string) (aws.Config, error) {
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String("service", "aws-ses")),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
	)
}
