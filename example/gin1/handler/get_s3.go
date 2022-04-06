package handler

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/pkg/utils"
)

// GetS3 handler
func (h *HTTPHandler) GetS3(c *gin.Context) {
	text := "aaaaaaaab"
	bucket := "test"

	// Create Bucket
	_, err := h.awsS3Repository.CreateBucket(bucket)
	if err != nil {
		h.Logger.WithError(err).Errorf("can't create s3 bucket")
	}

	// ListBuckets
	listBuckets, err := h.awsS3Repository.ListBuckets()
	if err == nil {
		for _, b := range listBuckets.Buckets {
			h.Logger.Infof("bucket = %s(%s)", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		}
	} else {
		h.Logger.WithError(err).Errorf("can't list of s3 bucket")
	}

	// Put Object
	_, err = h.awsS3Repository.PutObjectText(bucket, "test.txt", &text)
	if err != nil {
		h.Logger.WithError(err).Errorf("can't put s3 object")
	}

	// Get Object
	object, err := h.awsS3Repository.GetObject(bucket, "test.txt")
	if err != nil {
		h.Logger.WithError(err).Errorf("can't get s3 object")
	}
	rc := object.Body
	defer func() {
		err = rc.Close()
		if err != nil {
			h.Logger.WithError(err).Errorf("can't close body")
		}
	}()
	text, err = utils.GetStringFromReadCloser(rc)
	if err != nil {
		h.Logger.WithError(err).Errorf("can't get text")
	}
	h.Logger.Infof("text.txt = %s", text)

	// ListObjectV2
	listObjects, err := h.awsS3Repository.ListObjectsV2(bucket, "")
	if err == nil {
		for _, o := range listObjects.Contents {
			h.Logger.Infof("Object key = %s", aws.StringValue(o.Key))
		}
	} else {
		h.Logger.WithError(err).Errorf("can't list of s3 object")
	}

	// Delete Object
	_, err = h.awsS3Repository.DeleteObject(bucket, "test.txt")
	if err != nil {
		h.Logger.WithError(err).Errorf("can't delete s3 object")
	}

	// Delete Bucket
	_, err = h.awsS3Repository.DeleteBucket(bucket)
	if err != nil {
		h.Logger.WithError(err).Errorf("can't delete s3 bucket")
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
