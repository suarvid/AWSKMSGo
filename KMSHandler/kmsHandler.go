package kmshandler

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// KMSHandler stores information for accessing
// and using a KMS-Key for encryption and decryption.
type KMSHandler struct {
	awsRegion string
	keyARN    string
}

// NewHandler returns a KMSHandler with the specified
// aws Region and ARN for the key to use.
func NewHandler(awsRegion string, keyARN string) KMSHandler {
	handler := new(KMSHandler)
	handler.awsRegion = awsRegion
	handler.keyARN = keyARN
	return *handler
}

// CreateKey creates a key in the aws region defined
// by the KMS Handler. Is possible to set more extensive
// createKeyInput than just tags.
func (k *KMSHandler) CreateKey(tags map[string]string, sess session.Session) *kms.CreateKeyOutput {
	service := kms.New(&sess)
	awsTags := convertKeyTags(tags)
	result, err := service.CreateKey(&kms.CreateKeyInput{
		Tags: *awsTags,
	})
	k.checkError(err)
	return result
}

// QOL function for converting map of tags to
// aws-types without breaking out each tag.
// Double pointers are a bit icky, should probably be fixed.
func convertKeyTags(tags map[string]string) *[]*kms.Tag {
	var Tags []*kms.Tag
	for _, tag := range tags {
		Tags = append(Tags, &kms.Tag{
			TagKey:   aws.String(tag),
			TagValue: aws.String(tags[tag]),
		})
	}
	return &Tags
}

// DisableKey disables the key with the ARN specified
// by the KMSHandler and prevents it from being used for encryption/decryption.
// Closest thing possible for deleting an individual key without deleting the entire
// Custom Key Store.
func (k *KMSHandler) DisableKey(sess session.Session) {
	service := kms.New(&sess)
	_, err := service.DisableKey(&kms.DisableKeyInput{
		KeyId: &k.keyARN,
	})
	if err != nil {
		log.Fatalf("Error disabling key %s: %v", k.keyARN, err)
	}

	fmt.Printf("Successfully deleted key with ARN %s", k.keyARN)
}

// EncryptData encrypts data given as a byte slice.
// Uses the key with the ARN defined by the KMS Handler for encryption.
func (k *KMSHandler) EncryptData(data []byte, sess session.Session) *kms.EncryptOutput {
	service := kms.New(&sess)
	result, err := service.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(k.keyARN),
		Plaintext: data,
	})
	k.checkError(err)
	return result
}

// DecryptData decrypts data given as a byte slice.
// Uses the key with the ARN defined by the KMSHandler.
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
