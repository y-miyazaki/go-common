package infrastructure

import (
	"bytes"
	"net/http"
	"testing"

	loggerPkg "github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/transport"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAWSServiceString(t *testing.T) {
	tests := []struct {
		service  AWSService
		expected string
	}{
		{AWSServiceCloudWatchLogs, "aws-cloudwatchlogs"},
		{AWSServiceCognito, "aws-cognito"},
		{AWSServiceS3, "aws-s3"},
		{AWSServiceSES, "aws-ses"},
		{AWSServiceSecretsManager, "aws-secretsmanager"},
		{AWSService(999), "unknown"},
	}

	for _, test := range tests {
		result := test.service.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestCreateHTTPClientWithLogger(t *testing.T) {
	tests := []struct {
		name     string
		logger   interface{}
		expected bool
	}{
		{
			name:     "Logger type",
			logger:   loggerPkg.NewLogger(logrus.New()),
			expected: true,
		},
		{
			name: "ZapLogger type",
			logger: func() *loggerPkg.ZapLogger {
				config := &zap.Config{}
				config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
				return loggerPkg.NewZapLogger(config)
			}(),
			expected: true,
		},
		{
			name:     "SlogLogger type",
			logger:   loggerPkg.NewSlogLogger(&loggerPkg.SlogConfig{Level: loggerPkg.LevelInfo, Output: &bytes.Buffer{}}),
			expected: true,
		},
		{
			name:     "unknown logger type",
			logger:   "unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createHTTPClientWithLogger(tt.logger)
			assert.NotNil(t, client)
			assert.IsType(t, &http.Client{}, client)

			if tt.expected {
				// Check if transport is wrapped with logger
				switch tt.logger.(type) {
				case *loggerPkg.Logger:
					assert.IsType(t, &transport.HTTPLogger{}, client.Transport)
				case *loggerPkg.ZapLogger:
					assert.IsType(t, &transport.HTTPZapLogger{}, client.Transport)
				case *loggerPkg.SlogLogger:
					assert.IsType(t, &transport.HTTPSlogLogger{}, client.Transport)
				}
			}
		})
	}
}

func TestCreateAWSConfig(t *testing.T) {
	tests := []struct {
		name        string
		params      *AWSConfigParams
		httpClient  *http.Client
		expectError bool
	}{
		{
			name: "valid config with credentials",
			params: &AWSConfigParams{
				Key:            "test-key",
				Secret:         "test-secret",
				SessionToken:   "test-token",
				Region:         "us-east-1",
				Endpoint:       "",
				UseCredentials: true,
			},
			httpClient:  nil,
			expectError: false,
		},
		{
			name: "valid config without credentials",
			params: &AWSConfigParams{
				Region:         "us-east-1",
				UseCredentials: false,
			},
			httpClient:  nil,
			expectError: false,
		},
		{
			name: "config with custom endpoint",
			params: &AWSConfigParams{
				Region:         "us-east-1",
				Endpoint:       "http://localhost:4566",
				UseCredentials: false,
			},
			httpClient:  nil,
			expectError: false,
		},
		{
			name: "config with HTTP client",
			params: &AWSConfigParams{
				Region:         "us-east-1",
				UseCredentials: false,
			},
			httpClient:  &http.Client{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := createAWSConfig(tt.params, tt.httpClient)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, aws.Config{}, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, aws.Config{}, cfg)
				assert.Equal(t, tt.params.Region, cfg.Region)
			}
		})
	}
}

func TestGetAWSConfigWithLogger(t *testing.T) {
	logger := loggerPkg.NewLogger(logrus.New())
	params := &AWSConfigParams{
		Key:          "test-key",
		Secret:       "test-secret",
		SessionToken: "test-token",
		Region:       "us-east-1",
		Service:      AWSServiceS3,
	}

	cfg, err := GetAWSConfigWithLogger(logger, params)

	assert.NoError(t, err)
	assert.NotEqual(t, aws.Config{}, cfg)
	assert.Equal(t, params.Region, cfg.Region)
}

func TestGetAWSConfigNoCredentialsWithLogger(t *testing.T) {
	logger := loggerPkg.NewLogger(logrus.New())
	params := &AWSConfigParams{
		Region:  "us-east-1",
		Service: AWSServiceS3,
	}

	cfg, err := GetAWSConfigNoCredentialsWithLogger(logger, params)

	assert.NoError(t, err)
	assert.NotEqual(t, aws.Config{}, cfg)
	assert.Equal(t, params.Region, cfg.Region)
}

func TestGetAWSConfig(t *testing.T) {
	logger := loggerPkg.NewLogger(logrus.New())

	cfg, err := GetAWSConfig(logger, AWSServiceS3, "test-key", "test-secret", "test-token", "us-east-1", "")

	assert.NoError(t, err)
	assert.NotEqual(t, aws.Config{}, cfg)
	assert.Equal(t, "us-east-1", cfg.Region)
}

func TestGetAWSConfigNoCredentials(t *testing.T) {
	logger := loggerPkg.NewLogger(logrus.New())

	cfg, err := GetAWSConfigNoCredentials(logger, AWSServiceS3, "us-east-1", "")

	assert.NoError(t, err)
	assert.NotEqual(t, aws.Config{}, cfg)
	assert.Equal(t, "us-east-1", cfg.Region)
}
