package utils

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

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
