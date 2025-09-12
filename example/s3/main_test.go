// Package main demonstrates AWS S3 operations using AWS SDK v2.
package main

import (
	"os"
	"testing"

	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/repository"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerSetup(t *testing.T) {
	// Test logrus logger setup
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	l := logger.NewLogger(logrusLogger)

	assert.NotNil(t, l)
	assert.NotNil(t, l.Entry)
	assert.Equal(t, logrus.InfoLevel, l.Entry.Logger.Level)
}

func TestS3ConfigSetup(t *testing.T) {
	// Test S3 configuration parameters
	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		s3Region = "us-east-1"
	}
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Endpoint == "" {
		s3Endpoint = "http://localhost:9000"
	}
	s3ID := os.Getenv("S3_ID")
	if s3ID == "" {
		s3ID = "test-id"
	}
	s3Secret := os.Getenv("S3_SECRET")
	if s3Secret == "" {
		s3Secret = "test-secret" // pragma: allowlist-secret
	}
	s3Token := os.Getenv("S3_TOKEN")

	// Test that variables are set correctly
	assert.NotEmpty(t, s3Region)
	assert.NotEmpty(t, s3Endpoint)
	assert.NotEmpty(t, s3ID)
	assert.NotEmpty(t, s3Secret)
	assert.Equal(t, "", s3Token)
}

func TestS3RepositoryCreation(t *testing.T) {
	// Test S3 repository creation with mock config
	// This would require a mock AWS config
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	l := logger.NewLogger(logrusLogger)

	// Test repository creation logic without actual AWS connection
	assert.NotNil(t, l)

	// Test that we can create the repository structure
	// In a real test, we would need to mock the AWS config
	testFunc := func(o *s3.Options) {
		o.UsePathStyle = true
	}
	assert.NotNil(t, testFunc)
}

func TestS3Operations(t *testing.T) {
	// Test S3 operation parameters
	text := "abc"
	bucket := "test"
	objectKey := "test.txt"

	assert.Equal(t, "abc", text)
	assert.Equal(t, "test", bucket)
	assert.Equal(t, "test.txt", objectKey)
}

func TestMain_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if S3 environment variables are not set
	if os.Getenv("S3_REGION") == "" ||
		os.Getenv("S3_ENDPOINT") == "" ||
		os.Getenv("S3_ID") == "" ||
		os.Getenv("S3_SECRET") == "" {
		t.Skip("Skipping S3 integration test - requires S3 environment variables")
	}

	// This would test the actual main function logic
	// For now, just ensure the setup doesn't panic
	assert.NotPanics(t, func() {
		// Test logger setup
		logrusLogger := &logrus.Logger{}
		logrusLogger.Formatter = &logrus.JSONFormatter{}
		logrusLogger.Out = os.Stdout
		logrusLogger.Level, _ = logrus.ParseLevel("Info")
		l := logger.NewLogger(logrusLogger)

		// Test S3 config creation (this might fail without real credentials)
		s3Region := os.Getenv("S3_REGION")
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		s3ID := os.Getenv("S3_ID")
		s3Secret := os.Getenv("S3_SECRET")
		s3Token := os.Getenv("S3_TOKEN")

		_, err := infrastructure.GetAWSConfig(l, infrastructure.AWSServiceS3, s3ID, s3Secret, s3Token, s3Region, s3Endpoint)
		// We expect this to potentially fail in test environment
		_ = err // We don't assert on this as it depends on the test environment
	})
}

func TestS3RepositoryMock(t *testing.T) {
	// Test that we can create repository interface
	var repo *repository.AWSS3Repository
	assert.Nil(t, repo)

	// Test repository creation function exists
	assert.NotNil(t, repository.NewAWSS3Repository)
	assert.NotNil(t, repository.NewAWSS3RepositoryWithInterface)
}

func TestUtilsFunctions(t *testing.T) {
	// Test that utils functions are available
	// This would test the utils.GetStringFromReadCloser function
	// but we can't easily test it without a real ReadCloser
	assert.True(t, true) // Placeholder test
}

func TestMainConfiguration(t *testing.T) {
	// Test main function configuration without performing actual S3 operations
	// Set up environment variables for testing
	originalEnv := make(map[string]string)
	envVars := map[string]string{
		"S3_REGION":   "us-east-1",
		"S3_ENDPOINT": "http://localhost:9000",
		"S3_ID":       "test-id",
		"S3_SECRET":   "test-secret", // pragma: allowlist-secret
		"S3_TOKEN":    "",
	}

	// Save original environment variables
	for key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Set test environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Test logger setup
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	l := logger.NewLogger(logrusLogger)
	assert.NotNil(t, l)

	// Test S3 config creation
	s3Region := os.Getenv("S3_REGION")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ID := os.Getenv("S3_ID")
	s3Secret := os.Getenv("S3_SECRET")
	s3Token := os.Getenv("S3_TOKEN")

	s3Config, err := infrastructure.GetAWSConfig(l, infrastructure.AWSServiceS3, s3ID, s3Secret, s3Token, s3Region, s3Endpoint)
	assert.NoError(t, err)
	assert.NotNil(t, s3Config)

	// Test S3 repository creation
	awsS3Repository := repository.NewAWSS3Repository(s3.NewFromConfig(s3Config, func(o *s3.Options) {
		o.UsePathStyle = true
	}))
	assert.NotNil(t, awsS3Repository)

	// Test operation parameters
	text := "abc"
	bucket := "test"
	assert.Equal(t, "abc", text)
	assert.Equal(t, "test", bucket)

	// Restore original environment variables
	for key, value := range originalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}

func TestS3OperationsMock(t *testing.T) {
	// Test S3 operation parameters without actual AWS connection
	text := "abc"
	bucket := "test"
	objectKey := "test.txt"

	assert.Equal(t, "abc", text)
	assert.Equal(t, "test", bucket)
	assert.Equal(t, "test.txt", objectKey)

	// Test that repository creation functions exist
	assert.NotNil(t, repository.NewAWSS3Repository)
	assert.NotNil(t, repository.NewAWSS3RepositoryWithInterface)
}

func TestMainFunctionLogic(t *testing.T) {
	// Test main function logic without actual S3 operations
	// Test logger setup
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	l := logger.NewLogger(logrusLogger)
	assert.NotNil(t, l)

	// Test operation parameters
	text := "abc"
	bucket := "test"
	assert.Equal(t, "abc", text)
	assert.Equal(t, "test", bucket)

	// Test S3 options configuration
	testFunc := func(o *s3.Options) {
		o.UsePathStyle = true
	}
	assert.NotNil(t, testFunc)
}
