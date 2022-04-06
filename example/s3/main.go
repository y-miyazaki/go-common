package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"github.com/y-miyazaki/go-common/pkg/utils"
)

func main() {
	// --------------------------------------------------------------
	// logrus
	// --------------------------------------------------------------
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	logger := infrastructure.NewLogger(logrusLogger)

	// --------------------------------------------------------------
	// S3(minio)
	// --------------------------------------------------------------
	s3Region := os.Getenv("S3_REGION")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ID := os.Getenv("S3_ID")
	s3Secret := os.Getenv("S3_SECRET")
	s3Token := os.Getenv("S3_TOKEN")

	s3SessionOptions := infrastructure.GetS3DefaultOptions()
	s3Config := infrastructure.GetS3Config(logger.Entry, s3ID, s3Secret, s3Token, s3Region, s3Endpoint, true)
	session := infrastructure.NewS3Session(s3SessionOptions)

	// --------------------------------------------------------------
	// example: S3
	// --------------------------------------------------------------
	awsS3Repository := repository.NewAWSS3Repository(logger.Entry, session, s3Config)
	text := "aaaaaaaab"
	bucket := "test"

	// Create Bucket
	_, err := awsS3Repository.CreateBucket(bucket)
	if err != nil {
		logger.WithError(err).Errorf("can't create s3 bucket")
	}

	// ListBuckets
	listBuckets, err := awsS3Repository.ListBuckets()
	if err == nil {
		for _, b := range listBuckets.Buckets {
			logger.Infof("bucket = %s(%s)", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		}
	} else {
		logger.WithError(err).Errorf("can't list of s3 bucket")
	}

	// Put Object
	_, err = awsS3Repository.PutObjectText(bucket, "test.txt", &text)
	if err != nil {
		logger.WithError(err).Errorf("can't put s3 object")
	}

	// Get Object
	object, err := awsS3Repository.GetObject(bucket, "test.txt")
	if err != nil {
		logger.WithError(err).Errorf("can't get s3 object")
	}
	rc := object.Body
	defer rc.Close()

	text, err = utils.GetStringFromReadCloser(rc)
	if err != nil {
		logger.WithError(err).Errorf("can't get text")
	}
	logger.Infof("text.txt = %s", text)

	// ListObjectV2
	listObjects, err := awsS3Repository.ListObjectsV2(bucket, "")
	if err == nil {
		for _, o := range listObjects.Contents {
			logger.Infof("Object key = %s", aws.StringValue(o.Key))
		}
	} else {
		logger.WithError(err).Errorf("can't list of s3 object")
	}

	// Delete Object
	_, err = awsS3Repository.DeleteObject(bucket, "test.txt")
	if err != nil {
		logger.WithError(err).Errorf("can't delete s3 object")
	}

	// Delete Bucket
	_, err = awsS3Repository.DeleteBucket(bucket)
	if err != nil {
		logger.WithError(err).Errorf("can't delete s3 bucket")
	}

	// // Upload
	// _, err = awsS3Repository.Upload(bucket, "test.txt", "./example/s3/cmd.zip")
	// if err != nil {
	// 	logger.WithError(err).Errorf("can't upload file")
	// }
}
