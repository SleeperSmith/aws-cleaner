package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Interface interface {
	ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error)
	ListObjectVersions(input *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error)
	DeleteObjects(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)
	DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)
}

var sess = session.Must(session.NewSession())
var s3client s3Interface = s3.New(sess)
var cfnClient = cloudformation.New(sess)

func getBuckets() []*s3.Bucket {
	result, err := s3client.ListBuckets(nil)

	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	return result.Buckets
}

func getBucketStackID(bucketName *string) *string {
	result, err := s3client.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: bucketName,
	})

	if err != nil {
		return nil
	}

	for _, ts := range result.TagSet {
		if *ts.Key == "aws:cloudformation:stack-id" {
			return ts.Value
		}
	}

	return nil
}

func getStack(stackID *string) *cloudformation.Stack {
	result, err := cfnClient.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: stackID,
	})
	if err != nil {
		return nil
	}

	if len(result.Stacks) != 1 {
		return nil
	}

	return result.Stacks[0]
}

func getAllS3Objects(bucketName *string) ([]*s3.ObjectVersion, error) {
	result, err := s3client.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, err
	}

	return result.Versions, nil
}

func deleteS3Object(bucketName *string, objToDelete []*s3.ObjectIdentifier) {
	s3client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: bucketName,
		Delete: &s3.Delete{
			Objects: objToDelete,
		},
	})
}

func deleteBucket(bucketName *string) {
	s3client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: bucketName,
	})
}
