package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
)

// AWSAccountClientInterface defines the interface for AWS Account client operations.
type AWSAccountClientInterface interface {
	// GetAccountInformation retrieves information about an AWS account.
	GetAccountInformation(_ context.Context, _ *account.GetAccountInformationInput, _ ...func(*account.Options)) (*account.GetAccountInformationOutput, error)
}

// AWSAccountRepository struct.
type AWSAccountRepository struct {
	Client AWSAccountClientInterface
}

// NewAWSAccountRepository returns AWSAccountRepository instance.
func NewAWSAccountRepository(c *account.Client) *AWSAccountRepository {
	return &AWSAccountRepository{Client: c}
}

// NewAWSAccountRepositoryWithInterface returns AWSAccountRepository instance with interface (for testing).
func NewAWSAccountRepositoryWithInterface(c AWSAccountClientInterface) *AWSAccountRepository {
	return &AWSAccountRepository{Client: c}
}

// GetAccountInformation retrieves account information from AWS Account Management.
// https://docs.aws.amazon.com/accounts/latest/reference/API_GetAccountInformation.html
func (r *AWSAccountRepository) GetAccountInformation(ctx context.Context, accountID string) (*account.GetAccountInformationOutput, error) {
	in := &account.GetAccountInformationInput{}
	if accountID != "" {
		in.AccountId = aws.String(accountID)
	}

	result, err := r.Client.GetAccountInformation(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("account GetAccountInformation: %w", err)
	}

	return result, nil
}
