package repository

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsManagerRepositoryInterface interface {
	GetSecretString(secretName string) (string, error)
}

type SecretsManagerRepository struct {
	sm      *secretsmanager.SecretsManager
	session *session.Session
}

func NewSecretsManagerRepository(sm *secretsmanager.SecretsManager, sess *session.Session) *SecretsManagerRepository {
	return &SecretsManagerRepository{
		sm:      sm,
		session: sess,
	}
}

// GetSecretString gets a secret string from secretsmanager.
func (r *SecretsManagerRepository) GetSecretString(secretName string) (string, error) {
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
