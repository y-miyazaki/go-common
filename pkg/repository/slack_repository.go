package repository

import (
	"github.com/slack-go/slack"
)

// SlackRepositoryInterface interface
type SlackRepositoryInterface interface {
	PostMessage(attachment *slack.Attachment) error
	PostMessageAttachment(attachment *slack.Attachment) error
}

// SlackRepository struct.
type SlackRepository struct {
	client    *slack.Client
	channelID string
}

// NewSlackRepository returns SlackRepository instance.
func NewSlackRepository(client *slack.Client, channelID string) *SlackRepository {
	return &SlackRepository{
		client:    client,
		channelID: channelID,
	}
}

// PostMessageAttachment sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (r *SlackRepository) PostMessage(options ...slack.MsgOption) error {
	_, _, err := r.client.PostMessage(r.channelID, options...)
	return err
}

// PostMessageText sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (r *SlackRepository) PostMessageText(text string) error {
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
