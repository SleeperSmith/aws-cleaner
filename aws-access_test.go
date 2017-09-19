package main

import (
	"testing"

	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Mock struct {
}

func (mock *s3Mock) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return &s3.ListBucketsOutput{}, nil
}
func (mock *s3Mock) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return nil, nil
}
func (mock *s3Mock) ListObjectVersions(input *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error) {
	return nil, nil
}
func (mock *s3Mock) DeleteObjects(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	return nil, nil
}
func (mock *s3Mock) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return nil, nil
}

func TestGetS3Buckets(t *testing.T) {
	s3client = new(s3Mock)

	result := getBuckets()

	if result == nil {
		t.Fail()
	}

	fmt.Printf("length %d", len(result))
	for _, bucket := range result {
		fmt.Printf("bucket %v", bucket)
	}
}
