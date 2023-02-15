package repository

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	maxBufferSize int = 512
)

// AWSS3RepositoryInterface interface.
type AWSS3RepositoryInterface interface {
	GetObject(bucket, path string) io.ReaderAt
	PutObjectFile(bucket, key, filePath string) (*s3.PutObjectOutput, error)
	PutObjectText(bucket, key string, text *string) (*s3.PutObjectOutput, error)
	DeleteObject(bucket, key string) (*s3.DeleteObjectOutput, error)
	DeleteObjects(bucket string, keys []string) (*s3.DeleteObjectsOutput, error)
	ListObjectsV2(bucket, prefix string) (*s3.ListObjectsV2Output, error)
	ListBuckets() (*s3.ListBucketsOutput, error)
	CreateBucket(bucket string) (*s3.CreateBucketOutput, error)
	DeleteBucket(bucket string) (*s3.DeleteBucketOutput, error)
	GetPresignedURL(bucket, key string, expire time.Duration) (string, error)
	Upload(bucket, key, filePath string) (*s3manager.UploadOutput, error)
	Download(bucket, key, filePath string) error
}

// AWSS3Repository struct.
type AWSS3Repository struct {
	s3      *s3.S3
	session *session.Session
}

// NewAWSS3Repository returns AWSS3Repository instance.
func NewAWSS3Repository(s *s3.S3, sess *session.Session) *AWSS3Repository {
	return &AWSS3Repository{
		s3:      s,
		session: sess,
	}
}

// GetObject retrieves objects from Amazon S3.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (r *AWSS3Repository) GetObject(bucket, key string) (*s3.GetObjectOutput, error) {
	return r.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
}

// PutObjectFile adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSS3Repository) PutObjectFile(bucket, key, filePath string) (*s3.PutObjectOutput, error) {
	path := filepath.Clean(filePath)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	// Get content-type
	buf := make([]byte, maxBufferSize)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buf)

	return r.s3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(r.normalizePath(key)),
		Body:        file,
		ContentType: &contentType,
	})
}

// PutObjectText adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSS3Repository) PutObjectText(bucket, key string, text *string) (*s3.PutObjectOutput, error) {
	contentType := http.DetectContentType([]byte(*text))
	return r.s3.PutObject(&s3.PutObjectInput{
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
func (r *AWSS3Repository) DeleteObject(bucket, key string) (*s3.DeleteObjectOutput, error) {
	o, err := r.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	if err != nil {
		return nil, err
	}

	err = r.s3.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	return o, err
}

// DeleteObjects action enables you to delete multiple objects from a bucket using a single HTTP request.
// If you know the object keys that you want to delete, then this action provides a suitable alternative to
// sending individual delete requests, reducing per-request overhead.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html
func (r *AWSS3Repository) DeleteObjects(bucket string, keys []string) (*s3.DeleteObjectsOutput, error) {
	objects := make([]*s3.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = &s3.ObjectIdentifier{
			Key: aws.String(r.normalizePath(key)),
		}
	}
	return r.s3.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: objects,
		},
	})
}

// ListObjectsV2 returns some or all (up to 1,000) of the objects in a bucket with each request.
// You can use the request parameters as selection criteria to return a subset of the objects in a bucket.
// A 200 OK response can contain valid or invalid XML. Make sure to design your application to parse the contents
// of the response and handle it appropriately. Objects are returned sorted in an ascending order of the respective
// key names in the list. For more information about listing objects, see Listing object keys programmatically
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
func (r *AWSS3Repository) ListObjectsV2(bucket, prefix string) (*s3.ListObjectsV2Output, error) {
	return r.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
}

// ListBuckets returns a list of all buckets owned by the authenticated sender of the request.
// To use this operation, you must have the s3:ListAllMyBuckets permission.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html
func (r *AWSS3Repository) ListBuckets() (*s3.ListBucketsOutput, error) {
	return r.s3.ListBuckets(&s3.ListBucketsInput{})
}

// CreateBucket creates a new S3 bucket. To create a bucket, you must register with Amazon S3
// and have a valid AWS Access Key ID to authenticate requests. Anonymous requests are never allowed
// to create buckets. By creating the bucket, you become the bucket owner.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (r *AWSS3Repository) CreateBucket(bucket string) (*s3.CreateBucketOutput, error) {
	return r.s3.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
}

// DeleteBucket deletes the S3 bucket. All objects (including all object versions and delete markers) in the bucket
// must be deleted before the bucket itself can be deleted.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (r *AWSS3Repository) DeleteBucket(bucket string) (*s3.DeleteBucketOutput, error) {
	o, err := r.s3.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	err = r.s3.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	return o, err
}

// GetPresignedURL creates a Pre-Singed URL.
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-presigned-urls.html
func (r *AWSS3Repository) GetPresignedURL(bucket, key string, expire time.Duration) (string, error) {
	req, _ := r.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	return req.Presign(expire)
}

// Upload adds an object to a bucket.
func (r *AWSS3Repository) Upload(bucket, key, filePath string) (*s3manager.UploadOutput, error) {
	path := filepath.Clean(filePath)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	// Get content-type
	buf := make([]byte, maxBufferSize)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buf)

	uploader := s3manager.NewUploader(r.session)
	return uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Body:        file,
		Key:         aws.String(r.normalizePath(key)),
		ContentType: &contentType,
	})
}

// Download retrieves objects from Amazon S3.
func (r *AWSS3Repository) Download(bucket, key, filePath string) error {
	path := filepath.Clean(filePath)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	downloader := s3manager.NewDownloader(r.session)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
	return err
}

// normalizePath normalizes the path (removing any leading slash)
func (r *AWSS3Repository) normalizePath(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	return path
}
