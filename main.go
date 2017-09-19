package main

import (
	"fmt"
	"os"

	"sync"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	buckets := getBuckets()

	wg := new(sync.WaitGroup)
	wg.Add(len(buckets))
	for _, bucket := range buckets {
		fmt.Printf("Checking bucket: %v\n", *bucket.Name)
		go func(bucket *s3.Bucket) {
			func(bucket *s3.Bucket) {
				stackID := getBucketStackID(bucket.Name)
				if stackID == nil {
					return
				}

				stack := getStack(stackID)
				if stack != nil && *stack.StackStatus != cloudformation.StackStatusDeleteComplete {
					return
				}

				fmt.Println("========================================")
				fmt.Printf("Deleting bucket: %v\n", *bucket.Name)
				fmt.Printf("Stack ID: %v\n", *stackID)
				if stack != nil {
					fmt.Printf("Stack Status: %v.\n", *stack.StackStatus)
				} else {
					fmt.Println("Stack no longer exist.")
				}
				versions, err := getAllS3Objects(bucket.Name)
				if err != nil {
					return
				}
				versionIdentifiers := make([]*s3.ObjectIdentifier, len(versions))
				for i, objVersion := range versions {
					versionIdentifiers[i] = &s3.ObjectIdentifier{
						Key:       objVersion.Key,
						VersionId: objVersion.VersionId,
					}
				}
				deleteS3Object(bucket.Name, versionIdentifiers)
				deleteBucket(bucket.Name)
			}(bucket)
			wg.Done()
		}(bucket)
	}
	wg.Wait()
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
