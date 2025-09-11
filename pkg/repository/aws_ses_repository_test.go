package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSESClient is a mock implementation of SES client for testing
type MockSESClient struct {
	mock.Mock
}

func (m *MockSESClient) SendEmail(ctx context.Context, input *sesv2.SendEmailInput, opts ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sesv2.SendEmailOutput), args.Error(1)
}

func (m *MockSESClient) SendBulkEmail(ctx context.Context, input *sesv2.SendBulkEmailInput, opts ...func(*sesv2.Options)) (*sesv2.SendBulkEmailOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sesv2.SendBulkEmailOutput), args.Error(1)
}

// AWSSESRepositoryWithMock for testing with mock client
type AWSSESRepositoryWithMock struct {
	c                    AWSSESClientInterface
	configurationSetName *string
}

// NewAWSSESRepositoryWithMock creates repository with mock client for testing
func NewAWSSESRepositoryWithMock(mockClient AWSSESClientInterface, configurationSetName *string) *AWSSESRepositoryWithMock {
	return &AWSSESRepositoryWithMock{
		c:                    mockClient,
		configurationSetName: configurationSetName,
	}
}

// SendTextEmail sends text email.
func (r *AWSSESRepositoryWithMock) SendTextEmail(ctx context.Context, from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(ctx, &sesv2.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		Destination: &types.Destination{
			ToAddresses: to,
		},
		ReplyToAddresses: replyTo,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(content),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendTextEmail: %w", err)
	}
	return res, nil
}

// SendHTMLEmail sends HTML email.
func (r *AWSSESRepositoryWithMock) SendHTMLEmail(ctx context.Context, from string, to, replyTo []string, subject, content string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(ctx, &sesv2.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		Destination: &types.Destination{
			ToAddresses: to,
		},
		ReplyToAddresses: replyTo,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(content),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendHTMLEmail: %w", err)
	}
	return res, nil
}

// SendEmail sends email with both text and HTML content.
func (r *AWSSESRepositoryWithMock) SendEmail(ctx context.Context, from string, to, replyTo []string, subject, contentText, contentHTML string) (*sesv2.SendEmailOutput, error) {
	res, err := r.c.SendEmail(ctx, &sesv2.SendEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		Destination: &types.Destination{
			ToAddresses: to,
		},
		ReplyToAddresses: replyTo,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(contentText),
					},
					Html: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(contentHTML),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendEmail: %w", err)
	}
	return res, nil
}

// SendBulkEmail sends bulk email.
func (r *AWSSESRepositoryWithMock) SendBulkEmail(ctx context.Context, from string, replyTo []string, defaultTemplateData string, bulkEmailEntries []types.BulkEmailEntry) (*sesv2.SendBulkEmailOutput, error) {
	res, err := r.c.SendBulkEmail(ctx, &sesv2.SendBulkEmailInput{
		ConfigurationSetName: r.configurationSetName,
		FromEmailAddress:     aws.String(from),
		ReplyToAddresses:     replyTo,
		DefaultContent: &types.BulkEmailContent{
			Template: &types.Template{
				TemplateData: aws.String(defaultTemplateData),
			},
		},
		BulkEmailEntries: bulkEmailEntries,
	})
	if err != nil {
		return nil, fmt.Errorf("ses SendBulkEmail: %w", err)
	}
	return res, nil
}

func TestNewAWSSESRepository(t *testing.T) {
	mockClient := &sesv2.Client{}
	configurationSetName := "test_config"
	repo := NewAWSSESRepository(mockClient, &configurationSetName)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.c)
	assert.Equal(t, &configurationSetName, repo.configurationSetName)
}

func TestAWSSESRepository_SendTextEmail(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "Test email content"

	expectedOutput := &sesv2.SendEmailOutput{
		MessageId: aws.String("test-message-id"),
	}

	mockClient.On("SendEmail", mock.Anything, mock.MatchedBy(func(input *sesv2.SendEmailInput) bool {
		return *input.FromEmailAddress == from &&
			len(input.Destination.ToAddresses) == 1 &&
			input.Destination.ToAddresses[0] == to[0] &&
			*input.Content.Simple.Subject.Data == subject &&
			*input.Content.Simple.Body.Text.Data == content
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendTextEmail(context.Background(), from, to, replyTo, subject, content)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendTextEmail_Error(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithMock(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "Test email content"

	expectedError := errors.New("ses service error")

	mockClient.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.SendTextEmail(context.Background(), from, to, replyTo, subject, content)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ses SendTextEmail")
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendHTMLEmail(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "<h1>Test HTML content</h1>"

	expectedOutput := &sesv2.SendEmailOutput{
		MessageId: aws.String("test-message-id"),
	}

	mockClient.On("SendEmail", mock.Anything, mock.MatchedBy(func(input *sesv2.SendEmailInput) bool {
		return *input.FromEmailAddress == from &&
			len(input.Destination.ToAddresses) == 1 &&
			input.Destination.ToAddresses[0] == to[0] &&
			*input.Content.Simple.Subject.Data == subject &&
			*input.Content.Simple.Body.Html.Data == content
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendHTMLEmail(context.Background(), from, to, replyTo, subject, content)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendEmail(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	contentText := "Test text content"
	contentHTML := "<h1>Test HTML content</h1>"

	expectedOutput := &sesv2.SendEmailOutput{
		MessageId: aws.String("test-message-id"),
	}

	mockClient.On("SendEmail", mock.Anything, mock.MatchedBy(func(input *sesv2.SendEmailInput) bool {
		return *input.FromEmailAddress == from &&
			len(input.Destination.ToAddresses) == 1 &&
			input.Destination.ToAddresses[0] == to[0] &&
			*input.Content.Simple.Subject.Data == subject &&
			*input.Content.Simple.Body.Text.Data == contentText &&
			*input.Content.Simple.Body.Html.Data == contentHTML
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendEmail(context.Background(), from, to, replyTo, subject, contentText, contentHTML)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendHTMLEmail_Error(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "<h1>Test HTML content</h1>"

	expectedError := errors.New("ses service error")

	mockClient.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.SendHTMLEmail(context.Background(), from, to, replyTo, subject, content)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ses SendHTMLEmail")
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendEmail_Error(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	contentText := "Test text content"
	contentHTML := "<h1>Test HTML content</h1>"

	expectedError := errors.New("ses service error")

	mockClient.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.SendEmail(context.Background(), from, to, replyTo, subject, contentText, contentHTML)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ses SendEmail")
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendBulkEmail(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	replyTo := []string{"reply@example.com"}
	defaultTemplateData := "{\"name\":\"Default Name\"}"
	bulkEmailEntries := []types.BulkEmailEntry{
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
	}

	expectedOutput := &sesv2.SendBulkEmailOutput{
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

	mockClient.On("SendBulkEmail", mock.Anything, mock.MatchedBy(func(input *sesv2.SendBulkEmailInput) bool {
		return *input.FromEmailAddress == from &&
			len(input.BulkEmailEntries) == 2 &&
			*input.DefaultContent.Template.TemplateData == defaultTemplateData
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendBulkEmail(context.Background(), from, replyTo, defaultTemplateData, bulkEmailEntries)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	assert.Len(t, result.BulkEmailEntryResults, 2)
	mockClient.AssertExpectations(t)
}

// Test actual AWSSESRepository functions
func TestAWSSESRepositoryReal_SendTextEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepositoryReal_SendHTMLEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepositoryReal_SendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

func TestAWSSESRepositoryReal_SendBulkEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping SES integration test - requires real AWS credentials")
}

// Test edge cases and validation
func TestAWSSESRepository_SendTextEmail_EmptyRecipients(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	to := []string{} // Empty recipients
	replyTo := []string{"reply@example.com"}
	subject := "Test Subject"
	content := "Test email content"

	expectedOutput := &sesv2.SendEmailOutput{
		MessageId: aws.String("test-message-id"),
	}

	mockClient.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendTextEmail(context.Background(), from, to, replyTo, subject, content)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendBulkEmail_EmptyEntries(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	replyTo := []string{"reply@example.com"}
	defaultTemplateData := "{\"name\":\"Default Name\"}"
	bulkEmailEntries := []types.BulkEmailEntry{} // Empty entries

	expectedOutput := &sesv2.SendBulkEmailOutput{
		BulkEmailEntryResults: []types.BulkEmailEntryResult{},
	}

	mockClient.On("SendBulkEmail", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendBulkEmail(context.Background(), from, replyTo, defaultTemplateData, bulkEmailEntries)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	assert.Len(t, result.BulkEmailEntryResults, 0)
	mockClient.AssertExpectations(t)
}

func TestAWSSESRepository_SendBulkEmail_LargeBatch(t *testing.T) {
	mockClient := &MockSESClient{}
	configurationSetName := "test_config"

	repo := NewAWSSESRepositoryWithInterface(mockClient, &configurationSetName)

	from := "sender@example.com"
	replyTo := []string{"reply@example.com"}
	defaultTemplateData := "{\"name\":\"Default Name\"}"

	// Create 50 entries (maximum allowed)
	bulkEmailEntries := make([]types.BulkEmailEntry, 50)
	for i := 0; i < 50; i++ {
		bulkEmailEntries[i] = types.BulkEmailEntry{
			Destination: &types.Destination{
				ToAddresses: []string{fmt.Sprintf("user%d@example.com", i)},
			},
		}
	}

	expectedOutput := &sesv2.SendBulkEmailOutput{
		BulkEmailEntryResults: make([]types.BulkEmailEntryResult, 50),
	}

	mockClient.On("SendBulkEmail", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	result, err := repo.SendBulkEmail(context.Background(), from, replyTo, defaultTemplateData, bulkEmailEntries)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	assert.Len(t, result.BulkEmailEntryResults, 50)
	mockClient.AssertExpectations(t)
}
