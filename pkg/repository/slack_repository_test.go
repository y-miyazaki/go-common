package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSlackClient is a mock implementation of Slack client for testing
type MockSlackClient struct {
	mock.Mock
}

func (m *MockSlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelID, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockSlackClient) PostMessageContext(ctx context.Context, channelID string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(ctx, channelID, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockSlackClient) GetUserInfo(user string) (*slack.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.User), args.Error(1)
}

func (m *MockSlackClient) GetChannelInfo(channel string) (*slack.Channel, error) {
	args := m.Called(channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockSlackClient) GetUserByEmail(email string) (*slack.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.User), args.Error(1)
}

func (m *MockSlackClient) GetUsers() ([]slack.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]slack.User), args.Error(1)
}

func (m *MockSlackClient) GetChannels(private bool, options ...interface{}) ([]slack.Channel, error) {
	args := m.Called(private, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]slack.Channel), args.Error(1)
}

func (m *MockSlackClient) JoinChannel(channelName string) (*slack.Channel, error) {
	args := m.Called(channelName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockSlackClient) LeaveChannel(channelID string) (bool, error) {
	args := m.Called(channelID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSlackClient) ArchiveChannel(channelID string) error {
	args := m.Called(channelID)
	return args.Error(0)
}

func (m *MockSlackClient) UnarchiveChannel(channelID string) error {
	args := m.Called(channelID)
	return args.Error(0)
}

func (m *MockSlackClient) GetConversations(params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]slack.Channel), args.String(1), args.Error(2)
}

func (m *MockSlackClient) GetConversationHistory(params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.GetConversationHistoryResponse), args.Error(1)
}

func (m *MockSlackClient) GetConversationReplies(params *slack.GetConversationRepliesParameters) ([]slack.Message, bool, string, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, false, "", args.Error(3)
	}
	return args.Get(0).([]slack.Message), args.Bool(1), args.String(2), args.Error(3)
}

// SlackRepositoryWithInterface for testing with interface
type SlackRepositoryWithInterface struct {
	client    SlackClientInterface
	channelID string
}

// PostMessage sends a message to a channel.
func (r *SlackRepositoryWithInterface) PostMessage(options ...slack.MsgOption) error {
	_, _, err := r.client.PostMessage(r.channelID, options...)
	return err
}

// PostMessageText sends a message to a channel.
func (r *SlackRepositoryWithInterface) PostMessageText(text string) error {
	_, _, err := r.client.PostMessage(r.channelID, slack.MsgOptionText(text, false))
	return err
}

// PostMessageAttachment sends a message with attachment to a channel.
func (r *SlackRepositoryWithInterface) PostMessageAttachment(attachment *slack.Attachment) error {
	_, _, err := r.client.PostMessage(r.channelID, slack.MsgOptionAttachments(*attachment))
	return err
}

func TestNewSlackRepository(t *testing.T) {
	mockClient := &slack.Client{}
	channelID := "test_channel"
	repo := NewSlackRepository(mockClient, channelID)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.client)
	assert.Equal(t, channelID, repo.channelID)
}

func TestSlackRepository_PostMessageText(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	testMessage := "Hello, World!"
	mockClient.On("PostMessage", channelID, mock.MatchedBy(func(options []slack.MsgOption) bool {
		// Check if options array is not empty (simplified check)
		return len(options) > 0
	})).Return("timestamp", "message_id", nil)

	err := repo.PostMessageText(testMessage)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessageText_Error(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	testMessage := "Hello, World!"
	expectedError := errors.New("slack API error")

	mockClient.On("PostMessage", channelID, mock.Anything).Return("", "", expectedError)

	err := repo.PostMessageText(testMessage)

	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("slack PostMessageText: %w", expectedError), err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessageAttachment(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	attachment := &slack.Attachment{
		Text:  "Test attachment",
		Color: "good",
	}

	mockClient.On("PostMessage", channelID, mock.Anything).Return("timestamp", "message_id", nil)

	err := repo.PostMessageAttachment(attachment)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessageAttachment_Error(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	attachment := &slack.Attachment{
		Text:  "Test attachment",
		Color: "good",
	}

	expectedError := errors.New("slack API error")

	mockClient.On("PostMessage", channelID, mock.Anything).Return("", "", expectedError)

	err := repo.PostMessageAttachment(attachment)

	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("slack PostMessageAttachment: %w", expectedError), err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessage_Error(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	options := []slack.MsgOption{
		slack.MsgOptionText("Test message", false),
		slack.MsgOptionUsername("TestBot"),
	}

	expectedError := errors.New("slack API error")

	mockClient.On("PostMessage", channelID, options).Return("", "", expectedError)

	err := repo.PostMessage(options...)

	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("slack PostMessage: %w", expectedError), err)
	mockClient.AssertExpectations(t)
}

// Test actual SlackRepository functions
func TestSlackRepositoryReal_PostMessageText(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Slack integration test - requires real Slack API token")
}

func TestSlackRepositoryReal_PostMessageAttachment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Slack integration test - requires real Slack API token")
}

func TestSlackRepositoryReal_PostMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Slack integration test - requires real Slack API token")
}

// Test edge cases and validation
func TestSlackRepository_PostMessageText_EmptyMessage(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	emptyMessage := ""
	mockClient.On("PostMessage", channelID, mock.Anything).Return("timestamp", "message_id", nil)

	err := repo.PostMessageText(emptyMessage)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessageText_LongMessage(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	// Create a long message (4000 characters)
	longMessage := ""
	for i := 0; i < 400; i++ {
		longMessage += "This is a test message that is very long. "
	}

	mockClient.On("PostMessage", channelID, mock.Anything).Return("timestamp", "message_id", nil)

	err := repo.PostMessageText(longMessage)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessageAttachment_EmptyAttachment(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	emptyAttachment := &slack.Attachment{}

	mockClient.On("PostMessage", channelID, mock.Anything).Return("timestamp", "message_id", nil)

	err := repo.PostMessageAttachment(emptyAttachment)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessageAttachment_WithFields(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	attachment := &slack.Attachment{
		Text:  "Test attachment with fields",
		Color: "good",
		Fields: []slack.AttachmentField{
			{
				Title: "Field 1",
				Value: "Value 1",
				Short: true,
			},
			{
				Title: "Field 2",
				Value: "Value 2",
				Short: true,
			},
		},
	}

	mockClient.On("PostMessage", channelID, mock.Anything).Return("timestamp", "message_id", nil)

	err := repo.PostMessageAttachment(attachment)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessage_MultipleOptions(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	options := []slack.MsgOption{
		slack.MsgOptionText("Test message", false),
		slack.MsgOptionUsername("TestBot"),
		slack.MsgOptionIconEmoji(":robot_face:"),
		slack.MsgOptionAsUser(false),
	}

	mockClient.On("PostMessage", channelID, options).Return("timestamp", "message_id", nil)

	err := repo.PostMessage(options...)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSlackRepository_PostMessage_EmptyOptions(t *testing.T) {
	mockClient := &MockSlackClient{}
	channelID := "test_channel"
	repo := NewSlackRepositoryWithInterface(mockClient, channelID)

	emptyOptions := []slack.MsgOption{}

	mockClient.On("PostMessage", channelID, emptyOptions).Return("timestamp", "message_id", nil)

	err := repo.PostMessage(emptyOptions...)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
