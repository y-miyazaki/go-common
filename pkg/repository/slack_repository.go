package repository

import (
	"fmt"

	"github.com/slack-go/slack"
)

// SlackClientInterface defines the interface for Slack client operations
type SlackClientInterface interface {
	// PostMessage sends a message to a channel
	PostMessage(_ string, _ ...slack.MsgOption) (string, string, error)
}

// SlackRepository struct.
// Note: client field is intentionally unexported to maintain encapsulation.
// Use NewSlackRepositoryWithInterface for testing with mock implementations.
type SlackRepository struct {
	client    SlackClientInterface // Unexported for encapsulation; use constructor for injection
	channelID string
}

// NewSlackRepository returns SlackRepository instance.
func NewSlackRepository(client *slack.Client, channelID string) *SlackRepository {
	return &SlackRepository{
		client:    client,
		channelID: channelID,
	}
}

// NewSlackRepositoryWithInterface returns SlackRepository instance with interface (for testing).
func NewSlackRepositoryWithInterface(client SlackClientInterface, channelID string) *SlackRepository {
	return &SlackRepository{
		client:    client,
		channelID: channelID,
	}
}

// PostMessage sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (r *SlackRepository) PostMessage(options ...slack.MsgOption) error {
	_, _, err := r.client.PostMessage(r.channelID, options...)
	if err != nil {
		return fmt.Errorf("slack PostMessage: %w", err)
	}
	return nil
}

// PostMessageText sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (r *SlackRepository) PostMessageText(text string) error {
	// text
	msgOptText := slack.MsgOptionText(text, true)
	_, _, err := r.client.PostMessage(r.channelID, msgOptText)
	if err != nil {
		return fmt.Errorf("slack PostMessageText: %w", err)
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("slack PostMessageAttachment: %w", err)
	}
	return nil
}
