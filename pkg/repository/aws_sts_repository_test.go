package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSTSClient is a mock implementation of the STS client.
type MockSTSClient struct {
	mock.Mock
}

func (m *MockSTSClient) GetAccessKeyInfo(ctx context.Context, params *sts.GetAccessKeyInfoInput, optFns ...func(*sts.Options)) (*sts.GetAccessKeyInfoOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetAccessKeyInfoOutput), args.Error(1)
}

func (m *MockSTSClient) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetCallerIdentityOutput), args.Error(1)
}

func (m *MockSTSClient) GetDelegatedAccessToken(ctx context.Context, params *sts.GetDelegatedAccessTokenInput, optFns ...func(*sts.Options)) (*sts.GetDelegatedAccessTokenOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetDelegatedAccessTokenOutput), args.Error(1)
}

func (m *MockSTSClient) GetFederationToken(ctx context.Context, params *sts.GetFederationTokenInput, optFns ...func(*sts.Options)) (*sts.GetFederationTokenOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetFederationTokenOutput), args.Error(1)
}

func (m *MockSTSClient) GetSessionToken(ctx context.Context, params *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetSessionTokenOutput), args.Error(1)
}

func (m *MockSTSClient) GetWebIdentityToken(ctx context.Context, params *sts.GetWebIdentityTokenInput, optFns ...func(*sts.Options)) (*sts.GetWebIdentityTokenOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetWebIdentityTokenOutput), args.Error(1)
}

func TestNewAWSSTSRepository(t *testing.T) {
	mockClient := &sts.Client{}
	repo := NewAWSSTSRepository(mockClient)

	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.Client)
}

func TestAWSSTSRepository_GetCallerIdentity(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expected := &sts.GetCallerIdentityOutput{Account: aws.String("123456789012")}

	mockClient.On("GetCallerIdentity", mock.Anything, mock.MatchedBy(func(input *sts.GetCallerIdentityInput) bool {
		return input != nil
	}), mock.Anything).Return(expected, nil)

	result, err := repo.GetCallerIdentity(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetAccessKeyInfo(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	expected := &sts.GetAccessKeyInfoOutput{Account: aws.String("123456789012")}

	mockClient.On("GetAccessKeyInfo", mock.Anything, mock.MatchedBy(func(input *sts.GetAccessKeyInfoInput) bool {
		return input.AccessKeyId != nil && *input.AccessKeyId == accessKeyID
	}), mock.Anything).Return(expected, nil)

	result, err := repo.GetAccessKeyInfo(context.Background(), accessKeyID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetAccessKeyInfo_Error(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expectedErr := errors.New("sts error")
	mockClient.On("GetAccessKeyInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetAccessKeyInfo(context.Background(), "AKIAIOSFODNN7EXAMPLE")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sts GetAccessKeyInfo")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetDelegatedAccessToken(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	tradeInToken := "trade-token"
	expected := &sts.GetDelegatedAccessTokenOutput{}

	mockClient.On("GetDelegatedAccessToken", mock.Anything, mock.MatchedBy(func(input *sts.GetDelegatedAccessTokenInput) bool {
		return input.TradeInToken != nil && *input.TradeInToken == tradeInToken
	}), mock.Anything).Return(expected, nil)

	result, err := repo.GetDelegatedAccessToken(context.Background(), tradeInToken)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetDelegatedAccessToken_Error(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expectedErr := errors.New("sts error")
	mockClient.On("GetDelegatedAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetDelegatedAccessToken(context.Background(), "trade-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sts GetDelegatedAccessToken")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetFederationToken(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	in := &sts.GetFederationTokenInput{Name: aws.String("federated-user")}
	expected := &sts.GetFederationTokenOutput{}

	mockClient.On("GetFederationToken", mock.Anything, in, mock.Anything).Return(expected, nil)

	result, err := repo.GetFederationToken(context.Background(), in)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetFederationToken_NilInput_Error(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expectedErr := errors.New("sts error")
	mockClient.On("GetFederationToken", mock.Anything, mock.MatchedBy(func(input *sts.GetFederationTokenInput) bool {
		return input != nil
	}), mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetFederationToken(context.Background(), nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sts GetFederationToken")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetSessionToken(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	in := &sts.GetSessionTokenInput{DurationSeconds: aws.Int32(3600)}
	expected := &sts.GetSessionTokenOutput{}

	mockClient.On("GetSessionToken", mock.Anything, in, mock.Anything).Return(expected, nil)

	result, err := repo.GetSessionToken(context.Background(), in)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetSessionToken_NilInput_Error(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expectedErr := errors.New("sts error")
	mockClient.On("GetSessionToken", mock.Anything, mock.MatchedBy(func(input *sts.GetSessionTokenInput) bool {
		return input != nil
	}), mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetSessionToken(context.Background(), nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sts GetSessionToken")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetWebIdentityToken(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	in := &sts.GetWebIdentityTokenInput{
		Audience:         []string{"service-a"},
		SigningAlgorithm: aws.String("RS256"),
	}
	expected := &sts.GetWebIdentityTokenOutput{}

	mockClient.On("GetWebIdentityToken", mock.Anything, in, mock.Anything).Return(expected, nil)

	result, err := repo.GetWebIdentityToken(context.Background(), in)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetWebIdentityToken_NilInput_Error(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expectedErr := errors.New("sts error")
	mockClient.On("GetWebIdentityToken", mock.Anything, mock.MatchedBy(func(input *sts.GetWebIdentityTokenInput) bool {
		return input != nil
	}), mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetWebIdentityToken(context.Background(), nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sts GetWebIdentityToken")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSTSRepository_GetCallerIdentity_Error(t *testing.T) {
	mockClient := &MockSTSClient{}
	repo := NewAWSSTSRepositoryWithInterface(mockClient)

	expectedErr := errors.New("sts error")
	mockClient.On("GetCallerIdentity", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetCallerIdentity(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sts GetCallerIdentity")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}
