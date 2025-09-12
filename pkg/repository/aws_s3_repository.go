package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	// nolint:revive
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	maxBufferSize int = 512
)

// AWSS3ClientInterface interface for mocking S3 client
type AWSS3ClientInterface interface {
	// GetObject retrieves an object from S3
	GetObject(_ context.Context, _ *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	// PutObject uploads an object to S3
	PutObject(_ context.Context, _ *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	// DeleteObject deletes an object from S3
	DeleteObject(_ context.Context, _ *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	// DeleteObjects deletes multiple objects from S3
	DeleteObjects(_ context.Context, _ *s3.DeleteObjectsInput, _ ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	// ListObjectsV2 lists objects in an S3 bucket
	ListObjectsV2(_ context.Context, _ *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	// ListBuckets lists all S3 buckets
	ListBuckets(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	// CreateBucket creates a new S3 bucket
	CreateBucket(_ context.Context, _ *s3.CreateBucketInput, _ ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	// DeleteBucket deletes an S3 bucket
	DeleteBucket(_ context.Context, _ *s3.DeleteBucketInput, _ ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

// AWSS3PresignClientInterface interface for mocking S3 presign client
type AWSS3PresignClientInterface interface {
	// PresignGetObject generates a presigned URL for getting an S3 object
	PresignGetObject(_ context.Context, _ *s3.GetObjectInput, _ ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// AWSS3UploaderClientInterface interface for mocking S3 uploader
type AWSS3UploaderClientInterface interface {
	// Upload uploads an object to S3 using the manager uploader
	Upload(_ context.Context, _ *s3.PutObjectInput, _ ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

// AWSS3DownloaderClientInterface interface for mocking S3 downloader
type AWSS3DownloaderClientInterface interface {
	// Download downloads an object from S3 using the manager downloader
	Download(_ context.Context, _ io.WriterAt, _ *s3.GetObjectInput, _ ...func(*manager.Downloader)) (int64, error)
}

// AWSS3Repository struct implements AWSS3RepositoryInterface using AWS SDK v2.
type AWSS3Repository struct {
	Client     AWSS3ClientInterface
	uploader   AWSS3UploaderClientInterface
	downloader AWSS3DownloaderClientInterface
	presigned  AWSS3PresignClientInterface
}

// NewAWSS3Repository returns AWSS3Repository instance backed by AWS SDK v2 client.
func NewAWSS3Repository(client *s3.Client) *AWSS3Repository {
	if client == nil {
		return &AWSS3Repository{}
	}
	return &AWSS3Repository{
		Client:     client,
		uploader:   manager.NewUploader(client),
		downloader: manager.NewDownloader(client),
		presigned:  s3.NewPresignClient(client),
	}
}

// NewAWSS3RepositoryWithOther returns AWSS3Repository instance backed by AWS SDK v2 client.
func NewAWSS3RepositoryWithOther(client *s3.Client, uploader *manager.Uploader, downloader *manager.Downloader, presigned *s3.PresignClient) *AWSS3Repository {
	return &AWSS3Repository{
		Client:     client,
		uploader:   uploader,
		downloader: downloader,
		presigned:  presigned,
	}
}

// NewAWSS3RepositoryWithInterface returns AWSS3Repository instance backed by interfaces for testing.
func NewAWSS3RepositoryWithInterface(client AWSS3ClientInterface, uploader AWSS3UploaderClientInterface, downloader AWSS3DownloaderClientInterface, presigned AWSS3PresignClientInterface) *AWSS3Repository {
	return &AWSS3Repository{
		Client:     client,
		uploader:   uploader,
		downloader: downloader,
		presigned:  presigned,
	}
}

// GetObject retrieves objects from Amazon S3.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (r *AWSS3Repository) GetObject(bucket, key string) (*s3.GetObjectOutput, error) {
	out, err := r.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject: %w", err)
	}
	return out, nil
}

// PutObjectFile adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSS3Repository) PutObjectFile(bucket, key, filePath string) (*s3.PutObjectOutput, error) {
	path := filepath.Clean(filePath)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if cErr := file.Close(); cErr != nil {
			log.Printf("warning: failed to close file: %v", cErr)
		}
	}()
	// Get content-type
	buf := make([]byte, maxBufferSize)
	_, err = file.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	contentType := http.DetectContentType(buf)

	out, err := r.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        file,
		ContentType: &contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 PutObject: %w", err)
	}
	return out, nil
}

// PutObjectText adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSS3Repository) PutObjectText(bucket, key string, text *string) (*s3.PutObjectOutput, error) {
	contentType := http.DetectContentType([]byte(*text))
	out, err := r.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        bytes.NewReader([]byte(*text)),
		ContentType: &contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 PutObject: %w", err)
	}
	return out, nil
}

// DeleteObject removes the null version (if there is one) of an object and inserts a delete marker,
// which becomes the latest version of the object. If there isn't a null version, Amazon S3 does not remove
// any objects but will still respond that the command was successful.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (r *AWSS3Repository) DeleteObject(bucket, key string) (*s3.DeleteObjectOutput, error) {
	out, err := r.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 DeleteObject: %w", err)
	}
	return out, nil
}

// DeleteObjects action enables you to delete multiple objects from a bucket using a single HTTP request.
// If you know the object keys that you want to delete, then this action provides a suitable alternative to
// sending individual delete requests, reducing per-request overhead.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html
func (r *AWSS3Repository) DeleteObjects(bucket string, keys []string) (*s3.DeleteObjectsOutput, error) {
	var objectIDs []types.ObjectIdentifier
	for _, key := range keys {
		objectIDs = append(objectIDs, types.ObjectIdentifier{Key: aws.String(r.normalizePath(key))})
	}
	out, err := r.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{
			Objects: objectIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("s3 DeleteObjects: %w", err)
	}
	return out, nil
}

// ListObjectsV2 returns some or all (up to 1,000) of the objects in a bucket with each request.
// You can use the request parameters as selection criteria to return a subset of the objects in a bucket.
// A 200 OK response can contain valid or invalid XML. Make sure to design your application to parse the contents
// of the response and handle it appropriately. Objects are returned sorted in an ascending order of the respective
// key names in the list. For more information about listing objects, see Listing object keys programmatically
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
func (r *AWSS3Repository) ListObjectsV2(bucket, prefix string) (*s3.ListObjectsV2Output, error) {
	out, err := r.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 ListObjectsV2: %w", err)
	}
	return out, nil
}

// ListBuckets returns a list of all buckets owned by the authenticated sender of the request.
// To use this operation, you must have the s3:ListAllMyBuckets permission.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html
func (r *AWSS3Repository) ListBuckets() (*s3.ListBucketsOutput, error) {
	out, err := r.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("s3 ListBuckets: %w", err)
	}
	return out, nil
}

// CreateBucket creates a new S3 bucket. To create a bucket, you must register with Amazon S3
// and have a valid AWS Access Key ID to authenticate requests. Anonymous requests are never allowed
// to create buckets. By creating the bucket, you become the bucket owner.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (r *AWSS3Repository) CreateBucket(bucket string) (*s3.CreateBucketOutput, error) {
	out, err := r.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 CreateBucket: %w", err)
	}
	return out, nil
}

// DeleteBucket deletes the S3 bucket. All objects (including all object versions and delete markers) in the bucket
// must be deleted before the bucket itself can be deleted.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (r *AWSS3Repository) DeleteBucket(bucket string) (*s3.DeleteBucketOutput, error) {
	out, err := r.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 DeleteBucket: %w", err)
	}
	return out, nil
}

// GetPresignedURL creates a Pre-Singed URL.
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-presigned-urls.html
func (r *AWSS3Repository) GetPresignedURL(bucket, key string, expire time.Duration) (*v4.PresignedHTTPRequest, error) {
	// Use provided expire duration, default to 1 minute when zero
	exp := expire
	if exp <= 0 {
		exp = 1 * time.Minute
	}
	presignDuration := func(options *s3.PresignOptions) {
		options.Expires = exp
	}
	out, err := r.presigned.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	}, presignDuration)
	if err != nil {
		return nil, fmt.Errorf("s3 PresignGetObject: %w", err)
	}
	return out, nil
}

// Upload adds an object to a bucket.
func (r *AWSS3Repository) Upload(bucket, key, filePath string) (*manager.UploadOutput, error) {
	path := filepath.Clean(filePath)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if cErr := file.Close(); cErr != nil {
			log.Printf("warning: failed to close file: %v", cErr)
		}
	}()

	// Get content-type
	buf := make([]byte, maxBufferSize)
	_, err = file.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	contentType := http.DetectContentType(buf)

	out, err := r.uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        file,
		ContentType: &contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 uploader Upload: %w", err)
	}
	return out, nil
}

// Download retrieves objects from Amazon S3.
func (r *AWSS3Repository) Download(bucket, key, filePath string) error {
	path := filepath.Clean(filePath)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if cErr := file.Close(); cErr != nil {
			log.Printf("warning: failed to close file: %v", cErr)
		}
	}()

	_, err = r.downloader.Download(context.Background(), file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	if err != nil {
		return fmt.Errorf("s3 downloader Download: %w", err)
	}
	return nil
}

// normalizePath normalizes the path (removing any leading slash)
func (*AWSS3Repository) normalizePath(p string) string {
	if p != "" && p[0] == '/' {
		return p[1:]
	}
	return p
}
