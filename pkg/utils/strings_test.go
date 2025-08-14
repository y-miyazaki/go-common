package utils

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetStringCount(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "normal1",
			args: args{
				str: "あいうえお",
			},
			want: 5,
		},
		{
			name: "normal2",
			args: args{
				str: "あいtesttest1うえお",
			},
			want: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringCount(tt.args.str); got != tt.want {
				t.Errorf("GetStringCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckStringCount(t *testing.T) {
	type args struct {
		str    string
		maxLen int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal1",
			args: args{
				str:    "あいうえお",
				maxLen: 6,
			},
			want: true,
		},
		{
			name: "normal2",
			args: args{
				str:    "あいうえお",
				maxLen: 5,
			},
			want: true,
		},
		{
			name: "normal3",
			args: args{
				str:    "あいうえお",
				maxLen: 4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckStringCount(tt.args.str, tt.args.maxLen); got != tt.want {
				t.Errorf("CheckStringCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceUTF8(t *testing.T) {
	type args struct {
		str string
		pos int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal1",
			args: args{
				str: "あいうえお",
				pos: 1,
			},
			want: "あ",
		},
		{
			name: "normal2",
			args: args{
				str: "あいうえお",
				pos: 2,
			},
			want: "あい",
		},
		{
			name: "normal3",
			args: args{
				str: "あいうえお",
				pos: 5,
			},
			want: "あいうえお",
		},
		{
			name: "normal4",
			args: args{
				str: "あいうえお",
				pos: 6,
			},
			want: "あいうえお",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceUTF8(tt.args.str, tt.args.pos); got != tt.want {
				t.Errorf("SliceUTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceUTF8AddString(t *testing.T) {
	type args struct {
		str       string
		pos       int
		addString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal1",
			args: args{
				str:       "あいうえお",
				pos:       1,
				addString: "...",
			},
			want: "あ...",
		},
		{
			name: "normal2",
			args: args{
				str:       "あいうえお",
				pos:       2,
				addString: "...",
			},
			want: "あい...",
		},
		{
			name: "normal3",
			args: args{
				str:       "あいうえお",
				pos:       5,
				addString: "...",
			},
			want: "あいうえお",
		},
		{
			name: "normal4",
			args: args{
				str:       "あいうえお",
				pos:       6,
				addString: "...",
			},
			want: "あいうえお",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceUTF8AddString(tt.args.str, tt.args.pos, tt.args.addString); got != tt.want {
				t.Errorf("SliceUTF8AddString() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestConvertToStringaa(t *testing.T) {
	type args struct {
		input interface{}
	}
	testArgs := args{}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal-int",
			args: args{
				input: 123,
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "normal-bool1",
			args: args{
				input: true,
			},
			want:    "true",
			wantErr: false,
		},
		{
			name: "normal-bool2",
			args: args{
				input: false,
			},
			want:    "false",
			wantErr: false,
		},
		{
			name: "normal-float32",
			args: args{
				input: float32(123),
			},
			want:    "123.000000",
			wantErr: false,
		},
		{
			name: "normal-float64",
			args: args{
				input: float64(123),
			},
			want:    "123.000000",
			wantErr: false,
		},
		{
			name: "normal-string",
			args: args{
				input: "あいうえお",
			},
			want:    "あいうえお",
			wantErr: false,
		},
		{
			name: "normal-string",
			args: args{
				input: testArgs,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringFromReadCloser(t *testing.T) {
	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				r: ioutil.NopCloser(strings.NewReader("Hello, world!")),
			},
			want:    "Hello, world!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStringFromReadCloser(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStringFromReadCloser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStringFromReadCloser() = %v, want %v", got, tt.want)
			}
		})
	}
}
