package AwsHandler

import (
	"KMSClient/KMSHandler"
	FileHandler "KMSClient/fileHandler"
	s3handler "KMSClient/s3Handler"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsHandler struct {
	awsRegion   string
	keyARN      string
	bucketname  string
	s3Handler   s3handler.S3Handler
	kmsHandler  KMSHandler.KMSHandler
	fileHandler FileHandler.FileHandler
}

func NewHandler(awsRegion string, keyARN string, bucketname string, fileHandler FileHandler.FileHandler) AwsHandler {
	handler := new(AwsHandler)
	handler.awsRegion = awsRegion
	handler.keyARN = keyARN
	handler.bucketname = bucketname
	handler.fileHandler = fileHandler
	handler.s3Handler = s3handler.NewHandler(awsRegion, bucketname)
	handler.kmsHandler = KMSHandler.NewHandler(awsRegion, keyARN)
	return *handler
}

func (a *AwsHandler) EncryptUpload(keyARN string, bucketname string) {
	sess := a.createSession()
	toEncrypt := a.fileHandler.ReadFile(a.fileHandler.GetPlaintextPath())
	encryptResult := a.kmsHandler.EncryptData(toEncrypt, sess)
	a.fileHandler.WriteFile(a.fileHandler.GetEncryptedPath(), encryptResult.CiphertextBlob)
	log.Println("File Encrypted")
	toUploadhandle := a.fileHandler.GetFileHandle(a.fileHandler.GetEncryptedPath())
	a.s3Handler.UploadFileToBucket(bucketname, a.fileHandler.GetPlaintextPath(), toUploadhandle, sess)
	log.Println("File uploaded")

}

func (a *AwsHandler) DownloadDecrypt(filename string, bucketname string) {
	sess := a.createSession()
	fileHandler := a.fileHandler
	downloadToHandle := a.fileHandler.GetFileHandle(a.fileHandler.GetDownloadPath())
	a.s3Handler.DownloadFileFromBucket(downloadToHandle, bucketname, filename, &sess)
	log.Println("File Downloaded")
	downloadedFile := fileHandler.ReadFile(downloadToHandle.Name())
	decryptResult := a.kmsHandler.DecryptData(downloadedFile, sess)
	fileHandler.WriteFile(fileHandler.GetDecryptedPath(), decryptResult.Plaintext)
}

func (a *AwsHandler) createSession() session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.awsRegion),
	})

	a.checkError(err)

	return *sess
}

func (a *AwsHandler) checkError(err error) {
	if err != nil {
		log.Println("Error using AWS: ")
		log.Fatal(err)
	}
}
