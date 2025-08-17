// Package infrastructure provides AWS Secrets Manager integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// NewAWSSecretsManager returns Secrets Manager client using AWS SDK v2.
func NewAWSSecretsManager(c *aws.Config, optFns ...func(*secretsmanager.Options)) *secretsmanager.Client { // nolint:gocritic
	return secretsmanager.NewFromConfig(*c, optFns...)
}
