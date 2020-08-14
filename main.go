package main

import (
	kmshandler "KMSClient/KMSHandler"
	awshandler "KMSClient/awsHandler"
	filehandler "KMSClient/fileHandler"
	"os"
)

// Constant filepaths for ease-of-use when just testing
var (
	plainTextPath string = "/home/arvid/go/src/KMSClients/AWS/GoKMSClient/files/testfile.json"
	encryptedPath string = "/home/arvid/go/src/KMSClients/AWS/GoKMSClient/files/encrypted"
	downloadPath  string = "/home/arvid/go/src/KMSClients/AWS/GoKMSClient/files/downloaded"
	decryptedPath string = "/home/arvid/go/src/KMSClients/AWS/GoKMSClient/files/decrypted.json"
	policyPath    string = "/home/arvid/go/src/KMSClients/AWS/GoKMSClient/files/examplepolicy.json"
)

// Encrypt a sample json-file, upload it to S3
// Then Download and decrypt the uploaded file, writing to decrypted.json
func main() {
	keyARN := os.Getenv("GO_AWS_KEY_ARN")
	bucketname := os.Getenv("GO_AWS_BUCKET_NAME")
	AWSRegion := os.Getenv("GO_AWS_REGION")
	fileHandler := filehandler.NewHandler(plainTextPath, encryptedPath, downloadPath, decryptedPath)
	awsHandler := awshandler.NewHandler(AWSRegion, keyARN, bucketname, fileHandler)
	sess := awsHandler.CreateSession()
	kmsHandler := kmshandler.NewHandler(AWSRegion, keyARN)
	tags := make(map[string]string)
	tags["CreatedBy"] = "suarvid"
	tags["Alias"] = "IsThisTheName"
	// kmsHandler.CreateKey(tags, sess)
}
