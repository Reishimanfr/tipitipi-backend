package ovh

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	maxAttempts       = flag.Int("multipart-upload-attempts", 3, "The max amount of attempts before giving up on uploading a file to S3")
	maxPartSize int64 = *flag.Int64("multipart-max-part-size", 5, "The max part size for multipart uploads (in MiB)") << 20
)

// Creates a new object in the S3 cluster in some bucket. Returns the URL to the file (assuming the file is meant to be public)
// TODO: implement some function to compress images and mp4s if compress is set to true
func (w *Worker) AddObject(bucket, key string, buffer []byte) (fileUrl *string, error error) {
	fileType := http.DetectContentType(buffer)

	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(fileType),
		ACL:         aws.String("public-read"),
	}

	resp, err := w.S3.CreateMultipartUpload(input)
	if err != nil {
		return nil, err
	}

	var partLen, cur int64
	remaining := int64(len(buffer))
	completedParts := []*s3.CompletedPart{}

	partIdx := 1

	for cur = 0; remaining != 0; cur += partLen {
		if remaining < maxPartSize {
			partLen = remaining
		} else {
			partLen = maxPartSize
		}

		completedPart, err := uploadPart(w.S3, resp, buffer[cur:cur+partLen], partIdx)
		if err != nil {
			err2 := abortMultipartUpload(w.S3, resp)
			if err2 != nil {
				return nil, fmt.Errorf("failed to both upload the file and cancel it's upload request. Errors: %v\n%v", err, err2)
			}

			return nil, err
		}

		remaining -= partLen
		partIdx++
		completedParts = append(completedParts, completedPart)
	}

	_, err = completeMultipartUpload(w.S3, resp, completedParts)
	if err != nil {
		return nil, err
	}

	endpointWithoutHttps := strings.TrimPrefix(os.Getenv("AWS_ENDPOINT"), "https://")

	return aws.String(fmt.Sprintf("http://%s.%s%s", bucket, endpointWithoutHttps, key)), nil
}

// Uploads a part of the provided data to S3
func uploadPart(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, b []byte, partIdx int) (*s3.CompletedPart, error) {
	attempt := 1

	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(b),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		UploadId:      resp.UploadId,
		PartNumber:    aws.Int64(int64(partIdx)),
		ContentLength: aws.Int64(int64(len(b))),
	}

	for attempt <= *maxAttempts {
		uploadResult, err := svc.UploadPart(partInput)
		if err != nil && attempt == *maxAttempts {
			if aerr, ok := err.(awserr.Error); ok {
				return nil, aerr
			}
			return nil, err
		}

		if err != nil {
			attempt++
			continue
		}

		return &s3.CompletedPart{
			ETag:       uploadResult.ETag,
			PartNumber: aws.Int64(int64(partIdx)),
		}, nil
	}

	return nil, nil
}

func abortMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput) error {
	_, err := svc.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	})

	return err
}

func completeMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	return svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
}
