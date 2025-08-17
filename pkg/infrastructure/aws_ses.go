// Package infrastructure provides AWS SES integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

// NewAWSSES returns SES client using AWS SDK v2.
func NewAWSSES(c aws.Config, optFns ...func(*sesv2.Options)) *sesv2.Client { // nolint:gocritic
	return sesv2.NewFromConfig(c, optFns...)
}
