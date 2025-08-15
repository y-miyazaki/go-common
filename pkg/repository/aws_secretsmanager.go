package repository

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecretsManagerRepositoryInterface interface.
// nolint:iface,revive,unused
type AWSSecretsManagerRepositoryInterface interface {
	GetSecretString(secretName string) (string, error)
	GetCacheSecretString(secretName string) (string, error)
}

// AWSSecretsManagerRepository struct.
type AWSSecretsManagerRepository struct {
	c *secretsmanager.Client
}

// NewAWSSecretsManagerRepository returns AWSSecretsManagerRepository instance.
func NewAWSSecretsManagerRepository(c *secretsmanager.Client) *AWSSecretsManagerRepository {
	return &AWSSecretsManagerRepository{
		c: c,
	}
}

// GetSecretString gets a secret string from secretsmanager.
func (r *AWSSecretsManagerRepository) GetSecretString(secretName string) (string, error) {
	result, err := r.c.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		// VersionStage defaults to AWSCURRENT if unspecified
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return "", fmt.Errorf("secretsmanager GetSecretValue: %w", err)
	}
	if result.SecretString != nil { // pragma: allowlist secret
		return *result.SecretString, nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	length, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return "", fmt.Errorf("secretsmanager decode secret binary: %w", err)
	}
	return string(decodedBinarySecretBytes[:length]), nil
}
