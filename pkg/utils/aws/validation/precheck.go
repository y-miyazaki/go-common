// Package validation provides pre-execution validation utilities.
package validation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	// ErrEmptyARN is returned when the caller identity ARN is empty.
	ErrEmptyARN = errors.New("aws credentials are not set or invalid: empty ARN")
	// ErrNilConfig is returned when the AWS config is nil.
	ErrNilConfig = errors.New("aws config is nil")
)

// CheckAWSCredentials validates AWS credentials by calling STS GetCallerIdentity.
// It returns the caller identity ARN on success.
func CheckAWSCredentials(ctx context.Context, cfg *aws.Config) (string, error) {
	if cfg == nil {
		return "", ErrNilConfig
	}

	svc := sts.NewFromConfig(*cfg)
	result, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		if strings.Contains(err.Error(), "sso session has expired") || strings.Contains(err.Error(), "SSO session has expired") {
			return "", fmt.Errorf("aws sso session has expired. please run 'aws sso login' to refresh your session: %w", err)
		}
		return "", fmt.Errorf("aws credentials are not set or invalid: %w", err)
	}

	arn := aws.ToString(result.Arn)
	if arn == "" {
		return "", ErrEmptyARN
	}

	return arn, nil
}
