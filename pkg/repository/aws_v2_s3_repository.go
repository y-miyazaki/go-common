package repository

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	maxBufferSizeV2 int = 512
)

// AWSV2S3RepositoryInterface interface.
type AWSV2S3RepositoryInterface interface {
	GetObject(bucket, path string) io.ReaderAt
	PutObjectFile(bucket, key, filePath string) (*s3.PutObjectOutput, error)
	PutObjectText(bucket, key string, text *string) (*s3.PutObjectOutput, error)
	DeleteObject(bucket, key string) (*s3.DeleteObjectOutput, error)
	DeleteObjects(bucket string, keys []string) (*s3.DeleteObjectsOutput, error)
	ListObjectsV2(bucket, prefix string) (*s3.ListObjectsV2Output, error)
	ListBuckets() (*s3.ListBucketsOutput, error)
	CreateBucket(bucket string) (*s3.CreateBucketOutput, error)
	DeleteBucket(bucket string) (*s3.DeleteBucketOutput, error)
	GetPresignedURL(bucket, key string, expire time.Duration) (*v4.PresignedHTTPRequest, error)
	Upload(bucket, key, filePath string) (*manager.UploadOutput, error)
	Download(bucket, key, filePath string) error
}

// AWSV2S3Repository struct.
type AWSV2S3Repository struct {
	c          *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	presigned  *s3.PresignClient
}

// NewAWSV2S3Repository returns AWSV2S3Repository instance.
func NewAWSV2S3Repository(client *s3.Client) *AWSV2S3Repository {
	return &AWSV2S3Repository{
		c:          client,
		uploader:   manager.NewUploader(client),
		downloader: manager.NewDownloader(client),
		presigned:  s3.NewPresignClient(client),
	}
}

// GetObject retrieves objects from Amazon S3.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (r *AWSV2S3Repository) GetObject(bucket, key string) (*s3.GetObjectOutput, error) {
	return r.c.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
}

// PutObjectFile adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSV2S3Repository) PutObjectFile(bucket, key, filePath string) (*s3.PutObjectOutput, error) {
	path := filepath.Clean(filePath)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	// Get content-type
	buf := make([]byte, maxBufferSizeV2)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buf)

	return r.c.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        file,
		ContentType: &contentType,
	})
}

// PutObjectText adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSV2S3Repository) PutObjectText(bucket, key string, text *string) (*s3.PutObjectOutput, error) {
	contentType := http.DetectContentType([]byte(*text))
	return r.c.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        bytes.NewReader([]byte(*text)),
		ContentType: &contentType,
	})
}

// DeleteObject removes the null version (if there is one) of an object and inserts a delete marker,
// which becomes the latest version of the object. If there isn't a null version, Amazon S3 does not remove
// any objects but will still respond that the command was successful.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (r *AWSV2S3Repository) DeleteObject(bucket, key string) (*s3.DeleteObjectOutput, error) {
	return r.c.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
}

// DeleteObjects action enables you to delete multiple objects from a bucket using a single HTTP request.
// If you know the object keys that you want to delete, then this action provides a suitable alternative to
// sending individual delete requests, reducing per-request overhead.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html
func (r *AWSV2S3Repository) DeleteObjects(bucket string, keys []string) (*s3.DeleteObjectsOutput, error) {
	var objectIds []types.ObjectIdentifier
	for _, key := range keys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	return r.c.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{
			Objects: objectIds,
		},
	})
}

// ListObjectsV2 returns some or all (up to 1,000) of the objects in a bucket with each request.
// You can use the request parameters as selection criteria to return a subset of the objects in a bucket.
// A 200 OK response can contain valid or invalid XML. Make sure to design your application to parse the contents
// of the response and handle it appropriately. Objects are returned sorted in an ascending order of the respective
// key names in the list. For more information about listing objects, see Listing object keys programmatically
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
func (r *AWSV2S3Repository) ListObjectsV2(bucket, prefix string) (*s3.ListObjectsV2Output, error) {
	return r.c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
}

// ListBuckets returns a list of all buckets owned by the authenticated sender of the request.
// To use this operation, you must have the s3:ListAllMyBuckets permission.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html
func (r *AWSV2S3Repository) ListBuckets() (*s3.ListBucketsOutput, error) {
	return r.c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
}

// CreateBucket creates a new S3 bucket. To create a bucket, you must register with Amazon S3
// and have a valid AWS Access Key ID to authenticate requests. Anonymous requests are never allowed
// to create buckets. By creating the bucket, you become the bucket owner.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (r *AWSV2S3Repository) CreateBucket(bucket string) (*s3.CreateBucketOutput, error) {
	return r.c.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
}

// DeleteBucket deletes the S3 bucket. All objects (including all object versions and delete markers) in the bucket
// must be deleted before the bucket itself can be deleted.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (r *AWSV2S3Repository) DeleteBucket(bucket string) (*s3.DeleteBucketOutput, error) {
	return r.c.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
}

// GetPresignedURL creates a Pre-Singed URL.
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-presigned-urls.html
func (r *AWSV2S3Repository) GetPresignedURL(bucket, key string, expire time.Duration) (*v4.PresignedHTTPRequest, error) {
	presignDuration := func(options *s3.PresignOptions) {
		options.Expires = 1 * time.Minute
	}
	return r.presigned.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, presignDuration)
}

// Upload adds an object to a bucket.
func (r *AWSV2S3Repository) Upload(bucket, key, filePath string) (*manager.UploadOutput, error) {
	path := filepath.Clean(filePath)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	// Get content-type
	buf := make([]byte, maxBufferSizeV2)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buf)

	return r.uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        file,
		ContentType: &contentType,
	})
}

// Download retrieves objects from Amazon S3.
func (r *AWSV2S3Repository) Download(bucket, key, filePath string) error {
	path := filepath.Clean(filePath)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = r.downloader.Download(context.Background(), file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	return err
}

// normalizePath normalizes the path (removing any leading slash)
func (r *AWSV2S3Repository) normalizePath(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	return path
}
