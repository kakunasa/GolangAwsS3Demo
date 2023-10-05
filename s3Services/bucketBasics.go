package s3Services

import "github.com/aws/aws-sdk-go-v2/service/s3"

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon s3Services) actions
// used in the examples.
// It contains S3Client, an Amazon s3Services service client that is used to perform bucket
// and object actions.
type BucketBasics struct {
	S3Client *s3.Client
}
