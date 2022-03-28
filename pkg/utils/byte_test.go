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
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{
				r: ioutil.NopCloser(strings.NewReader("Hello, world!")),
			},
			want: []byte("Hello, world!"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBufferFromReadCloser(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBufferFromReadCloser() = %v, want %v", got, tt.want)
			}
		})
	}
}
