package infrastructure

import (
	"github.com/slack-go/slack"
)

// SlackConfig sets configurations.
type SlackConfig struct {
	OauthAccessToken string
}

// NewSlack returns an SlackClient instance.
func NewSlack(c *SlackConfig) *slack.Client {
	return slack.New(c.OauthAccessToken)
}
