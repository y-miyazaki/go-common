// Package secretsmanager provides utilities for AWS Secrets Manager operations.
package secretsmanager

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// GetSecretString gets a secret string from secretsmanager.
func GetSecretString(secretName, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %w", err)
	}
	// Create a Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		// VersionStage defaults to AWSCURRENT if unspecified
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(context.Background(), input)
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
