package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountClient is a mock implementation of the AWS Account client.
type MockAccountClient struct {
	mock.Mock
}

func (m *MockAccountClient) GetAccountInformation(ctx context.Context, params *account.GetAccountInformationInput, optFns ...func(*account.Options)) (*account.GetAccountInformationOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.GetAccountInformationOutput), args.Error(1)
}

func TestNewAWSAccountRepository(t *testing.T) {
	mockClient := &account.Client{}
	repo := NewAWSAccountRepository(mockClient)

	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.Client)
}

func TestAWSAccountRepository_GetAccountInformation_WithAccountID(t *testing.T) {
	mockClient := &MockAccountClient{}
	repo := NewAWSAccountRepositoryWithInterface(mockClient)

	accountID := "123456789012"
	expected := &account.GetAccountInformationOutput{
		AccountId:   aws.String(accountID),
		AccountName: aws.String("sample-account"),
	}

	mockClient.On("GetAccountInformation", mock.Anything, mock.MatchedBy(func(input *account.GetAccountInformationInput) bool {
		return input.AccountId != nil && *input.AccountId == accountID
	}), mock.Anything).Return(expected, nil)

	result, err := repo.GetAccountInformation(context.Background(), accountID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSAccountRepository_GetAccountInformation_WithoutAccountID(t *testing.T) {
	mockClient := &MockAccountClient{}
	repo := NewAWSAccountRepositoryWithInterface(mockClient)

	expected := &account.GetAccountInformationOutput{}

	mockClient.On("GetAccountInformation", mock.Anything, mock.MatchedBy(func(input *account.GetAccountInformationInput) bool {
		return input.AccountId == nil
	}), mock.Anything).Return(expected, nil)

	result, err := repo.GetAccountInformation(context.Background(), "")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestAWSAccountRepository_GetAccountInformation_Error(t *testing.T) {
	mockClient := &MockAccountClient{}
	repo := NewAWSAccountRepositoryWithInterface(mockClient)

	expectedErr := errors.New("account service error")

	mockClient.On("GetAccountInformation", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	result, err := repo.GetAccountInformation(context.Background(), "123456789012")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account GetAccountInformation")
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}
