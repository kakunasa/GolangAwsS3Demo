package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"goAwsS3/s3Actions"
	"os"
	"path/filepath"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func main() {
	loadEnv()
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	bucketBasics := s3Actions.BucketBasics{S3Client: s3Client}

	err = bucketBasics.ListBuckets()
	if err != nil {
		fmt.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
		return
	}

	fileNameToFullNameMap := scannerAssetDir()

	for fileName, filePath := range fileNameToFullNameMap {
		fmt.Println(fileName)
		err = bucketBasics.UploadFile("kaku-golang-bucket", fileName, filePath)
		if err != nil {
			fmt.Printf("Couldn't upload your file. Here's why: %v\n", err)
			return
		}
	}
}

func scannerAssetDir() map[string]string {
	// Get the asset folder in the root directory
	assetDir, err := filepath.Abs("asset")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Iterate through all the files in the asset folder
	files, err := os.ReadDir(assetDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Get the file name and path
	fileNmToFullNmMap := make(map[string]string)
	for _, file := range files {
		// Filter out folder
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		filePath := filepath.Join(assetDir, fileName)

		fileNmToFullNmMap[fileName] = filePath
	}

	return fileNmToFullNmMap
}
