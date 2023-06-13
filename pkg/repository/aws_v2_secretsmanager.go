package repository

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type AWSV2SecretsManagerRepositoryInterface interface {
	GetSecretString(secretName string) (string, error)
}

type AWSV2SecretsManagerRepository struct {
	c     *secretsmanager.Client
	cache *secretcache.Cache
}

// NewAWSV2SecretsManagerRepository returns NewAWSV2SecretsManagerRepository instance.
func NewAWSV2SecretsManagerRepository(c *secretsmanager.Client, cache *secretcache.Cache) *AWSV2SecretsManagerRepository {
	return &AWSV2SecretsManagerRepository{
		c:     c,
		cache: cache,
	}
}

// GetSecretString gets a secret string from secretsmanager.
func (r *AWSV2SecretsManagerRepository) GetSecretString(secretName string) (string, error) {
	result, err := r.c.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		// VersionStage defaults to AWSCURRENT if unspecified
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return "", err
	}
	if result.SecretString != nil { // pragma: allowlist secret
		return *result.SecretString, nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	length, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return "", err
	}
	return string(decodedBinarySecretBytes[:length]), nil
}

// GetCacheSecretString gets a cache secret string from secretsmanager.
func (r *AWSV2SecretsManagerRepository) GetCacheSecretString(secretName string) (string, error) {
	return r.cache.GetSecretString(secretName)
}
