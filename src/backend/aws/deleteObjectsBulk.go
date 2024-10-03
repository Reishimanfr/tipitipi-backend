package ovh

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (w *Worker) DeleteObjectsBulk(bucket string, keys []string) error {
	deleteObjects := make([]*s3.ObjectIdentifier, len(keys))

	for _, key := range keys {
		deleteObjects = append(deleteObjects, &s3.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	_, err := w.S3.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: &bucket,
		Delete: &s3.Delete{
			Objects: deleteObjects,
		},
	})

	return err
}