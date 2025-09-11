package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockS3Client is a mock implementation of S3ClientInterface for testing
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func (m *MockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (m *MockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}

func (m *MockS3Client) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, opts ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.DeleteObjectsOutput), args.Error(1)
}

func (m *MockS3Client) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func (m *MockS3Client) ListBuckets(ctx context.Context, input *s3.ListBucketsInput, opts ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListBucketsOutput), args.Error(1)
}

func (m *MockS3Client) CreateBucket(ctx context.Context, input *s3.CreateBucketInput, opts ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.CreateBucketOutput), args.Error(1)
}

func (m *MockS3Client) DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput, opts ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.DeleteBucketOutput), args.Error(1)
}

// MockS3PresignClient is a mock implementation of S3PresignClientInterface for testing
type MockS3PresignClient struct {
	mock.Mock
}

func (m *MockS3PresignClient) PresignGetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v4.PresignedHTTPRequest), args.Error(1)
}

// MockS3Uploader is a mock implementation of S3UploaderClientInterface for testing
type MockS3Uploader struct {
	mock.Mock
}

func (m *MockS3Uploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*manager.UploadOutput), args.Error(1)
}

// MockS3Downloader is a mock implementation of S3DownloaderClientInterface for testing
type MockS3Downloader struct {
	mock.Mock
}

func (m *MockS3Downloader) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, opts ...func(*manager.Downloader)) (int64, error) {
	args := m.Called(ctx, w, input, opts)
	return args.Get(0).(int64), args.Error(1)
}

func TestNewAWSS3RepositoryWithInterface(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.c)
	assert.Equal(t, mockUploader, repo.uploader)
	assert.Equal(t, mockDownloader, repo.downloader)
	assert.Equal(t, mockPresigned, repo.presigned)
}

func TestAWSS3Repository_GetObject(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte("test content"))),
	}

	mockClient.On("GetObject", mock.Anything, mock.MatchedBy(func(input *s3.GetObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.GetObject("test-bucket", "test-key")

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_GetObject_Error(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedError := errors.New("object not found")

	mockClient.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.GetObject("test-bucket", "test-key")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "s3 GetObject")
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_PutObjectText(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.PutObjectOutput{}

	mockClient.On("PutObject", mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(expectedOutput, nil)

	text := "test content"
	result, err := repo.PutObjectText("test-bucket", "test-key", &text)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_PutObjectFile(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-file-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testContent := "test file content"
	_, err = tempFile.WriteString(testContent)
	assert.NoError(t, err)
	tempFile.Close()

	expectedOutput := &s3.PutObjectOutput{}

	mockClient.On("PutObject", mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.PutObjectFile("test-bucket", "test-key", tempFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_DeleteObject(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.DeleteObjectOutput{}

	mockClient.On("DeleteObject", mock.Anything, mock.MatchedBy(func(input *s3.DeleteObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.DeleteObject("test-bucket", "test-key")

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_ListObjectsV2(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key:  aws.String("test-key"),
				Size: aws.Int64(100),
			},
		},
	}

	mockClient.On("ListObjectsV2", mock.Anything, mock.MatchedBy(func(input *s3.ListObjectsV2Input) bool {
		return *input.Bucket == "test-bucket" && *input.Prefix == "test-prefix"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.ListObjectsV2("test-bucket", "test-prefix")

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_GetPresignedURL(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &v4.PresignedHTTPRequest{
		URL: "https://test-bucket.s3.amazonaws.com/test-key?X-Amz-Expires=60",
	}

	mockPresigned.On("PresignGetObject", mock.Anything, mock.MatchedBy(func(input *s3.GetObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.GetPresignedURL("test-bucket", "test-key", time.Minute)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockPresigned.AssertExpectations(t)
}

func TestAWSS3Repository_NormalizePath(t *testing.T) {
	repo := &AWSS3Repository{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path without leading slash",
			input:    "test/path",
			expected: "test/path",
		},
		{
			name:     "path with leading slash",
			input:    "/test/path",
			expected: "test/path",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.normalizePath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewAWSS3Repository(t *testing.T) {
	// Test with nil client
	repo := NewAWSS3Repository(nil)
	assert.NotNil(t, repo)
	assert.Nil(t, repo.c)
	assert.Nil(t, repo.uploader)
	assert.Nil(t, repo.downloader)
	assert.Nil(t, repo.presigned)
}

func TestAWSS3Repository_ListBuckets(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.ListBucketsOutput{
		Buckets: []types.Bucket{
			{
				Name:         aws.String("test-bucket-1"),
				CreationDate: aws.Time(time.Now()),
			},
			{
				Name:         aws.String("test-bucket-2"),
				CreationDate: aws.Time(time.Now()),
			},
		},
	}

	mockClient.On("ListBuckets", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	result, err := repo.ListBuckets()

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_CreateBucket(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.CreateBucketOutput{}

	mockClient.On("CreateBucket", mock.Anything, mock.MatchedBy(func(input *s3.CreateBucketInput) bool {
		return *input.Bucket == "test-bucket"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.CreateBucket("test-bucket")

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_DeleteBucket(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.DeleteBucketOutput{}

	mockClient.On("DeleteBucket", mock.Anything, mock.MatchedBy(func(input *s3.DeleteBucketInput) bool {
		return *input.Bucket == "test-bucket"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.DeleteBucket("test-bucket")

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_DeleteObjects(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &s3.DeleteObjectsOutput{}

	mockClient.On("DeleteObjects", mock.Anything, mock.MatchedBy(func(input *s3.DeleteObjectsInput) bool {
		return *input.Bucket == "test-bucket" && len(input.Delete.Objects) == 2
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.DeleteObjects("test-bucket", []string{"key1", "key2"})

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSS3Repository_Upload(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	expectedOutput := &manager.UploadOutput{}

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_upload_*")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("test content")
	assert.NoError(t, err)
	tempFile.Close()

	mockUploader.On("Upload", mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.Upload("test-bucket", "test-key", tempFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockUploader.AssertExpectations(t)
}

func TestAWSS3Repository_Download(t *testing.T) {
	mockClient := &MockS3Client{}
	mockUploader := &MockS3Uploader{}
	mockDownloader := &MockS3Downloader{}
	mockPresigned := &MockS3PresignClient{}

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_download_*")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	mockDownloader.On("Download", mock.Anything, mock.Anything, mock.MatchedBy(func(input *s3.GetObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-key"
	}), mock.Anything).Return(int64(100), nil)

	err = repo.Download("test-bucket", "test-key", tempFile.Name())

	assert.NoError(t, err)
	mockDownloader.AssertExpectations(t)
}
