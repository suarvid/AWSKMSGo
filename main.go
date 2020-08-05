package main

import (
	AwsHandler "KMSClient/awsHandler"
	FileHandler "KMSClient/fileHandler"
	"fmt"
	"os"
)

// Constant filepaths for ease-of-use when just testing
var (
	plainTextPath string = "/home/arvid/go/src/GoKMSClient/files/testfile.json"
	encryptedPath string = "/home/arvid/go/src/GoKMSClient/files/encrypted"
	downloadPath  string = "/home/arvid/go/src/GoKMSClient/files/downloaded"
	decryptedPath string = "/home/arvid/go/src/GoKMSClient/files/decrypted.json"
)

// Encrypt a sample json-file, upload it to S3
// Then Download and decrypt the uploaded file, writing to decrypted.json
func main() {
	keyARN := os.Getenv("GO_AWS_KEY_ARN")
	bucketname := os.Getenv("GO_AWS_BUCKET_NAME")
	AWSRegion := os.Getenv("GO_AWS_REGION")
	filehandler := FileHandler.NewHandler(plainTextPath, encryptedPath, downloadPath, decryptedPath)
	awsHandler := AwsHandler.NewHandler(AWSRegion, keyARN, bucketname, filehandler)

	fmt.Printf("Key ARN: %s\n", keyARN)
	fmt.Printf("Bucket Name: %s\n", bucketname)

	awsHandler.EncryptUpload(keyARN, bucketname)
	awsHandler.DownloadDecrypt(plainTextPath, bucketname)
}
