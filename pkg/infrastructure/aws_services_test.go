package infrastructure

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSCloudWatchLogs(t *testing.T) {
	cfg := &aws.Config{Region: "us-east-1"}
	client := NewAWSCloudWatchLogs(cfg)

	assert.NotNil(t, client)
	assert.IsType(t, &cloudwatchlogs.Client{}, client)
}

func TestNewAWSCognito(t *testing.T) {
	cfg := &aws.Config{Region: "us-east-1"}
	client := NewAWSCognito(cfg)

	assert.NotNil(t, client)
	assert.IsType(t, &cognitoidentityprovider.Client{}, client)
}

func TestNewAWSS3(t *testing.T) {
	cfg := &aws.Config{Region: "us-east-1"}
	client := NewAWSS3(cfg)

	assert.NotNil(t, client)
	assert.IsType(t, &s3.Client{}, client)
}

func TestNewAWSS3Downloader(t *testing.T) {
	cfg := &aws.Config{Region: "us-east-1"}
	s3Client := NewAWSS3(cfg)
	downloader := NewAWSS3Downloader(s3Client)

	assert.NotNil(t, downloader)
	assert.IsType(t, &manager.Downloader{}, downloader)
}

func TestNewAWSS3Uploader(t *testing.T) {
	cfg := &aws.Config{Region: "us-east-1"}
	s3Client := NewAWSS3(cfg)
	uploader := NewAWSS3Uploader(s3Client)

	assert.NotNil(t, uploader)
	assert.IsType(t, &manager.Uploader{}, uploader)
}

func TestNewAWSSecretsManager(t *testing.T) {
	cfg := &aws.Config{Region: "us-east-1"}
	client := NewAWSSecretsManager(cfg)

	assert.NotNil(t, client)
	assert.IsType(t, &secretsmanager.Client{}, client)
}

func TestNewAWSSES(t *testing.T) {
	cfg := aws.Config{Region: "us-east-1"}
	client := NewAWSSES(cfg)

	assert.NotNil(t, client)
	assert.IsType(t, &sesv2.Client{}, client)
}
