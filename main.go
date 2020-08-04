package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	toEncryptPath := "./testFile.json"
	toDecryptPath := "./encrypted"
	downloadPath := "./downloaded.txt"
	decryptedPath := "./decrypted.json"
	keyARN := os.Args[1]
	bucketName := os.Args[2]

	fmt.Printf("Key ARN: %s\n", keyARN)
	fmt.Printf("Bucket Name: %s\n", bucketName)

	sess := createSession()

	toEncrypt := readFile(toEncryptPath)
	encryptResult := encryptData(toEncrypt, sess, keyARN)
	writeFile(toDecryptPath, encryptResult.CiphertextBlob)
	fmt.Println("File encrypted")
	fileToUploadHandle := getFileHandle(toDecryptPath)
	uploadFileToBucket(bucketName, toDecryptPath, fileToUploadHandle, sess)
	downloadToHandle := getFileHandle(downloadPath)
	downloadFileFromBucket(downloadToHandle, bucketName, toDecryptPath, &sess)
	downloadedFile := readFile(downloadToHandle.Name())
	decryptResult := decryptData(downloadedFile, sess)
	fmt.Println("File Decrypted")
	fmt.Println(string(decryptResult.Plaintext))
	writeFile(decryptedPath, decryptResult.Plaintext)
}

// creates session for accessing AWS
func createSession() session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})

	handleError(err)
	return *sess
}

// creates AWS KMS Key with specified tags
// can access metadata such as ARN or KeyID through result
func createKey(sess session.Session) *kms.CreateKeyOutput {
	service := kms.New(&sess)
	result, err := service.CreateKey(&kms.CreateKeyInput{
		Tags: []*kms.Tag{
			{
				TagKey:   aws.String("CreatedBy"),
				TagValue: aws.String("suarvid"),
			},
		},
	})
	handleError(err)
	return result
}

// encrypts the given data using the key with the supplied ARN
func encryptData(data []byte, sess session.Session, keyARN string) *kms.EncryptOutput {
	service := kms.New(&sess)
	result, err := service.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyARN),
		Plaintext: []byte(data),
	})
	handleError(err)
	return result
}

// Decrypts the given data with the given session
// have to call string() on the plaintext of the result for it to be readable
func decryptData(data []byte, sess session.Session) *kms.DecryptOutput {
	service := kms.New(&sess)
	result, err := service.Decrypt(&kms.DecryptInput{
		CiphertextBlob: data,
	})
	handleError(err)
	return result
}

func displayS3Buckets(sess *session.Session) {
	service := s3.New(sess)
	result, err := service.ListBuckets(nil)
	handleError(err)
	fmt.Println("Buckets: ")
	for _, bucket := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(bucket.Name), aws.TimeValue(bucket.CreationDate))
	}
}

func uploadFileToBucket(bucket string, filename string, filecontent *os.File, sess session.Session) {
	uploader := s3manager.NewUploader(&sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   filecontent,
	})
	handleError(err)
	fmt.Printf("Uploaded %q to %q\n", filename, bucket)
}

// Download file with specified ID from bucket with specified name
// Writes downloaded file to the provided *File
func downloadFileFromBucket(downloadTo *os.File, bucketName string, toDownload string, sess *session.Session) {
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(downloadTo,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(toDownload),
		})
	handleError(err)
	fmt.Println("Downloadaded ", downloadTo.Name(), " ", numBytes, " bytes")
}

// returns a handle to a file from which its data can be read
// Handle is opened for both reading and writing
func getFileHandle(path string) *os.File {
	if !fileExists(path) {
		os.Create(path)
	}
	fileHandle, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error getting handle for file %s ", path)
		log.Fatal(err)
	}
	return fileHandle
}

// func readFileWithHandle(fileHandle *os.File) []byte {
// 	defer fileHandle.Close()
// 	readIncrement := make([]byte, 10)
// 	fileHandle.
// 	content, err := fileHandle.Read(readIncrement)
// 	handleError(err)
// 	return content
// }

// directly reads the contents of a file
func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	handleError(err)
	return data
}

func writeFile(path string, data []byte) {
	file, err := os.Create(path)
	handleError(err)
	defer file.Close()
	file.Write(data)
}

// Checks if specified file exists
// Cannot be a directory
func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// panics and logs error information
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
