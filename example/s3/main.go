package main

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
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
	l := logger.NewLogger(logrusLogger)

	// --------------------------------------------------------------
	// S3(minio)
	// --------------------------------------------------------------
	s3Region := os.Getenv("S3_REGION")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ID := os.Getenv("S3_ID")
	s3Secret := os.Getenv("S3_SECRET")
	s3Token := os.Getenv("S3_TOKEN")

	s3Config, err := infrastructure.GetAWSV2S3Config(l, s3ID, s3Secret, s3Token, s3Region, s3Endpoint, true)
	if err != nil {
		panic(err)
	}

	// --------------------------------------------------------------
	// example: S3
	// --------------------------------------------------------------
	awsS3Repository := repository.NewAWSV2S3Repository(s3.NewFromConfig(s3Config, func(o *s3.Options) {
		o.UsePathStyle = true
	}))
	text := "abc"
	bucket := "test"

	// Create Bucket
	_, err = awsS3Repository.CreateBucket(bucket)
	if err != nil {
		l.WithError(err).Errorf("can't create s3 bucket")
	}

	// ListBuckets
	listBuckets, err := awsS3Repository.ListBuckets()
	if err == nil {
		for _, b := range listBuckets.Buckets {
			l.Infof("bucket = %s(%s)", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		}
	} else {
		l.WithError(err).Errorf("can't list of s3 bucket")
	}

	// Put Object
	_, err = awsS3Repository.PutObjectText(bucket, "test.txt", &text)
	if err != nil {
		l.WithError(err).Errorf("can't put s3 object")
	}

	// Get Object
	object, err := awsS3Repository.GetObject(bucket, "test.txt")
	if err != nil {
		l.WithError(err).Errorf("can't get s3 object")
	}
	rc := object.Body
	defer func() {
		err = rc.Close()
		if err != nil {
			l.WithError(err).Errorf("can't close body")
		}
	}()
	text, err = utils.GetStringFromReadCloser(rc)
	if err != nil {
		l.WithError(err).Errorf("can't get text")
	}
	l.Infof("text.txt = %s", text)

	// ListObjectV2
	listObjects, err := awsS3Repository.ListObjectsV2(bucket, "")
	if err == nil {
		for _, o := range listObjects.Contents {
			l.Infof("Object key = %s", aws.StringValue(o.Key))
		}
	} else {
		l.WithError(err).Errorf("can't list of s3 object")
	}

	// Delete Object
	_, err = awsS3Repository.DeleteObject(bucket, "test.txt")
	if err != nil {
		l.WithError(err).Errorf("can't delete s3 object")
	}

	// Delete Bucket
	_, err = awsS3Repository.DeleteBucket(bucket)
	if err != nil {
		l.WithError(err).Errorf("can't delete s3 bucket")
	}
}
