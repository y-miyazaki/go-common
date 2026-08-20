package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
	"go.uber.org/mock/gomock"
)

var (
	errTestS3GetObject     = errors.New("object not found")
	errTestS3DownloadObject = errors.New("download object failed")
)

func newS3Repo(t *testing.T, ctrl *gomock.Controller) (
	*AWSS3Repository,
	*mocks.MockAWSS3ClientInterface,
	*mocks.MockAWSS3UploaderClientInterface,
	*mocks.MockAWSS3DownloaderClientInterface,
	*mocks.MockAWSS3PresignClientInterface,
) {
	t.Helper()

	mockClient := mocks.NewMockAWSS3ClientInterface(ctrl)
	mockUploader := mocks.NewMockAWSS3UploaderClientInterface(ctrl)
	mockDownloader := mocks.NewMockAWSS3DownloaderClientInterface(ctrl)
	mockPresigned := mocks.NewMockAWSS3PresignClientInterface(ctrl)
	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)
	return repo, mockClient, mockUploader, mockDownloader, mockPresigned
}

func newS3RepoWithTransfer(t *testing.T, ctrl *gomock.Controller) (
	*AWSS3Repository,
	*mocks.MockAWSS3ClientInterface,
	*mocks.MockAWSS3UploaderClientInterface,
	*mocks.MockAWSS3DownloaderClientInterface,
	*mocks.MockAWSS3PresignClientInterface,
	*mocks.MockAWSS3TransferClientInterface,
) {
	t.Helper()

	mockClient := mocks.NewMockAWSS3ClientInterface(ctrl)
	mockUploader := mocks.NewMockAWSS3UploaderClientInterface(ctrl)
	mockDownloader := mocks.NewMockAWSS3DownloaderClientInterface(ctrl)
	mockPresigned := mocks.NewMockAWSS3PresignClientInterface(ctrl)
	mockTransfer := mocks.NewMockAWSS3TransferClientInterface(ctrl)
	repo := NewAWSS3RepositoryWithTransferClient(mockClient, mockUploader, mockDownloader, mockPresigned, mockTransfer)
	return repo, mockClient, mockUploader, mockDownloader, mockPresigned, mockTransfer
}

func TestNewAWSS3RepositoryWithInterface(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSS3ClientInterface(ctrl)
	mockUploader := mocks.NewMockAWSS3UploaderClientInterface(ctrl)
	mockDownloader := mocks.NewMockAWSS3DownloaderClientInterface(ctrl)
	mockPresigned := mocks.NewMockAWSS3PresignClientInterface(ctrl)

	repo := NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)
	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
	require.Equal(t, mockUploader, repo.uploader)
	require.Equal(t, mockDownloader, repo.downloader)
	require.Equal(t, mockPresigned, repo.presigned)
}

func TestAWSS3Repository_GetObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte("test content"))),
	}

	mockClient.EXPECT().
		GetObject(gomock.Any(), gomock.AssignableToTypeOf(&s3.GetObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	result, err := repo.GetObject("test-bucket", "test-key")

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_GetObject_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	mockClient.EXPECT().
		GetObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errTestS3GetObject)

	result, err := repo.GetObject("test-bucket", "test-key")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "s3 GetObject")
	require.ErrorIs(t, err, errTestS3GetObject)
}

func TestAWSS3Repository_PutObjectText(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.PutObjectOutput{}

	mockClient.EXPECT().
		PutObject(gomock.Any(), gomock.AssignableToTypeOf(&s3.PutObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	text := "test content"
	result, err := repo.PutObjectText("test-bucket", "test-key", &text)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_PutObjectFile(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	tempFile, err := os.CreateTemp("", "test-file-*.txt")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testContent := "test file content"
	_, err = tempFile.WriteString(testContent)
	require.NoError(t, err)
	require.NoError(t, tempFile.Close())

	expectedOutput := &s3.PutObjectOutput{}

	mockClient.EXPECT().
		PutObject(gomock.Any(), gomock.AssignableToTypeOf(&s3.PutObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	result, err := repo.PutObjectFile("test-bucket", "test-key", tempFile.Name())

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_DeleteObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.DeleteObjectOutput{}

	mockClient.EXPECT().
		DeleteObject(gomock.Any(), gomock.AssignableToTypeOf(&s3.DeleteObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	result, err := repo.DeleteObject("test-bucket", "test-key")

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_ListObjectsV2(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key:  aws.String("test-key"),
				Size: aws.Int64(100),
			},
		},
	}

	mockClient.EXPECT().
		ListObjectsV2(gomock.Any(), gomock.AssignableToTypeOf(&s3.ListObjectsV2Input{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Prefix)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-prefix", *input.Prefix)
			return expectedOutput, nil
		})

	result, err := repo.ListObjectsV2("test-bucket", "test-prefix")

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_GetPresignedURL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, _, _, _, mockPresigned := newS3Repo(t, ctrl)

	expectedOutput := &v4.PresignedHTTPRequest{
		URL: "https://test-bucket.s3.amazonaws.com/test-key?X-Amz-Expires=60",
	}

	mockPresigned.EXPECT().
		PresignGetObject(gomock.Any(), gomock.AssignableToTypeOf(&s3.GetObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.GetObjectInput, _ ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	result, err := repo.GetPresignedURL("test-bucket", "test-key", time.Minute)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_NormalizePath(t *testing.T) {
	t.Parallel()

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

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := repo.normalizePath(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNewAWSS3Repository(t *testing.T) {
	t.Parallel()

	repo := NewAWSS3Repository(nil)
	require.NotNil(t, repo)
	require.Nil(t, repo.Client)
	require.Nil(t, repo.uploader)
	require.Nil(t, repo.downloader)
	require.Nil(t, repo.presigned)
	require.Nil(t, repo.transferClient)
}

func TestAWSS3Repository_ListBuckets(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

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

	mockClient.EXPECT().
		ListBuckets(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	result, err := repo.ListBuckets()

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_CreateBucket(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.CreateBucketOutput{}

	mockClient.EXPECT().
		CreateBucket(gomock.Any(), gomock.AssignableToTypeOf(&s3.CreateBucketInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.CreateBucketInput, _ ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
			require.NotNil(t, input.Bucket)
			require.Equal(t, "test-bucket", *input.Bucket)
			return expectedOutput, nil
		})

	result, err := repo.CreateBucket("test-bucket")

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_DeleteBucket(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.DeleteBucketOutput{}

	mockClient.EXPECT().
		DeleteBucket(gomock.Any(), gomock.AssignableToTypeOf(&s3.DeleteBucketInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.DeleteBucketInput, _ ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
			require.NotNil(t, input.Bucket)
			require.Equal(t, "test-bucket", *input.Bucket)
			return expectedOutput, nil
		})

	result, err := repo.DeleteBucket("test-bucket")

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_DeleteObjects(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, mockClient, _, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &s3.DeleteObjectsOutput{}

	mockClient.EXPECT().
		DeleteObjects(gomock.Any(), gomock.AssignableToTypeOf(&s3.DeleteObjectsInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.DeleteObjectsInput, _ ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
			require.NotNil(t, input.Bucket)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Len(t, input.Delete.Objects, 2)
			return expectedOutput, nil
		})

	result, err := repo.DeleteObjects("test-bucket", []string{"key1", "key2"})

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_Upload(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, _, mockUploader, _, _ := newS3Repo(t, ctrl)

	expectedOutput := &manager.UploadOutput{}

	tempFile, err := os.CreateTemp("", "test_upload_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("test content")
	require.NoError(t, err)
	require.NoError(t, tempFile.Close())

	mockUploader.EXPECT().
		Upload(gomock.Any(), gomock.AssignableToTypeOf(&s3.PutObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *s3.PutObjectInput, _ ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	result, err := repo.Upload("test-bucket", "test-key", tempFile.Name())

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_Download(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, _, _, mockDownloader, _ := newS3Repo(t, ctrl)

	tempFile, err := os.CreateTemp("", "test_download_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	require.NoError(t, tempFile.Close())

	mockDownloader.EXPECT().
		Download(gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf(&s3.GetObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ io.WriterAt, input *s3.GetObjectInput, _ ...func(*manager.Downloader)) (int64, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return int64(100), nil
		})

	err = repo.Download("test-bucket", "test-key", tempFile.Name())

	require.NoError(t, err)
}

func TestAWSS3Repository_UploadObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, _, _, _, _, mockTransfer := newS3RepoWithTransfer(t, ctrl)

	expectedOutput := &transfermanager.UploadObjectOutput{}

	tempFile, err := os.CreateTemp("", "test_upload_object_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("test content")
	require.NoError(t, err)
	require.NoError(t, tempFile.Close())

	mockTransfer.EXPECT().
		UploadObject(gomock.Any(), gomock.AssignableToTypeOf(&transfermanager.UploadObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *transfermanager.UploadObjectInput, _ ...func(*transfermanager.Options)) (*transfermanager.UploadObjectOutput, error) {
			if input == nil || input.Bucket == nil || input.Key == nil || input.ContentType == nil {
				return nil, errors.New("invalid input")
			}
			body, readErr := io.ReadAll(input.Body)
			if readErr != nil {
				return nil, readErr
			}
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			require.Equal(t, "test content", string(body))
			require.Equal(t, "text/plain; charset=utf-8", *input.ContentType)
			return expectedOutput, nil
		})

	result, err := repo.UploadObject("test-bucket", "test-key", tempFile.Name())

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}

func TestAWSS3Repository_DownloadObject_Error_NoPartialFile(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, _, _, _, _, mockTransfer := newS3RepoWithTransfer(t, ctrl)

	tempDir := t.TempDir()
	destinationPath := filepath.Join(tempDir, "downloaded.txt")

	mockTransfer.EXPECT().
		DownloadObject(gomock.Any(), gomock.AssignableToTypeOf(&transfermanager.DownloadObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *transfermanager.DownloadObjectInput, _ ...func(*transfermanager.Options)) (*transfermanager.DownloadObjectOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return nil, errTestS3DownloadObject
		})

	result, err := repo.DownloadObject("test-bucket", "test-key", destinationPath)

	require.Error(t, err)
	require.Nil(t, result)
	_, statErr := os.Stat(destinationPath)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestAWSS3Repository_DownloadObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo, _, _, _, _, mockTransfer := newS3RepoWithTransfer(t, ctrl)

	tempFile, err := os.CreateTemp("", "test_download_object_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	require.NoError(t, tempFile.Close())

	expectedOutput := &transfermanager.DownloadObjectOutput{}
	mockTransfer.EXPECT().
		DownloadObject(gomock.Any(), gomock.AssignableToTypeOf(&transfermanager.DownloadObjectInput{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, input *transfermanager.DownloadObjectInput, _ ...func(*transfermanager.Options)) (*transfermanager.DownloadObjectOutput, error) {
			require.NotNil(t, input.Bucket)
			require.NotNil(t, input.Key)
			require.Equal(t, "test-bucket", *input.Bucket)
			require.Equal(t, "test-key", *input.Key)
			return expectedOutput, nil
		})

	result, err := repo.DownloadObject("test-bucket", "test-key", tempFile.Name())

	require.NoError(t, err)
	require.Equal(t, expectedOutput, result)
}
