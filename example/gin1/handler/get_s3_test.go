//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
	"go.uber.org/mock/gomock"
)

var errTestS3GetObject = errors.New("get object failed")

func expectS3HandleFlow(mockClient *mocks.MockAWSS3ClientInterface) {
	const bucket = "test"
	const objectKey = "test.txt"
	const objectText = "aaaaaaaab"
	now := time.Now()

	mockClient.EXPECT().
		CreateBucket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.CreateBucketOutput{}, nil)
	mockClient.EXPECT().
		ListBuckets(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.ListBucketsOutput{
			Buckets: []types.Bucket{{Name: aws.String(bucket), CreationDate: &now}},
		}, nil)
	mockClient.EXPECT().
		PutObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.PutObjectOutput{}, nil)
	mockClient.EXPECT().
		GetObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.GetObjectOutput{
			Body: io.NopCloser(strings.NewReader(objectText)),
		}, nil)
	mockClient.EXPECT().
		ListObjectsV2(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.ListObjectsV2Output{
			Contents: []types.Object{{Key: aws.String(objectKey)}},
		}, nil)
	mockClient.EXPECT().
		DeleteObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.DeleteObjectOutput{}, nil)
	mockClient.EXPECT().
		DeleteBucket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.DeleteBucketOutput{}, nil)
}

func newS3Handler(t *testing.T, ctrl *gomock.Controller) *HTTPHandler {
	t.Helper()

	mockClient := mocks.NewMockAWSS3ClientInterface(ctrl)
	mockUploader := mocks.NewMockAWSS3UploaderClientInterface(ctrl)
	mockDownloader := mocks.NewMockAWSS3DownloaderClientInterface(ctrl)
	mockPresigned := mocks.NewMockAWSS3PresignClientInterface(ctrl)

	expectS3HandleFlow(mockClient)

	repo := repository.NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)
	return NewHTTPHandler(newTestLogger(), nil, nil, repo, nil)
}

func newS3HandlerGetObjectError(t *testing.T, ctrl *gomock.Controller) *HTTPHandler {
	t.Helper()

	mockClient := mocks.NewMockAWSS3ClientInterface(ctrl)
	mockUploader := mocks.NewMockAWSS3UploaderClientInterface(ctrl)
	mockDownloader := mocks.NewMockAWSS3DownloaderClientInterface(ctrl)
	mockPresigned := mocks.NewMockAWSS3PresignClientInterface(ctrl)

	const bucket = "test"
	const objectKey = "test.txt"
	now := time.Now()

	mockClient.EXPECT().
		CreateBucket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.CreateBucketOutput{}, nil)
	mockClient.EXPECT().
		ListBuckets(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.ListBucketsOutput{
			Buckets: []types.Bucket{{Name: aws.String(bucket), CreationDate: &now}},
		}, nil)
	mockClient.EXPECT().
		PutObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.PutObjectOutput{}, nil)
	mockClient.EXPECT().
		GetObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errTestS3GetObject)
	mockClient.EXPECT().
		ListObjectsV2(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.ListObjectsV2Output{
			Contents: []types.Object{{Key: aws.String(objectKey)}},
		}, nil)
	mockClient.EXPECT().
		DeleteObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.DeleteObjectOutput{}, nil)
	mockClient.EXPECT().
		DeleteBucket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&s3.DeleteBucketOutput{}, nil)

	repo := repository.NewAWSS3RepositoryWithInterface(mockClient, mockUploader, mockDownloader, mockPresigned)
	return NewHTTPHandler(newTestLogger(), nil, nil, repo, nil)
}

func TestHTTPHandler_HandleS3(t *testing.T) {
	tests := []struct {
		setup      func(t *testing.T, ctrl *gomock.Controller) *HTTPHandler
		name       string
		wantBody   string
		wantStatus int
	}{
		{
			name: "runs s3 workflow and returns ok",
			setup: func(t *testing.T, ctrl *gomock.Controller) *HTTPHandler {
				return newS3Handler(t, ctrl)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"message":"ok"`,
		},
		{
			name: "get object error does not panic",
			setup: func(t *testing.T, ctrl *gomock.Controller) *HTTPHandler {
				return newS3HandlerGetObjectError(t, ctrl)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"message":"ok"`,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			h := tt.setup(t, ctrl)
			status, body := invokeHandler(t, h.HandleS3)
			require.Equal(t, tt.wantStatus, status)
			require.Contains(t, body, tt.wantBody)
		})
	}
}

func TestHTTPHandler_HandleS3_Route(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{name: "post method returns not found", method: http.MethodPost, wantStatus: http.StatusNotFound},
	}

	h := &HTTPHandler{}
	router := gin.New()
	router.GET("/s3", h.HandleS3)

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/s3", http.NoBody)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
