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
	sess := createSession()
	displayS3Buckets(&sess)
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
				TagValue: aws.String("ExampleUser"),
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

func uploadFileToBucket(bucket string, filename string, filecontent []byte, sess session.Session) {
	uploader := s3manager.NewUploader(&sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   filecontent,
	})
}

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

// panics and logs error information
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
