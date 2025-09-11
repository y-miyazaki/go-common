// Package repository provides repository implementations for various AWS services and databases.
package repository

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecretsManagerClientInterface interface for AWS Secrets Manager operations
type AWSSecretsManagerClientInterface interface {
	// GetSecretValue retrieves the value of a secret from AWS Secrets Manager
	GetSecretValue(_ context.Context, _ *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// AWSSecretsManagerRepository struct.
type AWSSecretsManagerRepository struct {
	c AWSSecretsManagerClientInterface
}

// NewAWSSecretsManagerRepository returns AWSSecretsManagerRepository instance.
func NewAWSSecretsManagerRepository(c *secretsmanager.Client) *AWSSecretsManagerRepository {
	return &AWSSecretsManagerRepository{
		c: c,
	}
}

// GetSecretString gets a secret string from secretsmanager.
func (r *AWSSecretsManagerRepository) GetSecretString(ctx context.Context, secretName string) (string, error) {
	result, err := r.c.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		// VersionStage defaults to AWSCURRENT if unspecified
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return "", fmt.Errorf("secretsmanager GetSecretValue: %w", err)
	}
	if result.SecretString != nil { // pragma: allowlist-secret
		return *result.SecretString, nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	length, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return "", fmt.Errorf("secretsmanager decode secret binary: %w", err)
	}
	return string(decodedBinarySecretBytes[:length]), nil
}
