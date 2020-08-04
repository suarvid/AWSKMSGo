package KMSHandler

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMSHandler struct {
	awsRegion string
	keyARN    string
}

func NewHandler(awsRegion string, keyARN string) KMSHandler {
	handler := new(KMSHandler)
	handler.awsRegion = awsRegion
	handler.keyARN = keyARN
	return *handler
}

func (k *KMSHandler) CreateKey(sess session.Session) *kms.CreateKeyOutput {
	service := kms.New(&sess)
	result, err := service.CreateKey(&kms.CreateKeyInput{
		Tags: []*kms.Tag{
			{
				TagKey:   aws.String("CreatedBy"),
				TagValue: aws.String("suarvid"),
			},
		},
	})
	k.checkError(err)
	return result
}

func (k *KMSHandler) EncryptData(data []byte, sess session.Session) *kms.EncryptOutput {
	service := kms.New(&sess)
	result, err := service.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(k.keyARN),
		Plaintext: data,
	})
	k.checkError(err)
	return result
}

func (k *KMSHandler) DecryptData(data []byte, sess session.Session) *kms.DecryptOutput {
	service := kms.New(&sess)
	result, err := service.Decrypt(&kms.DecryptInput{
		CiphertextBlob: data,
	})
	k.checkError(err)
	return result
}

func (k *KMSHandler) checkError(err error) {
	if err != nil {
		log.Println("Error in kmsHandler.go: ")
		log.Fatal(err)
	}
}
