package iamhandler

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAMUserHandler defines methods for CRUD operations
// on IAM Users as well as defining access keys for them.
type IAMUserHandler struct {
	region string
}

// NewUserHandler returns an IAMHandler with the
// region given as a parameter. Holds a policymanager
// with the same region internally.
func NewUserHandler(region string) *IAMUserHandler {
	handler := new(IAMUserHandler)
	handler.region = region
	return handler
}

// GetUser returns an object representing an IAM User
// with the given username.
func (i *IAMUserHandler) GetUser(username string) iam.User {
	service := i.createIAMService()

	result, err := service.GetUser(&iam.GetUserInput{
		UserName: &username,
	})

	if err != nil {
		log.Fatalf("Error getting IAM user: %v", err)
	}
	return *result.User
}

// CreateUser creates an IAM User with the given username.
func (i *IAMUserHandler) CreateUser(username string) iam.User {
	service := i.createIAMService()

	result, err := service.CreateUser(&iam.CreateUserInput{
		UserName: &username,
	})

	if err != nil {
		log.Fatalf("Error creating IAM user: %v", err)
	}

	return *result.User
}

// DeleteUser tries to delete an IAM User with the given
// username. Fatal error is thrown if the user doesn't exist.
func (i *IAMUserHandler) DeleteUser(username string) {
	service := i.createIAMService()

	_, err := service.DeleteUser(&iam.DeleteUserInput{
		UserName: &username,
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		log.Fatalf("Error: User %s does not exist", username)
	} else if err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}

	fmt.Printf("Deleted user %s\n", username)
}

// ListIAMUsers lists up to 15 IAM Users
// in the region defined by the IAMHandler.
func (i *IAMUserHandler) ListIAMUsers() {
	service := i.createIAMService()

	result, err := service.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(15),
	})

	if err != nil {
		log.Fatalf("Error listing IAM Users: %v", err)
	}

	for index, user := range result.Users {
		fmt.Printf("%d User %s, date created: %v\n", index, *user.UserName, user.CreateDate)
	}
}

// CreateAccessKey creates an access key for signing programmatic requests
// to the AWS API. Takes username of access key owner as parameter.
func (i *IAMUserHandler) CreateAccessKey(owner string) {
	service := i.createIAMService()

	result, err := service.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(owner),
	})

	if err != nil {
		log.Fatalf("Error creating Access Key: %v", err)
	}

	fmt.Println("Created Access Key: ", *result.AccessKey)
}

func (i *IAMUserHandler) createIAMService() iam.IAM {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(i.region)},
	)

	if err != nil {
		log.Fatalf("Error getting IAM service: %v", err)
	}

	return *iam.New(sess)
}
