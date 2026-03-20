// Package infrastructure provides AWS Account integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
)

// NewAWSAccount returns AWS Account client using AWS SDK v2.
func NewAWSAccount(c *aws.Config, optFns ...func(*account.Options)) *account.Client { // nolint:gocritic
	return account.NewFromConfig(*c, optFns...)
}
