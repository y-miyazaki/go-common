package repository

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAWSBedrockClient is a mock implementation of AWSBedrockClientInterface
type MockAWSBedrockClient struct {
	mock.Mock
}

func (m *MockAWSBedrockClient) InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bedrockruntime.InvokeModelOutput), args.Error(1)
}

func (m *MockAWSBedrockClient) InvokeModelWithResponseStream(ctx context.Context, params *bedrockruntime.InvokeModelWithResponseStreamInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bedrockruntime.InvokeModelWithResponseStreamOutput), args.Error(1)
}

func (m *MockAWSBedrockClient) Converse(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bedrockruntime.ConverseOutput), args.Error(1)
}

func TestNewAWSBedrockRepository(t *testing.T) {
	// Test with nil client
	repo := NewAWSBedrockRepository(nil)
	assert.NotNil(t, repo)
	assert.Nil(t, repo.Client)

	// Test with real client (without mocking AWS calls)
	mockClient := new(MockAWSBedrockClient)
	repo = NewAWSBedrockRepositoryWithInterface(mockClient)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.Client)
}

func TestAWSBedrockRepository_InvokeModel(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	payload := []byte(`{"prompt":"test"}`)
	expectedResponse := []byte(`{"completion":"response"}`)

	mockClient.On("InvokeModel", mock.Anything, mock.MatchedBy(func(input *bedrockruntime.InvokeModelInput) bool {
		return *input.ModelId == modelID && string(input.Body) == string(payload)
	})).Return(&bedrockruntime.InvokeModelOutput{
		Body: expectedResponse,
	}, nil)

	result, err := repo.InvokeModel(context.Background(), modelID, payload)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_InvokeModelWithStream(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	payload := []byte(`{"prompt":"test"}`)
	expectedOutput := &bedrockruntime.InvokeModelWithResponseStreamOutput{}

	mockClient.On("InvokeModelWithResponseStream", mock.Anything, mock.MatchedBy(func(input *bedrockruntime.InvokeModelWithResponseStreamInput) bool {
		return *input.ModelId == modelID && string(input.Body) == string(payload)
	})).Return(expectedOutput, nil)

	result, err := repo.InvokeModelWithStream(context.Background(), modelID, payload)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_Converse(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	message := map[string]string{"role": "user", "content": "Hello"}
	expectedOutput := &bedrockruntime.ConverseOutput{}

	mockClient.On("Converse", mock.Anything, mock.MatchedBy(func(input *bedrockruntime.ConverseInput) bool {
		return *input.ModelId == modelID
	})).Return(expectedOutput, nil)

	result, err := repo.Converse(context.Background(), modelID, message)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_InvokeModel_Error(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	payload := []byte(`{"prompt":"test"}`)

	mockClient.On("InvokeModel", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	_, err := repo.InvokeModel(context.Background(), modelID, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invoke model")
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_InvokeModelWithStream_Error(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	payload := []byte(`{"prompt":"test"}`)

	mockClient.On("InvokeModelWithResponseStream", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	_, err := repo.InvokeModelWithStream(context.Background(), modelID, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invoke model with stream")
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_Converse_Error(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	message := map[string]string{"role": "user", "content": "Hello"}

	mockClient.On("Converse", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	_, err := repo.Converse(context.Background(), modelID, message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "converse")
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_InvokeModelWithFileData(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	fileData := []byte("test image data")
	mimeType := "image/png"
	additionalPayload := map[string]any{"prompt": "Describe this image"}
	expectedResponse := []byte(`{"completion":"response"}`)

	mockClient.On("InvokeModel", mock.Anything, mock.MatchedBy(func(input *bedrockruntime.InvokeModelInput) bool {
		return *input.ModelId == modelID
	})).Return(&bedrockruntime.InvokeModelOutput{
		Body: expectedResponse,
	}, nil)

	result, err := repo.InvokeModelWithFileData(context.Background(), modelID, fileData, mimeType, additionalPayload)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockClient.AssertExpectations(t)
}

func TestAWSBedrockRepository_InvokeModelWithFileData_Error(t *testing.T) {
	mockClient := new(MockAWSBedrockClient)
	repo := NewAWSBedrockRepositoryWithInterface(mockClient)

	modelID := "anthropic.claude-v2"
	fileData := []byte("test image data")
	mimeType := "image/png"
	additionalPayload := map[string]any{"prompt": "Describe this image"}

	mockClient.On("InvokeModel", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	_, err := repo.InvokeModelWithFileData(context.Background(), modelID, fileData, mimeType, additionalPayload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invoke model with file data")
	mockClient.AssertExpectations(t)
}
