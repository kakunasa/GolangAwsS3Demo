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
