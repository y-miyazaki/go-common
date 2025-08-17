// Package infrastructure provides AWS SES integration helpers using AWS SDK v2.
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
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"go.uber.org/zap"
)

// NewAWSSES returns SES client using AWS SDK v2.
func NewAWSSES(c aws.Config, optFns ...func(*sesv2.Options)) *sesv2.Client { // nolint:gocritic
	return sesv2.NewFromConfig(c, optFns...)
}

// GetAWSSESConfig get SES config using AWS SDK v2.
func GetAWSSESConfig(l *logger.Logger, key, secret, sessionToken, region, _ string) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField(serviceFieldKey, awsSESService),
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
		return aws.Config{}, fmt.Errorf("config load ses: %w", err)
	}
	return cfg, nil
}

// GetSESConfigNoCredentials get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSSESConfigNoCredentials(l *logger.Logger, region, _ string) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPLogger(
				l.WithField(serviceFieldKey, awsSESService),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load ses no creds: %w", err)
	}
	return cfg, nil
}

// GetSESConfigZap get config.
func GetAWSSESConfigZap(l *logger.ZapLogger, key, secret, sessionToken, region, _ string) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String(serviceFieldKey, awsSESService)),
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
		return aws.Config{}, fmt.Errorf("config load ses zap: %w", err)
	}
	return cfg, nil
}

// GetAWSSESConfigNoCredentialsZap get no credentials config.
// If AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are environment variables and are in the execution environment, Credentials is not required.
func GetAWSSESConfigNoCredentialsZap(l *logger.ZapLogger, region, _ string) (aws.Config, error) { // nolint:unused
	var httpClient *http.Client
	if l != nil {
		httpClient = &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(
				l.With(zap.String(serviceFieldKey, awsSESService)),
				transport.HTTPLoggerTypeExternal,
			),
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load ses no creds zap: %w", err)
	}
	return cfg, nil
}
