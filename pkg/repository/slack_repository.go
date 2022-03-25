package repository

import (
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

// SlackRepositoryInterface interface
type SlackRepositoryInterface interface {
	PostMessage(attachment *slack.Attachment) error
}

// SlackRepository struct.
type SlackRepository struct {
	e         *logrus.Entry
	client    *slack.Client
	channelID string
}

// NewSlackRepository returns SlackRepository instance.
func NewSlackRepository(
	e *logrus.Entry,
	client *slack.Client,
	channelID string,
) *SlackRepository {
	return &SlackRepository{
		e:         e,
		client:    client,
		channelID: channelID,
	}
}

// PostMessage sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (r *SlackRepository) PostMessage(text string) error {
	// text
	msgOptText := slack.MsgOptionText(text, true)
	_, _, err := r.client.PostMessage(r.channelID, msgOptText)
	return err
}

// PostMessageAttachment sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (r *SlackRepository) PostMessageAttachment(attachment *slack.Attachment) error {
	// text
	msgOptText := slack.MsgOptionText("", true)
	// attachment
	msgOptAttachments := slack.MsgOptionAttachments(*attachment)
	_, _, err := r.client.PostMessage(r.channelID, msgOptAttachments, msgOptText)
	return err
}