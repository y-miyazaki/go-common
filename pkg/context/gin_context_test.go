package context

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetGinContextError(t *testing.T) {
	c := &gin.Context{}
	errData := errors.New("test")
	SetGinContextError(c, errData)

	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				c: c,
			},
			want:    errData,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGinContextError(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGinContextError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !errors.Is(tt.want, got) && tt.want != got {
				t.Errorf("GetGinContextError() = %v, want %v", got, tt.want)
			}
		})
	}

	c = &gin.Context{}
	tests = []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "test2",
			args: args{
				c: c,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGinContextError(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGinContextError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !errors.Is(tt.want, got) && tt.want != got {
				t.Errorf("GetGinContextError() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestGetGinContextErrorMessage(t *testing.T) {
	c := &gin.Context{}
	SetGinContextErrorMessage(c, "test")
	type args struct {
		c *gin.Context
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
				c: c,
			},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGinContextErrorMessage(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGinContextErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGinContextErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
	c = &gin.Context{}
	tests = []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				c: c,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGinContextErrorMessage(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGinContextErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGinContextErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetGinContextError(t *testing.T) {
	c := &gin.Context{}
	type args struct {
		c   *gin.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				c:   c,
				err: errors.New("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGinContextError(tt.args.c, tt.args.err)
		})
	}
}
func TestSetGinContextErrorMessage(t *testing.T) {
	c := &gin.Context{}
	type args struct {
		c       *gin.Context
		message interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				c:       c,
				message: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGinContextErrorMessage(tt.args.c, tt.args.message)
		})
	}
}
