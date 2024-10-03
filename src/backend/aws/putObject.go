package ovh

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Creates a new object in the S3 cluster in some bucket. Returns the URL to the file (assuming the file is meant to be public)
// TODO: implement some function to compress images and mp4s if compress is set to true
func (w *Worker) AddObject(buffer *bytes.Buffer, bucket, key, contentType string) (fileUrl string, error error) {
	reader := bytes.NewReader(buffer.Bytes())

	_, err := w.S3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
		Body:        reader,
	})

	endpointWithoutHttps := strings.TrimPrefix(os.Getenv("AWS_ENDPOINT"), "https://")

	return fmt.Sprintf("http://%s.%s%s", bucket, endpointWithoutHttps, key), err
}
