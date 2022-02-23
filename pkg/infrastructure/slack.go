package infrastructure

import (
	"github.com/slack-go/slack"
)

// SlackConfigSetting sets configurations.
type SlackConfigSetting struct {
	OauthAccessToken string
}

// NewSlack returns an SlackClient instance.
func NewSlack(c *SlackConfigSetting) *slack.Client {
	return slack.New(c.OauthAccessToken)
}
