package s3handler

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Handler struct {
	awsRegion  string
	bucketname string
}

func NewHandler(awsRegion string, bucketname string) S3Handler {
	handler := new(S3Handler)
	handler.awsRegion = awsRegion
	handler.bucketname = bucketname
	return *handler
}

func (s *S3Handler) UploadFileToBucket(bucket string, filename string, filecontent *os.File, sess session.Session) {
	uploader := s3manager.NewUploader(&sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   filecontent,
	})
	s.checkError(err)
	log.Printf("Uploaded %q to %q\n", filename, bucket)
}

func (s *S3Handler) DownloadFileFromBucket(downloadTo *os.File, bucketName string, toDownload string, sess *session.Session) {
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(downloadTo,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(toDownload),
		})
	s.checkError(err)
	log.Println("Downloadaded ", downloadTo.Name(), " ", numBytes, " bytes")
}

func (s *S3Handler) checkError(err error) {
	if err != nil {
		log.Println("Error in s3Handler.go: ")
		log.Fatal(err)
	}
}
