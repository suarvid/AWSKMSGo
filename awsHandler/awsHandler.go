package awshandler

import (
	kmshandler "KMSClient/KMSHandler"
	FileHandler "KMSClient/fileHandler"
	s3handler "KMSClient/s3Handler"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AwsHandler holds handlers for the individual AWS Services and
// a path to a specific bucket to use for uploading and downloading files,
// and to a key for encrypting and decrypting.
// Defines behaviour for encrypting and then uploading files
// or decrypting and then downloading files.
type AwsHandler struct {
	awsRegion   string
	keyARN      string
	bucketname  string
	s3Handler   s3handler.S3Handler
	kmsHandler  kmshandler.KMSHandler
	fileHandler FileHandler.FileHandler
}

// NewHandler returns a new AwsHandler with paths to the
// specified s3 bucket and KMS key.
// Individual handlers of Aws Handler also have the same paths.
func NewHandler(awsRegion string, keyARN string, bucketname string, fileHandler FileHandler.FileHandler) AwsHandler {
	handler := new(AwsHandler)
	handler.awsRegion = awsRegion
	handler.keyARN = keyARN
	handler.bucketname = bucketname
	handler.fileHandler = fileHandler
	handler.s3Handler = s3handler.NewHandler(awsRegion, bucketname)
	handler.kmsHandler = kmshandler.NewHandler(awsRegion, keyARN)
	return *handler
}

// EncryptUpload encrypts and uploads a file.
// Which file is encrypted and what it is named when uploaded is defined
// by the FileHandler.
func (a *AwsHandler) EncryptUpload(keyARN string, bucketname string) {
	sess := a.CreateSession()
	toEncrypt := a.fileHandler.ReadFile(a.fileHandler.GetPlaintextPath())
	encryptResult := a.kmsHandler.EncryptData(toEncrypt, sess)
	a.fileHandler.WriteFile(a.fileHandler.GetEncryptedPath(), encryptResult.CiphertextBlob)
	log.Println("File Encrypted")
	toUploadhandle := a.fileHandler.GetFileHandle(a.fileHandler.GetEncryptedPath())
	a.s3Handler.UploadFileToBucket(a.fileHandler.GetPlaintextPath(), toUploadhandle, sess)
	log.Println("File uploaded")

}

// DownloadDecrypt downloads a file and tries to decrypt it.
// Key used for decryption is defined by the kmsHandler, path to
// file for downloading is defined by the S3 Handler.
func (a *AwsHandler) DownloadDecrypt(filename string, bucketname string) {
	sess := a.CreateSession()
	fileHandler := a.fileHandler
	downloadToHandle := a.fileHandler.GetFileHandle(a.fileHandler.GetDownloadPath())
	a.s3Handler.DownloadFileFromBucket(downloadToHandle, filename, &sess)
	log.Println("File Downloaded")
	downloadedFile := fileHandler.ReadFile(downloadToHandle.Name())
	decryptResult := a.kmsHandler.DecryptData(downloadedFile, sess)
	fileHandler.WriteFile(fileHandler.GetDecryptedPath(), decryptResult.Plaintext)
}

// CreateSession should probably not really be exported,
// but I did so for ease of use when testing.
func (a *AwsHandler) CreateSession() session.Session {
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
