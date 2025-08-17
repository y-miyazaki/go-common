// Package infrastructure provides AWS infrastructure configuration utilities.
package infrastructure

import (
	"context"
	"fmt"
	"net/http"

	loggerPkg "go-common/pkg/logger"
	"go-common/pkg/transport"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// AWSService represents AWS service types
type AWSService int

const (
	// AWSServiceCloudWatchLogs represents CloudWatch Logs service
	AWSServiceCloudWatchLogs AWSService = iota
	// AWSServiceCognito represents Cognito Identity Provider service
	AWSServiceCognito
	// AWSServiceS3 represents S3 service
	AWSServiceS3
	// AWSServiceSES represents SES service
	AWSServiceSES
	// AWSServiceSecretsManager represents Secrets Manager service
	AWSServiceSecretsManager
)

// String returns the string representation of the AWS service
func (s AWSService) String() string {
	switch s {
	case AWSServiceCloudWatchLogs:
		return "aws-cloudwatchlogs"
	case AWSServiceCognito:
		return "aws-cognito"
	case AWSServiceS3:
		return "aws-s3"
	case AWSServiceSES:
		return "aws-ses"
	case AWSServiceSecretsManager:
		return "aws-secretsmanager"
	default:
		return "unknown"
	}
}

// AWSConfigParams holds parameters for AWS configuration
type AWSConfigParams struct {
	Key            string
	Secret         string
	SessionToken   string
	Region         string
	Endpoint       string
	Service        AWSService
	UseCredentials bool
}

// createHTTPClientWithLogger creates HTTP client with appropriate logger transport
func createHTTPClientWithLogger(logger any) *http.Client {
	httpTransport := &http.Transport{}

	switch l := logger.(type) {
	case *loggerPkg.Logger:
		return &http.Client{
			Transport: transport.NewTransportHTTPLogger(l, transport.HTTPLoggerTypeExternal),
		}
	case *loggerPkg.ZapLogger:
		return &http.Client{
			Transport: transport.NewTransportHTTPZapLogger(l, transport.HTTPLoggerTypeExternal),
		}
	case *loggerPkg.SlogLogger:
		return &http.Client{
			Transport: transport.NewTransportHTTPSlogLogger(l, transport.HTTPLoggerTypeExternal),
		}
	default:
		// Fallback to no logging if logger type is unknown
		return &http.Client{
			Transport: httpTransport,
		}
	}
}

// createAWSConfig creates an AWS configuration with the given parameters
func createAWSConfig(params *AWSConfigParams, httpClient *http.Client) (aws.Config, error) {
	var configOptions []func(*config.LoadOptions) error

	// Add credentials if required
	if params.UseCredentials {
		configOptions = append(configOptions,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				params.Key, params.Secret, params.SessionToken)),
		)
	}

	// Add region
	configOptions = append(configOptions, config.WithRegion(params.Region))

	// Add custom endpoint if specified
	if params.Endpoint != "" {
		// Use BaseEndpoint for newer AWS SDK v2 approach
		configOptions = append(configOptions, config.WithBaseEndpoint(params.Endpoint))
	}

	// Add HTTP client if provided
	if httpClient != nil {
		configOptions = append(configOptions, config.WithHTTPClient(httpClient))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("config load: %w", err)
	}

	return cfg, nil
}

// GetAWSConfigWithLogger returns AWS config with logger support using struct parameters
func GetAWSConfigWithLogger(logger any, params *AWSConfigParams) (aws.Config, error) {
	httpClient := createHTTPClientWithLogger(logger)
	params.UseCredentials = true
	// nolint: wrapcheck
	return createAWSConfig(params, httpClient)
}

// GetAWSConfigNoCredentialsWithLogger returns AWS config without explicit credentials using struct parameters
func GetAWSConfigNoCredentialsWithLogger(logger any, params *AWSConfigParams) (aws.Config, error) {
	httpClient := createHTTPClientWithLogger(logger)
	params.UseCredentials = false
	// nolint: wrapcheck
	return createAWSConfig(params, httpClient)
}

// GetAWSConfig returns AWS config with explicit credentials using simplified parameters
func GetAWSConfig(logger any, service AWSService, key, secret, sessionToken, region, endpoint string) (aws.Config, error) {
	// nolint: wrapcheck
	return GetAWSConfigWithLogger(logger, &AWSConfigParams{
		Key:          key,
		Secret:       secret, // pragma: allowlist secret
		SessionToken: sessionToken,
		Region:       region,
		Endpoint:     endpoint,
		Service:      service,
	})
}

// GetAWSConfigNoCredentials returns AWS config without explicit credentials using simplified parameters
func GetAWSConfigNoCredentials(logger any, service AWSService, region, endpoint string) (aws.Config, error) {
	// nolint: wrapcheck
	return GetAWSConfigNoCredentialsWithLogger(logger, &AWSConfigParams{
		Region:   region,
		Endpoint: endpoint,
		Service:  service,
	})
}
