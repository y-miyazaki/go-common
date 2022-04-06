package repository

import (
	"github.com/slack-go/slack"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// SlackRepositoryInterface interface
type SlackRepositoryInterface interface {
	PostMessage(attachment *slack.Attachment) error
	PostMessageAttachment(attachment *slack.Attachment) error
}

// SlackRepository struct.
type SlackRepository struct {
	logger    *logger.Logger
	client    *slack.Client
	channelID string
}

// NewSlackRepository returns SlackRepository instance.
func NewSlackRepository(logger *logger.Logger, client *slack.Client, channelID string) *SlackRepository {
	return &SlackRepository{
		logger:    logger,
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
