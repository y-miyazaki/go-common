package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// AWSSTSClientInterface defines the interface for STS client operations.
type AWSSTSClientInterface interface {
	// GetAccessKeyInfo returns the account identifier for the specified access key.
	GetAccessKeyInfo(_ context.Context, _ *sts.GetAccessKeyInfoInput, _ ...func(*sts.Options)) (*sts.GetAccessKeyInfoOutput, error)
	// GetCallerIdentity returns details about the IAM identity whose credentials are used.
	GetCallerIdentity(_ context.Context, _ *sts.GetCallerIdentityInput, _ ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
	// GetDelegatedAccessToken requests a delegated access token from STS.
	GetDelegatedAccessToken(_ context.Context, _ *sts.GetDelegatedAccessTokenInput, _ ...func(*sts.Options)) (*sts.GetDelegatedAccessTokenOutput, error)
	// GetFederationToken requests temporary security credentials for federation.
	GetFederationToken(_ context.Context, _ *sts.GetFederationTokenInput, _ ...func(*sts.Options)) (*sts.GetFederationTokenOutput, error)
	// GetSessionToken requests temporary security credentials for an IAM user.
	GetSessionToken(_ context.Context, _ *sts.GetSessionTokenInput, _ ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error)
	// GetWebIdentityToken requests a web identity token from STS.
	GetWebIdentityToken(_ context.Context, _ *sts.GetWebIdentityTokenInput, _ ...func(*sts.Options)) (*sts.GetWebIdentityTokenOutput, error)
}

// AWSSTSRepository struct.
type AWSSTSRepository struct {
	Client AWSSTSClientInterface
}

// NewAWSSTSRepository returns AWSSTSRepository instance.
func NewAWSSTSRepository(c *sts.Client) *AWSSTSRepository {
	return &AWSSTSRepository{Client: c}
}

// NewAWSSTSRepositoryWithInterface returns AWSSTSRepository instance with interface (for testing).
func NewAWSSTSRepositoryWithInterface(c AWSSTSClientInterface) *AWSSTSRepository {
	return &AWSSTSRepository{Client: c}
}

// GetAccessKeyInfo returns information about the specified access key.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetAccessKeyInfo.html
func (r *AWSSTSRepository) GetAccessKeyInfo(ctx context.Context, accessKeyID string) (*sts.GetAccessKeyInfoOutput, error) {
	result, err := r.Client.GetAccessKeyInfo(ctx, &sts.GetAccessKeyInfoInput{
		AccessKeyId: aws.String(accessKeyID),
	})
	if err != nil {
		return nil, fmt.Errorf("sts GetAccessKeyInfo: %w", err)
	}

	return result, nil
}

// GetCallerIdentity returns details about the current caller identity.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html
func (r *AWSSTSRepository) GetCallerIdentity(ctx context.Context) (*sts.GetCallerIdentityOutput, error) {
	result, err := r.Client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("sts GetCallerIdentity: %w", err)
	}

	return result, nil
}

// GetDelegatedAccessToken requests a delegated access token.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetDelegatedAccessToken.html
func (r *AWSSTSRepository) GetDelegatedAccessToken(ctx context.Context, tradeInToken string) (*sts.GetDelegatedAccessTokenOutput, error) {
	result, err := r.Client.GetDelegatedAccessToken(ctx, &sts.GetDelegatedAccessTokenInput{
		TradeInToken: aws.String(tradeInToken),
	})
	if err != nil {
		return nil, fmt.Errorf("sts GetDelegatedAccessToken: %w", err)
	}

	return result, nil
}

// GetFederationToken requests temporary credentials for a federated user.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html
func (r *AWSSTSRepository) GetFederationToken(ctx context.Context, in *sts.GetFederationTokenInput) (*sts.GetFederationTokenOutput, error) {
	req := in
	if req == nil {
		req = &sts.GetFederationTokenInput{}
	}

	result, err := r.Client.GetFederationToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sts GetFederationToken: %w", err)
	}

	return result, nil
}

// GetSessionToken requests temporary credentials for an IAM user.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html
func (r *AWSSTSRepository) GetSessionToken(ctx context.Context, in *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
	req := in
	if req == nil {
		req = &sts.GetSessionTokenInput{}
	}

	result, err := r.Client.GetSessionToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sts GetSessionToken: %w", err)
	}

	return result, nil
}

// GetWebIdentityToken requests a web identity token.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetWebIdentityToken.html
func (r *AWSSTSRepository) GetWebIdentityToken(ctx context.Context, in *sts.GetWebIdentityTokenInput) (*sts.GetWebIdentityTokenOutput, error) {
	req := in
	if req == nil {
		req = &sts.GetWebIdentityTokenInput{}
	}

	result, err := r.Client.GetWebIdentityToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sts GetWebIdentityToken: %w", err)
	}

	return result, nil
}
