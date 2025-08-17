// Package infrastructure provides AWS CloudWatch Logs integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// NewAWSCloudWatchLogs returns CloudWatch Logs client using AWS SDK v2.
func NewAWSCloudWatchLogs(c *aws.Config, optFns ...func(*cloudwatchlogs.Options)) *cloudwatchlogs.Client { // nolint:gocritic
	return cloudwatchlogs.NewFromConfig(*c, optFns...)
}
