// Package infrastructure provides AWS STS integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// NewAWSSTS returns STS client using AWS SDK v2.
func NewAWSSTS(c *aws.Config, optFns ...func(*sts.Options)) *sts.Client { // nolint:gocritic
	return sts.NewFromConfig(*c, optFns...)
}
