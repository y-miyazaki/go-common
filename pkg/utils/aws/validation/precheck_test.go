package validation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/y-miyazaki/go-common/pkg/utils/aws/validation"
)

func TestCheckAWSCredentials_NilConfig(t *testing.T) {
	t.Parallel()

	_, err := validation.CheckAWSCredentials(context.Background(), nil)
	if !errors.Is(err, validation.ErrNilConfig) {
		t.Fatalf("expected ErrNilConfig, got: %v", err)
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
			name:        "ErrEmptyARN",
			err:         validation.ErrEmptyARN,
			expectedMsg: "aws credentials are not set or invalid: empty ARN",
		},
		{
			name:        "ErrNilConfig",
			err:         validation.ErrNilConfig,
			expectedMsg: "aws config is nil",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedMsg {
				t.Fatalf("expected error message %q, got %q", tt.expectedMsg, tt.err.Error())
			}
		})
	}
}
