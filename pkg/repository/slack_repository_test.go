package repository

import (
	"errors"
	"strings"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestSlack = errors.New("slack API error")

func TestNewSlackRepository(t *testing.T) {
	t.Parallel()

	mockClient := &slack.Client{}
	channelID := "test_channel"
	repo := NewSlackRepository(mockClient, channelID)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
	require.Equal(t, channelID, repo.channelID)
}

func TestSlackRepository_PostMessageText(t *testing.T) {
	t.Parallel()

	channelID := "test_channel"

	tests := []struct {
		name       string
		message    string
		setupMock  func(m *mocks.MockSlackClientInterface)
		wantErrMsg string
	}{
		{
			name:    "success",
			message: "Hello, World!",
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
		{
			name:    "client error",
			message: "Hello, World!",
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any()).
					Return("", "", errTestSlack)
			},
			wantErrMsg: "slack PostMessageText",
		},
		{
			name:    "empty message",
			message: "",
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
		{
			name:    "long message",
			message: strings.Repeat("This is a test message that is very long. ", 400),
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockSlackClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewSlackRepositoryWithInterface(mockClient, channelID)
			err := repo.PostMessageText(tc.message)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Contains(t, err.Error(), errTestSlack.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSlackRepository_PostMessageAttachment(t *testing.T) {
	t.Parallel()

	channelID := "test_channel"

	tests := []struct {
		name       string
		attachment *slack.Attachment
		setupMock  func(m *mocks.MockSlackClientInterface)
		wantErrMsg string
	}{
		{
			name: "success",
			attachment: &slack.Attachment{
				Text:  "Test attachment",
				Color: "good",
			},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any(), gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
		{
			name: "client error",
			attachment: &slack.Attachment{
				Text:  "Test attachment",
				Color: "good",
			},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any(), gomock.Any()).
					Return("", "", errTestSlack)
			},
			wantErrMsg: "slack PostMessageAttachment",
		},
		{
			name:       "empty attachment",
			attachment: &slack.Attachment{},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any(), gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
		{
			name: "with fields",
			attachment: &slack.Attachment{
				Text:  "Test attachment with fields",
				Color: "good",
				Fields: []slack.AttachmentField{
					{Title: "Field 1", Value: "Value 1", Short: true},
					{Title: "Field 2", Value: "Value 2", Short: true},
				},
			},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any(), gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockSlackClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewSlackRepositoryWithInterface(mockClient, channelID)
			err := repo.PostMessageAttachment(tc.attachment)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Contains(t, err.Error(), errTestSlack.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSlackRepository_PostMessage(t *testing.T) {
	t.Parallel()

	channelID := "test_channel"

	tests := []struct {
		name       string
		options    []slack.MsgOption
		setupMock  func(m *mocks.MockSlackClientInterface)
		wantErrMsg string
	}{
		{
			name: "client error",
			options: []slack.MsgOption{
				slack.MsgOptionText("Test message", false),
				slack.MsgOptionUsername("TestBot"),
			},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any()).
					Return("", "", errTestSlack)
			},
			wantErrMsg: "slack PostMessage",
		},
		{
			name: "multiple options",
			options: []slack.MsgOption{
				slack.MsgOptionText("Test message", false),
				slack.MsgOptionUsername("TestBot"),
				slack.MsgOptionIconEmoji(":robot_face:"),
				slack.MsgOptionAsUser(false),
			},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID, gomock.Any()).
					Return("timestamp", "message_id", nil)
			},
		},
		{
			name:    "empty options",
			options: []slack.MsgOption{},
			setupMock: func(m *mocks.MockSlackClientInterface) {
				m.EXPECT().
					PostMessage(channelID).
					Return("timestamp", "message_id", nil)
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockSlackClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewSlackRepositoryWithInterface(mockClient, channelID)
			err := repo.PostMessage(tc.options...)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Contains(t, err.Error(), errTestSlack.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

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
