package ovh

import (
	"fmt"
	"os"
	"slices"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Worker struct {
	Session *session.Session
	S3      *s3.S3
}

var (
	blogBucket    = os.Getenv("AWS_BLOG_BUCKET_NAME")
	galleryBucket = os.Getenv("AWS_GALLERY_BUCKET_NAME")
)

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

	// Check if buckets exist
	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	bucketNames := make([]string, len(result.Buckets))

	for _, b := range result.Buckets {
		bucketNames = append(bucketNames, *b.Name)
	}

	if !slices.Contains(bucketNames, blogBucket) {
		return nil, fmt.Errorf("bucket %s doesn't exist on S3", blogBucket)
	}

	if !slices.Contains(bucketNames, galleryBucket) {
		return nil, fmt.Errorf("bucket %s doesn't exist on S3", galleryBucket)
	}

	return &Worker{
		Session: s,
		S3:      s3Client,
	}, nil
}
