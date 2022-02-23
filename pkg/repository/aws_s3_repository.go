package repository

import (
	"bytes"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

// AWSS3RepositoryInterface interface.
type AWSS3RepositoryInterface interface {
	GetObject(bucket, path string) io.ReaderAt
}

// AWSS3Repository struct.
type AWSS3Repository struct {
	e  *logrus.Entry
	s3 *s3.S3
}

// NewAWSS3Repository returns AWSS3Repository instance.
func NewAWSS3Repository(
	e *logrus.Entry,
	s3 *s3.S3,

) *AWSS3Repository {
	return &AWSS3Repository{
		e:  e,
		s3: s3,
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

// PutObject adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSS3Repository) PutObjectFile(bucket, key, filePath string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return r.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
		Body:   file,
	})
}

// PutObject adds an object to a bucket.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (r *AWSS3Repository) PutObjectText(bucket, key, text string) (*s3.PutObjectOutput, error) {
	return r.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
		Body:   bytes.NewReader([]byte(text)),
	})
}

// DeleteObject removes the null version (if there is one) of an object and inserts a delete marker,
// which becomes the latest version of the object. If there isn't a null version, Amazon S3 does not remove
// any objects but will still respond that the command was successful.
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (r *AWSS3Repository) DeleteObject(bucket, key string) (*s3.DeleteObjectOutput, error) {
	return r.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(r.normalizePath(key)),
	})
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

// Normalizes normalizes the path (removing any leading slash)
func (fur *AWSS3Repository) normalizePath(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	return path
}
