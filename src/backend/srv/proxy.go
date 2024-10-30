package srv

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

const (
	errOptKeyNil       = "No image key provided"
	errOptBucketInv    = "Invalid bucket type provided"
	errFileReadFailure = "Error while reading file"
)

// Proxies S3 images to be used in html <img> tags
func (s *Server) Proxy(c *gin.Context) {
	key := c.Query("key")
	bucket := c.DefaultQuery("bucket", "blog")

	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptKeyNil,
			"message": nil,
		})
		return
	}

	switch bucket {
	case "blog":
		bucket = os.Getenv("AWS_BLOG_BUCKET_NAME")
	case "gallery":
		bucket = os.Getenv("AWS_GALLERY_BUCKET_NAME")
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptBucketInv,
			"message": nil,
		})
		return
	}

	obj, err := s.Ovh.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to proxy file from S3",
		})
		return
	}

	defer obj.Body.Close()

	c.Header("Content-Type", *obj.ContentType)
	c.Header("Content-Length", strconv.Itoa(int(*obj.ContentLength)))
	c.Header("Cache-Control", "max-age=3600")

	body, err := io.ReadAll(obj.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   errFileReadFailure,
			"message": nil,
		})
		return
	}

	c.Data(http.StatusOK, *obj.ContentType, body)
}
