package s3Services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ListBuckets lists the buckets in the current account.
func (basics BucketBasics) ListBuckets() error {
	result, err := basics.S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	count := 10
	fmt.Printf("Let's list up to %v buckets for your account.\n", count)
	if err != nil {
		fmt.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
		return err
	}
	if len(result.Buckets) == 0 {
		fmt.Println("You don't have any buckets!")
		return nil
	} else {
		if count > len(result.Buckets) {
			count = len(result.Buckets)
		}
		for _, bucket := range result.Buckets[:count] {
			fmt.Printf("\t%v\n", *bucket.Name)
		}
		return nil
	}
}
