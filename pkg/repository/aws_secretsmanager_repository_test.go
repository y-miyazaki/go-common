package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSecretsManagerClient is a mock implementation of the SecretsManager client
type MockSecretsManagerClient struct {
	mock.Mock
}

func (m *MockSecretsManagerClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*secretsmanager.GetSecretValueOutput), args.Error(1)
}

func TestNewAWSSecretsManagerRepository(t *testing.T) {
	mockClient := &secretsmanager.Client{}
	repo := NewAWSSecretsManagerRepository(mockClient)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.Client)
}

func TestAWSSecretsManagerRepository_GetSecretString_SecretString(t *testing.T) {
	mockClient := new(MockSecretsManagerClient)
	repo := &AWSSecretsManagerRepository{Client: mockClient}

	secretName := "test-secret"
	expectedSecret := "secret-value"

	mockClient.On("GetSecretValue", mock.Anything, mock.MatchedBy(func(input *secretsmanager.GetSecretValueInput) bool {
		return *input.SecretId == secretName && *input.VersionStage == "AWSCURRENT"
	}), mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(expectedSecret),
	}, nil)

	result, err := repo.GetSecretString(context.Background(), secretName)

	assert.NoError(t, err)
	assert.Equal(t, expectedSecret, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSecretsManagerRepository_GetSecretString_SecretBinary(t *testing.T) {
	mockClient := new(MockSecretsManagerClient)
	repo := &AWSSecretsManagerRepository{Client: mockClient}

	secretName := "test-secret-binary"
	secretValue := "binary-secret-value"
	encodedSecret := base64.StdEncoding.EncodeToString([]byte(secretValue))

	mockClient.On("GetSecretValue", mock.Anything, mock.MatchedBy(func(input *secretsmanager.GetSecretValueInput) bool {
		return *input.SecretId == secretName && *input.VersionStage == "AWSCURRENT"
	}), mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: nil,
		SecretBinary: []byte(encodedSecret),
	}, nil)

	result, err := repo.GetSecretString(context.Background(), secretName)

	assert.NoError(t, err)
	assert.Equal(t, secretValue, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSecretsManagerRepository_GetSecretString_Error(t *testing.T) {
	mockClient := new(MockSecretsManagerClient)
	repo := &AWSSecretsManagerRepository{Client: mockClient}

	secretName := "test-secret-error"
	expectedError := errors.New("secretsmanager error")

	mockClient.On("GetSecretValue", mock.Anything, mock.MatchedBy(func(input *secretsmanager.GetSecretValueInput) bool {
		return *input.SecretId == secretName && *input.VersionStage == "AWSCURRENT"
	}), mock.Anything).Return(nil, expectedError)

	result, err := repo.GetSecretString(context.Background(), secretName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secretsmanager GetSecretValue")
	assert.Empty(t, result)
	mockClient.AssertExpectations(t)
}
