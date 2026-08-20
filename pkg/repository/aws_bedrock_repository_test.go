package repository

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestBedrock = errors.New("bedrock error")

func TestNewAWSBedrockRepository(t *testing.T) {
	t.Parallel()

	repo := NewAWSBedrockRepository(nil)
	require.NotNil(t, repo)
	require.Nil(t, repo.Client)

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSBedrockClientInterface(ctrl)
	repo = NewAWSBedrockRepositoryWithInterface(mockClient)
	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
}

func TestAWSBedrockRepository_InvokeModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		modelID    string
		payload    []byte
		setupMock  func(m *mocks.MockAWSBedrockClientInterface)
		want       []byte
		wantErrMsg string
	}{
		{
			name:    "success",
			modelID: "anthropic.claude-v2",
			payload: []byte(`{"prompt":"test"}`),
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				expectedResponse := []byte(`{"completion":"response"}`)
				m.EXPECT().
					InvokeModel(gomock.Any(), gomock.AssignableToTypeOf(&bedrockruntime.InvokeModelInput{})).
					DoAndReturn(func(_ context.Context, input *bedrockruntime.InvokeModelInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
						require.Equal(t, "anthropic.claude-v2", *input.ModelId)
						require.True(t, bytes.Equal(input.Body, []byte(`{"prompt":"test"}`)))
						return &bedrockruntime.InvokeModelOutput{Body: expectedResponse}, nil
					})
			},
			want: []byte(`{"completion":"response"}`),
		},
		{
			name:    "client error",
			modelID: "anthropic.claude-v2",
			payload: []byte(`{"prompt":"test"}`),
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				m.EXPECT().
					InvokeModel(gomock.Any(), gomock.Any()).
					Return(nil, errTestBedrock)
			},
			wantErrMsg: "invoke model",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSBedrockClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSBedrockRepositoryWithInterface(mockClient)
			got, err := repo.InvokeModel(context.Background(), tc.modelID, tc.payload)

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

func TestAWSBedrockRepository_InvokeModelWithStream(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		modelID    string
		payload    []byte
		setupMock  func(m *mocks.MockAWSBedrockClientInterface)
		want       *bedrockruntime.InvokeModelWithResponseStreamOutput
		wantErrMsg string
	}{
		{
			name:    "success",
			modelID: "anthropic.claude-v2",
			payload: []byte(`{"prompt":"test"}`),
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				expectedOutput := &bedrockruntime.InvokeModelWithResponseStreamOutput{}
				m.EXPECT().
					InvokeModelWithResponseStream(gomock.Any(), gomock.AssignableToTypeOf(&bedrockruntime.InvokeModelWithResponseStreamInput{})).
					DoAndReturn(func(_ context.Context, input *bedrockruntime.InvokeModelWithResponseStreamInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
						require.Equal(t, "anthropic.claude-v2", *input.ModelId)
						require.True(t, bytes.Equal(input.Body, []byte(`{"prompt":"test"}`)))
						return expectedOutput, nil
					})
			},
			want: &bedrockruntime.InvokeModelWithResponseStreamOutput{},
		},
		{
			name:    "client error",
			modelID: "anthropic.claude-v2",
			payload: []byte(`{"prompt":"test"}`),
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				m.EXPECT().
					InvokeModelWithResponseStream(gomock.Any(), gomock.Any()).
					Return(nil, errTestBedrock)
			},
			wantErrMsg: "invoke model with stream",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSBedrockClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSBedrockRepositoryWithInterface(mockClient)
			got, err := repo.InvokeModelWithStream(context.Background(), tc.modelID, tc.payload)

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

func TestAWSBedrockRepository_Converse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		modelID    string
		message    map[string]string
		setupMock  func(m *mocks.MockAWSBedrockClientInterface)
		wantErrMsg string
	}{
		{
			name:    "success",
			modelID: "anthropic.claude-v2",
			message: map[string]string{"role": "user", "content": "Hello"},
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				m.EXPECT().
					Converse(gomock.Any(), gomock.AssignableToTypeOf(&bedrockruntime.ConverseInput{})).
					DoAndReturn(func(_ context.Context, input *bedrockruntime.ConverseInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
						require.Equal(t, "anthropic.claude-v2", *input.ModelId)
						return &bedrockruntime.ConverseOutput{}, nil
					})
			},
		},
		{
			name:    "client error",
			modelID: "anthropic.claude-v2",
			message: map[string]string{"role": "user", "content": "Hello"},
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				m.EXPECT().
					Converse(gomock.Any(), gomock.Any()).
					Return(nil, errTestBedrock)
			},
			wantErrMsg: "converse",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSBedrockClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSBedrockRepositoryWithInterface(mockClient)
			got, err := repo.Converse(context.Background(), tc.modelID, tc.message)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Empty(t, got)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got)
		})
	}
}

func TestAWSBedrockRepository_InvokeModelWithFileData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		modelID    string
		fileData   []byte
		payload    []byte
		setupMock  func(m *mocks.MockAWSBedrockClientInterface)
		want       []byte
		wantErrMsg string
	}{
		{
			name:     "success",
			modelID:  "anthropic.claude-v2",
			fileData: []byte("test image data"),
			payload:  []byte(`{"prompt": "Describe this image"}`),
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				expectedResponse := []byte(`{"completion":"response"}`)
				m.EXPECT().
					InvokeModel(gomock.Any(), gomock.AssignableToTypeOf(&bedrockruntime.InvokeModelInput{})).
					DoAndReturn(func(_ context.Context, input *bedrockruntime.InvokeModelInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
						require.Equal(t, "anthropic.claude-v2", *input.ModelId)
						return &bedrockruntime.InvokeModelOutput{Body: expectedResponse}, nil
					})
			},
			want: []byte(`{"completion":"response"}`),
		},
		{
			name:     "client error",
			modelID:  "anthropic.claude-v2",
			fileData: []byte("test image data"),
			payload:  []byte(`{"prompt": "Describe this image"}`),
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				m.EXPECT().
					InvokeModel(gomock.Any(), gomock.Any()).
					Return(nil, errTestBedrock)
			},
			wantErrMsg: "invoke model with file data",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSBedrockClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSBedrockRepositoryWithInterface(mockClient)
			got, err := repo.InvokeModelWithFileData(context.Background(), tc.modelID, tc.fileData, tc.payload)

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

func TestAWSBedrockRepository_InvokeModelWithFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		filePath   string
		setupFile  func(t *testing.T) string
		setupMock  func(m *mocks.MockAWSBedrockClientInterface)
		want       []byte
		wantErrMsg string
	}{
		{
			name: "success",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpFile, err := os.CreateTemp("", "test-image-*.png")
				require.NoError(t, err)
				t.Cleanup(func() { os.Remove(tmpFile.Name()) })

				_, err = tmpFile.Write([]byte("test image data"))
				require.NoError(t, err)
				require.NoError(t, tmpFile.Close())
				return tmpFile.Name()
			},
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				expectedResponse := []byte(`{"completion":"response"}`)
				m.EXPECT().
					InvokeModel(gomock.Any(), gomock.AssignableToTypeOf(&bedrockruntime.InvokeModelInput{})).
					DoAndReturn(func(_ context.Context, input *bedrockruntime.InvokeModelInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
						require.Equal(t, "anthropic.claude-v2", *input.ModelId)
						return &bedrockruntime.InvokeModelOutput{Body: expectedResponse}, nil
					})
			},
			want: []byte(`{"completion":"response"}`),
		},
		{
			name:       "file not found",
			filePath:   "/nonexistent/file.png",
			wantErrMsg: "open file",
		},
		{
			name: "invoke error",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpFile, err := os.CreateTemp("", "test-image-*.png")
				require.NoError(t, err)
				t.Cleanup(func() { os.Remove(tmpFile.Name()) })

				_, err = tmpFile.Write([]byte("test image data"))
				require.NoError(t, err)
				require.NoError(t, tmpFile.Close())
				return tmpFile.Name()
			},
			setupMock: func(m *mocks.MockAWSBedrockClientInterface) {
				m.EXPECT().
					InvokeModel(gomock.Any(), gomock.Any()).
					Return(nil, errTestBedrock)
			},
			wantErrMsg: "invoke model with file",
		},
	}

	modelID := "anthropic.claude-v2"
	payload := []byte(`{"prompt": "Describe this image"}`)

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSBedrockClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			filePath := tc.filePath
			if tc.setupFile != nil {
				filePath = tc.setupFile(t)
			}

			repo := NewAWSBedrockRepositoryWithInterface(mockClient)
			got, err := repo.InvokeModelWithFile(context.Background(), modelID, filePath, payload)

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
