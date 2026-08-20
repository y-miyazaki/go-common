//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package utils

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

var errTestRead = errors.New("read error")

// errReader is an io.ReadCloser that always returns an error on Read.
type errReader struct{}

func (e *errReader) Read(_ []byte) (int, error) {
	return 0, errTestRead
}

func (e *errReader) Close() error { return nil }

func TestGetBufferFromReadCloser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       io.ReadCloser
		want    []byte
		wantErr bool
	}{
		{
			name:    "reads all bytes from reader",
			r:       io.NopCloser(strings.NewReader("Hello, world!")),
			want:    []byte("Hello, world!"),
			wantErr: false,
		},
		{
			name:    "returns error when read fails",
			r:       &errReader{},
			want:    nil,
			wantErr: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetBufferFromReadCloser(tt.r)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetBufferFromReadCloser() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetBufferFromReadCloser() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetBufferFromReadCloser() = %v, want %v", got, tt.want)
			}
		})
	}
}
