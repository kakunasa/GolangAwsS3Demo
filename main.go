package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"goAwsS3/s3Actions"
	"goAwsS3/services"
	"net/http"
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

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		dir := "tempAsset"
		err = os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(dir)
		fileMap := services.UploadHandler(w, r)
		for fileName, fileByte := range fileMap {
			// 将 byte 数据写入文件
			err := os.WriteFile(filepath.Join("tempAsset", fileName), fileByte[1:], 0644)
			if err != nil {
				fmt.Println(err)
				return
			}

			// 打开文件
			file, err := os.Open("tempAsset/" + fileName)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 上传文件到 S3
			err = bucketBasics.UploadFile(os.Getenv("AWS_BUCKET_NAME"), fileName, file)
			if err != nil {
				fmt.Printf("Couldn't upload your file. Here's why: %v\n", err)
				return
			} else {
				fmt.Printf("Upload file: %v success\n", fileName)
			}
		}
	})

	fmt.Printf("Server is running on :%v...\n", os.Getenv("SERVER_LISTEN_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("SERVER_LISTEN_PORT"), nil)
	if err != nil {
		fmt.Printf("Error is:%v\n", err)
		return
	}
}
