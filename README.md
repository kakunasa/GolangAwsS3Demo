# GolangAwsS3Demo

### Project Description
Use GoLang to configure AWS SDK and control file upload, file download and file deletion in S3 bucket

### Pre-requirements
+ Use AWS IAM roles generate access key ID and secret access key，then set corresponding configuration items in the .env file
+ reference: <https://aws.github.io/aws-sdk-go-v2/docs/>

### Project Initialization
```bash
git clone https://github.com/kakunasa/GolangAwsS3Demo.git
go mod init GolangAwsS3Demo
go get github.com/joho/godotenv
go get github.com/aws/aws-sdk-go-v2
go mod tidy
go run .
```

### Test Api ~/upload
Start this project,after see
```bash
Server is running on :<SERVER_LISTEN_PORT>...
```
run the command line below
```bash
curl -F "files=@testData/send_request/csvNo1.csv" \
-F "files=@testData/send_request/jpgNo1.jpg" \
-F "files=@testData/send_request/jpgNo2.jpg" \
-F "files=@testData/send_request/txtNo1.txt" \
-F "files=@testData/send_request/txtNo2.txt" \
http://localhost:18088/upload
```

### Api ~/upload output
```json
{
  "data": {
    "fileList": [
      {
        "fileName": "test-file1.txt",
        "fileUrl": "https://s3.amazonaws.com/my-bucket/test-file1.txt"
      },
      {
        "fileName": "test-file2.jpg",
        "fileUrl": "https://s3.amazonaws.com/my-bucket/test-file2.jpg"
      }
    ]
  },
  "status": "OK"
}
```