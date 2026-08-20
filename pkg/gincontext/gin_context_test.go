//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package gincontext

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

var errTestSample = errors.New("test")

func TestGetGinContextError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want    error
		wantErr error
		setup   func(c *gin.Context)
		name    string
	}{
		{
			name: "returns stored error",
			setup: func(c *gin.Context) {
				SetGinContextError(c, errTestSample)
			},
			want: errTestSample,
		},
		{
			name:  "returns nil when no error stored",
			setup: func(c *gin.Context) {},
			want:  nil,
		},
		{
			name: "returns ErrCannotGetError for non-error value",
			setup: func(c *gin.Context) {
				c.Set(contextKeyError, "not an error")
			},
			wantErr: ErrCannotGetError,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &gin.Context{}
			tt.setup(c)

			got, err := GetGinContextError(c)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("GetGinContextError() error = nil, want %v", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetGinContextError() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetGinContextError() unexpected error: %v", err)
			}
			if !errors.Is(tt.want, got) && !errors.Is(tt.want, got) {
				t.Fatalf("GetGinContextError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGinContextErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantErr error
		setup   func(c *gin.Context)
		name    string
		want    string
	}{
		{
			name: "returns stored message",
			setup: func(c *gin.Context) {
				SetGinContextErrorMessage(c, "test")
			},
			want: "test",
		},
		{
			name:  "returns empty string when no message stored",
			setup: func(c *gin.Context) {},
			want:  "",
		},
		{
			name: "returns ErrCannotGetMessage for non-string value",
			setup: func(c *gin.Context) {
				c.Set(contextKeyErrorMessage, 12345)
			},
			wantErr: ErrCannotGetMessage,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &gin.Context{}
			tt.setup(c)

			got, err := GetGinContextErrorMessage(c)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("GetGinContextErrorMessage() error = nil, want %v", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetGinContextErrorMessage() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetGinContextErrorMessage() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("GetGinContextErrorMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetGinContextError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		name string
	}{
		{name: "stores error on context", err: errTestSample},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &gin.Context{}
			SetGinContextError(c, tt.err)

			got, err := GetGinContextError(c)
			if err != nil {
				t.Fatalf("GetGinContextError() unexpected error: %v", err)
			}
			if !errors.Is(got, tt.err) {
				t.Fatalf("GetGinContextError() = %v, want %v", got, tt.err)
			}
		})
	}
}

func TestSetGinContextErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
	}{
		{name: "stores message on context", message: "test"},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &gin.Context{}
			SetGinContextErrorMessage(c, tt.message)

			got, err := GetGinContextErrorMessage(c)
			if err != nil {
				t.Fatalf("GetGinContextErrorMessage() unexpected error: %v", err)
			}
			if got != tt.message {
				t.Fatalf("GetGinContextErrorMessage() = %q, want %q", got, tt.message)
			}
		})
	}
}
