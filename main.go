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
		fileMap := services.UploadHandler(w, r)

		dataChan := make(chan []byte)

		go func() {
			for _, data := range fileMap {
				dataChan <- data
			}
			close(dataChan)
		}()

		go func() {
			for data := range dataChan {
				// 创建一个临时文件
				file, err := os.CreateTemp("", "file")
				if err != nil {
					fmt.Println(err)
					return
				}

				// 将 []byte 写入临时文件
				err = os.WriteFile(file.Name(), data, 0644)
				if err != nil {
					fmt.Println(err)
					return
				}

				// 将临时文件转换为 *os.File
				file, err = os.OpenFile(file.Name(), os.O_RDONLY, 0644)
				if err != nil {
					fmt.Println(err)
					return
				}

				fileName := file.Name()

				// 上传文件到 S3
				err = bucketBasics.UploadFile("kaku-golang-bucket", fileName, file)
				if err != nil {
					fmt.Printf("Couldn't upload your file. Here's why: %v\n", err)
					return
				}
			}
		}()

	})

	fmt.Printf("Server is running on :%v...\n", os.Getenv("SERVER_LISTEN_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("SERVER_LISTEN_PORT"), nil)
	if err != nil {
		fmt.Printf("Error is:%v\n", err)
		return
	}
}
