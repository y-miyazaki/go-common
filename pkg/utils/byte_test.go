package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

// errReader is an io.ReadCloser that always returns an error on Read.
type errReader struct{}

func (e *errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errReader) Close() error { return nil }

func TestGetBufferFromReadCloser(t *testing.T) {
	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				r: ioutil.NopCloser(strings.NewReader("Hello, world!")),
			},
			want:    []byte("Hello, world!"),
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				r: &errReader{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBufferFromReadCloser(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBufferFromReadCloser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBufferFromReadCloser() = %v, want %v", got, tt.want)
			}
		})
	}
}
