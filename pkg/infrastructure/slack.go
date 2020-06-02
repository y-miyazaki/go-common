package infrastructure

import (
	"github.com/nlopes/slack"
)

// SlackClient sets slack instance.
type SlackClient struct {
	API       *slack.Client
	ChannelID string
}

// NewSlack returns an instance of logger
func NewSlack(oauthAccessToken, channelID string) *SlackClient {
	slackClient := &SlackClient{}
	slackClient.API = slack.New(oauthAccessToken)
	slackClient.ChannelID = channelID
	return slackClient
}

// PostMessage sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (s *SlackClient) PostMessage(attachment *slack.Attachment) error {
	// text
	msgOptText := slack.MsgOptionText("", true)
	// attachment
	msgOptAttachments := slack.MsgOptionAttachments(*attachment)
	_, _, err := s.API.PostMessage(s.ChannelID, msgOptAttachments, msgOptText)
	return err
}
