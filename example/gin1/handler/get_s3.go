package handler

import (
	"net/http"

	"go-common/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"
)

// HandleS3 demonstrates AWS S3 operations including bucket creation, object management, and cleanup.
func (h *HTTPHandler) HandleS3(c *gin.Context) {
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
		for i := range listBuckets.Buckets {
			h.Logger.Infof("bucket = %s(%s)", aws.ToString(listBuckets.Buckets[i].Name), aws.ToTime(listBuckets.Buckets[i].CreationDate))
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
		for i := range listObjects.Contents {
			h.Logger.Infof("Object key = %s", aws.ToString(listObjects.Contents[i].Key))
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
