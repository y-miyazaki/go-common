//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package utils

import (
	"io"
	"strings"
	"testing"
)

func TestGetStringCount(t *testing.T) {
	t.Parallel()
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
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetStringCount(tt.args.str); got != tt.want {
				t.Fatalf("GetStringCount(%q) = %v, want %v", tt.args.str, got, tt.want)
			}
		})
	}
}

func TestCheckStringCount(t *testing.T) {
	t.Parallel()
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
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := CheckStringCount(tt.args.str, tt.args.maxLen); got != tt.want {
				t.Fatalf("CheckStringCount(%q, %d) = %v, want %v", tt.args.str, tt.args.maxLen, got, tt.want)
			}
		})
	}
}

func TestSliceUTF8(t *testing.T) {
	t.Parallel()
	type args struct {
		str string
		pos int
	}
	tests := []struct {
		name string
		want string
		args args
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
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := SliceUTF8(tt.args.str, tt.args.pos); got != tt.want {
				t.Fatalf("SliceUTF8(%q, %d) = %v, want %v", tt.args.str, tt.args.pos, got, tt.want)
			}
		})
	}
}

func TestSliceUTF8AddString(t *testing.T) {
	t.Parallel()
	type args struct {
		str       string
		addString string
		pos       int
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
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := SliceUTF8AddString(tt.args.str, tt.args.pos, tt.args.addString); got != tt.want {
				t.Fatalf("SliceUTF8AddString(%q, %d, %q) = %v, want %v", tt.args.str, tt.args.pos, tt.args.addString, got, tt.want)
			}
		})
	}
}

func TestConvertToStringaa(t *testing.T) {
	t.Parallel()
	type args struct {
		input any
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
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ConvertToString(tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ConvertToString() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ConvertToString() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("ConvertToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringFromReadCloser(t *testing.T) {
	t.Parallel()
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
				r: io.NopCloser(strings.NewReader("Hello, world!")),
			},
			want:    "Hello, world!",
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				r: &errReader{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetStringFromReadCloser(tt.args.r)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetStringFromReadCloser() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetStringFromReadCloser() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("GetStringFromReadCloser() = %q, want %q", got, tt.want)
			}
		})
	}
}
