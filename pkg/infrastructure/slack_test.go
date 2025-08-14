package infrastructure

import (
	"testing"

	"github.com/slack-go/slack"
)

func TestNewSlack(t *testing.T) {
	type args struct {
		c *SlackConfig
	}
	c := &SlackConfig{
		OauthAccessToken: "test",
	}
	tests := []struct {
		name string
		args args
		want *slack.Client
	}{
		{
			name: "test1",
			args: args{
				c: c,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlack(tt.args.c); got == tt.want {
				t.Errorf("NewSlack() = %v, want %v", got, tt.want)
			}
		})
	}
}
