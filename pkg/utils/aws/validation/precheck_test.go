//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package validation_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/y-miyazaki/go-common/pkg/utils/aws/validation"
)

func TestCheckAWSCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cfg        *aws.Config
		wantErr    error
		wantSubstr string
	}{
		{
			name:    "nil config returns ErrNilConfig",
			cfg:     nil,
			wantErr: validation.ErrNilConfig,
		},
		{
			name:       "invalid credentials return wrapped error",
			cfg:        &aws.Config{Region: "us-east-1"},
			wantSubstr: "aws credentials are not set or invalid",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := validation.CheckAWSCredentials(context.Background(), tt.cfg)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("CheckAWSCredentials() error = nil, want %v", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("CheckAWSCredentials() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err == nil {
				t.Fatalf("CheckAWSCredentials() error = nil, want error containing %q", tt.wantSubstr)
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Fatalf("CheckAWSCredentials() error = %v, want substring %q", err, tt.wantSubstr)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		expectedMsg string
	}{
		{
			name:        "ErrEmptyARN message",
			err:         validation.ErrEmptyARN,
			expectedMsg: "aws credentials are not set or invalid: empty ARN",
		},
		{
			name:        "ErrNilConfig message",
			err:         validation.ErrNilConfig,
			expectedMsg: "aws config is nil",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedMsg {
				t.Fatalf("error message = %q, want %q", tt.err.Error(), tt.expectedMsg)
			}
		})
	}
}
