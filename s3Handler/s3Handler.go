package s3handler

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Handler stores information for accessing
// a S3 bucket in a specified region.
// Defines behaviour for uploading and downloading files from S3 buckets.
type S3Handler struct {
	awsRegion  string
	bucketname string
}

// NewHandler returns a new S3 Handler for the specified bucket.
func NewHandler(awsRegion string, bucketname string) S3Handler {
	handler := new(S3Handler)
	handler.awsRegion = awsRegion
	handler.bucketname = bucketname
	return *handler
}

// UploadFileToBucket uploads a file to the bucket specified by the S3Handler.
// Uploaded file is given the specified filename in the S3 bucket.
func (s *S3Handler) UploadFileToBucket(filename string, filecontent *os.File, sess session.Session) {
	uploader := s3manager.NewUploader(&sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucketname),
		Key:    aws.String(filename),
		Body:   filecontent,
	})
	s.checkError(err)
	log.Printf("Uploaded %q to %q\n", filename, s.bucketname)
}

// DownloadFileFromBucket downloads a file with the specified name from
// the bucket specified by the S3Handler.
func (s *S3Handler) DownloadFileFromBucket(downloadTo *os.File, toDownload string, sess *session.Session) {
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(downloadTo,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketname),
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
