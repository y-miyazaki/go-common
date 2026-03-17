// Package infrastructure provides AWS S3 integration helpers using AWS SDK v2.
package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewAWSS3 returns S3 client using AWS SDK v2.
func NewAWSS3(c *aws.Config, optFns ...func(*s3.Options)) *s3.Client { // nolint:gocritic
	return s3.NewFromConfig(*c, optFns...)
}

// NewAWSS3Downloader returns S3 Downloader using AWS SDK v2.
// TODO: migrate to github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager when ready.
func NewAWSS3Downloader(c manager.DownloadAPIClient, options ...func(*manager.Downloader)) *manager.Downloader { //nolint:staticcheck // SA1019: manager.Downloader deprecated, migration pending
	return manager.NewDownloader(c, options...) //nolint:staticcheck // SA1019: manager.NewDownloader deprecated, migration pending
}

// NewAWSS3Uploader returns S3 Uploader using AWS SDK v2.
// TODO: migrate to github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager when ready.
func NewAWSS3Uploader(c manager.UploadAPIClient, options ...func(*manager.Uploader)) *manager.Uploader { //nolint:staticcheck // SA1019: manager.Uploader deprecated, migration pending
	return manager.NewUploader(c, options...) //nolint:staticcheck // SA1019: manager.NewUploader deprecated, migration pending
}
