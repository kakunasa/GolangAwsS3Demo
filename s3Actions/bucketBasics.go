package s3Actions

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon s3Actions) actions
// used in the examples.
// It contains S3Client, an Amazon s3Actions service client that is used to perform bucket
// and object actions.
type BucketBasics struct {
	S3Client *s3.Client
}

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

// UploadFile reads from a file and puts the data into an object in a bucket.
// This function is a member of the BucketBasics structure.
func (basics BucketBasics) UploadFile(bucketName string, objectKey string, file *os.File) error {
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Couldn't close file. Here's why: %v\n", err)
		}
	}(file)
	_, err := basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v. Here's why: %v\n",
			objectKey, bucketName, err)
	}

	return nil
}
