package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestSES = errors.New("ses service error")

func TestNewAWSSESRepository(t *testing.T) {
	t.Parallel()

	mockClient := &sesv2.Client{}
	configurationSetName := "test_config"
	repo := NewAWSSESRepository(mockClient, &configurationSetName)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
	require.Equal(t, &configurationSetName, repo.configurationSetName)
}

func TestAWSSESRepository_SendTextEmail(t *testing.T) {
	t.Parallel()

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "Test email content"
	configurationSetName := "test_config"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSSESClientInterface)
		want       *sesv2.SendEmailOutput
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSSESClientInterface) {
				expected := &sesv2.SendEmailOutput{
					MessageId: aws.String("test-message-id"),
				}
				m.EXPECT().
					SendEmail(gomock.Any(), gomock.AssignableToTypeOf(&sesv2.SendEmailInput{})).
					DoAndReturn(func(_ context.Context, input *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
						require.Equal(t, from, *input.FromEmailAddress)
						require.Len(t, input.Destination.ToAddresses, 1)
						require.Equal(t, to[0], input.Destination.ToAddresses[0])
						require.Equal(t, subject, *input.Content.Simple.Subject.Data)
						require.Equal(t, content, *input.Content.Simple.Body.Text.Data)
						return expected, nil
					})
			},
			want: &sesv2.SendEmailOutput{
				MessageId: aws.String("test-message-id"),
			},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSSESClientInterface) {
				m.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil, errTestSES)
			},
			wantErrMsg: "ses SendTextEmail",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSESClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)
			got, err := repo.SendTextEmail(context.Background(), from, to, replyTo, subject, content)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestAWSSESRepository_SendHTMLEmail(t *testing.T) {
	t.Parallel()

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "<h1>Test HTML content</h1>"
	configurationSetName := "test_config"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSSESClientInterface)
		want       *sesv2.SendEmailOutput
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSSESClientInterface) {
				expected := &sesv2.SendEmailOutput{
					MessageId: aws.String("test-message-id"),
				}
				m.EXPECT().
					SendEmail(gomock.Any(), gomock.AssignableToTypeOf(&sesv2.SendEmailInput{})).
					DoAndReturn(func(_ context.Context, input *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
						require.Equal(t, from, *input.FromEmailAddress)
						require.Len(t, input.Destination.ToAddresses, 1)
						require.Equal(t, to[0], input.Destination.ToAddresses[0])
						require.Equal(t, subject, *input.Content.Simple.Subject.Data)
						require.Equal(t, content, *input.Content.Simple.Body.Html.Data)
						return expected, nil
					})
			},
			want: &sesv2.SendEmailOutput{
				MessageId: aws.String("test-message-id"),
			},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSSESClientInterface) {
				m.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil, errTestSES)
			},
			wantErrMsg: "ses SendHTMLEmail",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSESClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)
			got, err := repo.SendHTMLEmail(context.Background(), from, to, replyTo, subject, content)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestAWSSESRepository_SendEmail(t *testing.T) {
	t.Parallel()

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	contentText := "Test text content"
	contentHTML := "<h1>Test HTML content</h1>"
	configurationSetName := "test_config"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSSESClientInterface)
		want       *sesv2.SendEmailOutput
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSSESClientInterface) {
				expected := &sesv2.SendEmailOutput{
					MessageId: aws.String("test-message-id"),
				}
				m.EXPECT().
					SendEmail(gomock.Any(), gomock.AssignableToTypeOf(&sesv2.SendEmailInput{})).
					DoAndReturn(func(_ context.Context, input *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
						require.Equal(t, from, *input.FromEmailAddress)
						require.Len(t, input.Destination.ToAddresses, 1)
						require.Equal(t, to[0], input.Destination.ToAddresses[0])
						require.Equal(t, subject, *input.Content.Simple.Subject.Data)
						require.Equal(t, contentText, *input.Content.Simple.Body.Text.Data)
						require.Equal(t, contentHTML, *input.Content.Simple.Body.Html.Data)
						return expected, nil
					})
			},
			want: &sesv2.SendEmailOutput{
				MessageId: aws.String("test-message-id"),
			},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSSESClientInterface) {
				m.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil, errTestSES)
			},
			wantErrMsg: "ses SendEmail",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSESClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)
			got, err := repo.SendEmail(context.Background(), from, to, replyTo, subject, contentText, contentHTML)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestAWSSESRepository_SendBulkEmail(t *testing.T) {
	t.Parallel()

	from := "sender@example.com"
	replyTo := []string{"reply@example.com"}
	defaultTemplateData := `{"name":"Default Name"}`
	configurationSetName := "test_config"

	tests := []struct {
		name             string
		bulkEmailEntries []types.BulkEmailEntry
		setupMock        func(m *mocks.MockAWSSESClientInterface, bulkEmailEntries []types.BulkEmailEntry)
		want             *sesv2.SendBulkEmailOutput
		wantErrMsg       string
	}{
		{
			name: "success",
			bulkEmailEntries: []types.BulkEmailEntry{
				{
					Destination: &types.Destination{
						ToAddresses: []string{"user1@example.com"},
					},
				},
				{
					Destination: &types.Destination{
						ToAddresses: []string{"user2@example.com"},
					},
				},
			},
			setupMock: func(m *mocks.MockAWSSESClientInterface, bulkEmailEntries []types.BulkEmailEntry) {
				expected := &sesv2.SendBulkEmailOutput{
					BulkEmailEntryResults: []types.BulkEmailEntryResult{
						{
							Error:     nil,
							MessageId: aws.String("msg-1"),
							Status:    "SUCCESS",
						},
						{
							Error:     nil,
							MessageId: aws.String("msg-2"),
							Status:    "SUCCESS",
						},
					},
				}
				m.EXPECT().
					SendBulkEmail(gomock.Any(), gomock.AssignableToTypeOf(&sesv2.SendBulkEmailInput{})).
					DoAndReturn(func(_ context.Context, input *sesv2.SendBulkEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendBulkEmailOutput, error) {
						require.Equal(t, from, *input.FromEmailAddress)
						require.Len(t, input.BulkEmailEntries, len(bulkEmailEntries))
						require.Equal(t, defaultTemplateData, *input.DefaultContent.Template.TemplateData)
						return expected, nil
					})
			},
			want: &sesv2.SendBulkEmailOutput{
				BulkEmailEntryResults: []types.BulkEmailEntryResult{
					{
						Error:     nil,
						MessageId: aws.String("msg-1"),
						Status:    "SUCCESS",
					},
					{
						Error:     nil,
						MessageId: aws.String("msg-2"),
						Status:    "SUCCESS",
					},
				},
			},
		},
		{
			name:             "empty entries",
			bulkEmailEntries: []types.BulkEmailEntry{},
			setupMock: func(m *mocks.MockAWSSESClientInterface, bulkEmailEntries []types.BulkEmailEntry) {
				expected := &sesv2.SendBulkEmailOutput{
					BulkEmailEntryResults: []types.BulkEmailEntryResult{},
				}
				m.EXPECT().SendBulkEmail(gomock.Any(), gomock.Any()).Return(expected, nil)
			},
			want: &sesv2.SendBulkEmailOutput{
				BulkEmailEntryResults: []types.BulkEmailEntryResult{},
			},
		},
		{
			name: "large batch",
			setupMock: func(m *mocks.MockAWSSESClientInterface, bulkEmailEntries []types.BulkEmailEntry) {
				expected := &sesv2.SendBulkEmailOutput{
					BulkEmailEntryResults: make([]types.BulkEmailEntryResult, len(bulkEmailEntries)),
				}
				m.EXPECT().SendBulkEmail(gomock.Any(), gomock.Any()).Return(expected, nil)
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bulkEmailEntries := tc.bulkEmailEntries
			if tc.name == "large batch" {
				bulkEmailEntries = make([]types.BulkEmailEntry, 50)
				for j := range 50 {
					bulkEmailEntries[j] = types.BulkEmailEntry{
						Destination: &types.Destination{
							ToAddresses: []string{fmt.Sprintf("user%d@example.com", j)},
						},
					}
				}
			}

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSESClientInterface(ctrl)
			tc.setupMock(mockClient, bulkEmailEntries)

			repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)
			got, err := repo.SendBulkEmail(context.Background(), from, replyTo, defaultTemplateData, bulkEmailEntries)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			if tc.want != nil {
				require.Equal(t, tc.want, got)
			}
			require.Len(t, got.BulkEmailEntryResults, len(bulkEmailEntries))
		})
	}
}

func TestAWSSESRepositoryReal_SendTextEmail(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepositoryReal_SendHTMLEmail(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepositoryReal_SendEmail(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepositoryReal_SendBulkEmail(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepository_SendTextEmail_EmptyRecipients(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSSESClientInterface(ctrl)
	configurationSetName := "test_config"

	expectedOutput := &sesv2.SendEmailOutput{
		MessageId: aws.String("test-message-id"),
	}

	mockClient.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(expectedOutput, nil)

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)
	got, err := repo.SendTextEmail(
		context.Background(),
		"sender@example.com",
		[]string{},
		[]string{"reply@example.com"},
		"Test Subject",
		"Test email content",
	)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, got)
}
