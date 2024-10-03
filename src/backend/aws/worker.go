package ovh

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Worker struct {
	Session *session.Session
	S3      *s3.S3
}

func NewWorker(accessKey, secretKey, region, endpoint string) (*Worker, error) {
	s, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	s3Client := s3.New(s)

	// TODO: implement checking if required buckets exist
	// result, err := s3Client.ListBuckets(nil)

	return &Worker{
		Session: s,
		S3:      s3Client,
	}, nil
}
