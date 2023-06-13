package repository

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type AWSSecretsManagerRepositoryInterface interface {
	GetSecretString(secretName string) (string, error)
}

type AWSSecretsManagerRepository struct {
	sm      *secretsmanager.SecretsManager
	cache   *secretcache.Cache
	session *session.Session
}

// NewAWSSecretsManagerRepository returns NewAWSSecretsManagerRepository instance.
func NewAWSSecretsManagerRepository(sm *secretsmanager.SecretsManager, sess *session.Session, cache *secretcache.Cache) *AWSSecretsManagerRepository {
	return &AWSSecretsManagerRepository{
		sm:      sm,
		cache:   cache,
		session: sess,
	}
}

// GetSecretString gets a secret string from secretsmanager.
func (r *AWSSecretsManagerRepository) GetSecretString(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		// VersionStage defaults to AWSCURRENT if unspecified
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := r.sm.GetSecretValue(input)
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
func (r *AWSSecretsManagerRepository) GetCacheSecretString(secretName string) (string, error) {
	return r.cache.GetSecretString(secretName)
}
