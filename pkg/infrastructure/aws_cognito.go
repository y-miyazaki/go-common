// Package infrastructure provides AWS Cognito integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// NewAWSCognito returns Cognito Identity Provider client using AWS SDK v2.
func NewAWSCognito(c *aws.Config, optFns ...func(*cognitoidentityprovider.Options)) *cognitoidentityprovider.Client { // nolint:gocritic
	return cognitoidentityprovider.NewFromConfig(*c, optFns...)
}
