package main

import (
	"context"
	"encoding/json"
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

type FileInfo struct {
	FileName string `json:"fileName"`
	FileUrl  string `json:"fileUrl"`
}

func main() {
	// Load the environment variables.
	loadEnv()

	// Create a new AWS SDK configuration.
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}

	// Create a new S3 client.
	s3Client := s3.NewFromConfig(sdkConfig)

	// Create a new bucket basics object.
	bucketBasics := s3Actions.BucketBasics{S3Client: s3Client}

	presignClient := s3.NewPresignClient(s3Client)
	presigner := s3Actions.Presigner{PresignClient: presignClient}

	// Create a new HTTP handler for the `/upload` endpoint.
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Create a temporary directory to store the uploaded files.
		dir := "tempAsset"
		err = os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Defer removing the temporary directory.
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(dir)

		// Get the uploaded files from the request.
		fileMap := services.UploadHandler(w, r)

		// Return the uploaded files url to request.
		fileList := make([]FileInfo, 0)

		// Iterate over the uploaded files and upload them to S3.
		for fileName, fileByte := range fileMap {
			// Write the byte data to a file.
			err := os.WriteFile(filepath.Join("tempAsset", fileName), fileByte[1:], 0644)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Open the file.
			file, err := os.Open("tempAsset/" + fileName)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Upload the file to S3.
			err = bucketBasics.UploadFile(os.Getenv("AWS_BUCKET_NAME"), fileName, file)
			if err != nil {
				fmt.Printf("Couldn't upload your file. Here's why: %v\n", err)
				return
			} else {
				fmt.Printf("Upload file: %v success\n", fileName)
				presignedGetRequest, err := presigner.GetObject(os.Getenv("AWS_BUCKET_NAME"), fileName, 60)
				if err != nil {
					fmt.Printf("Couldn't get file presign. Here's why: %v\n", err)
					return
				}
				fileList = append(fileList, FileInfo{
					FileName: fileName,
					FileUrl:  presignedGetRequest.URL,
				})
			}
		}

		// Set the response header
		w.Header().Set("Content-Type", "application/json")

		// Create a JSON response object.
		response := struct {
			Data   []FileInfo `json:"data"`
			Status string     `json:"status"`
		}{
			Data:   fileList,
			Status: "OK",
		}

		// Marshal the JSON response object into a JSON string.
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write the JSON response to the response writer.
		_, err = w.Write(jsonResponse)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	// Start the HTTP server.
	fmt.Printf("Server is running on :%v...\n", os.Getenv("SERVER_LISTEN_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("SERVER_LISTEN_PORT"), nil)
	if err != nil {
		fmt.Printf("Error is:%v\n", err)
		return
	}
}
