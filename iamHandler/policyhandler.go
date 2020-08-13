package iamhandler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAMPolicyHandler defines methods for managing
// IAM Policies in a certain region.
type IAMPolicyHandler struct {
	region string
}

// NewPolicyHandler creates a new policy handler for
// the specified region.
func NewPolicyHandler(region string) *IAMPolicyHandler {
	handler := new(IAMPolicyHandler)
	handler.region = region
	return handler
}

// PolicyDocument represents the info necessary to construct
// a new IAM Policy through the PolicyHandler.
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// StatementEntry contains information for statements
// which are used to construct an IAM Policy.
type StatementEntry struct {
	Effect   string
	Action   []string
	Resource string
}

// BuildStatementEntry returns a StatementEntry with the specified
// effect and actions on the given resource.
// Mainly improves QOL and limits amount of formatting necessary to create a policy.
func (p *IAMPolicyHandler) BuildStatementEntry(effect string, action []string, resource string) *StatementEntry {
	entry := StatementEntry{
		Effect:   effect,
		Action:   action,
		Resource: resource,
	}
	return &entry
}

// CreatePolicy creates an IAM Policy with the given name and policy version.
// Statement Entries that make up the policy should be created beforehand by
// using BuildStatementEntry. Policy names must be unique, cannot update an already
// existing policy by creating a new policy with the same name.
// Recommended version to use is 2012-10-17 as it is the latest version.
func (p *IAMPolicyHandler) CreatePolicy(policyName, version string, statements []StatementEntry) *iam.CreatePolicyOutput {
	service := p.createIAMService()

	policy := PolicyDocument{
		Version:   version,
		Statement: statements,
	}
	bytes, err := json.Marshal(&policy)
	if err != nil {
		log.Fatalf("Error marshalling policy to JSON: %v", err)
	}

	result, err := service.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(bytes)),
		PolicyName:     aws.String(policyName),
	})

	if err != nil {
		log.Fatalf("Error creating policy %s: %v", policyName, err)
	}

	fmt.Println("Successfully created new policy", policyName)
	return result
}

// DeletePolicy deletes the policy with the given ARN to access it.
func (p *IAMPolicyHandler) DeletePolicy(policyARN string) *iam.DeletePolicyOutput {
	service := p.createIAMService()
	result, err := service.DeletePolicy(&iam.DeletePolicyInput{
		PolicyArn: &policyARN,
	})
	if err != nil {
		log.Fatalf("Error deleting policy: %v", err)
	}
	fmt.Printf("Successfully deleted policy with ARN %s", policyARN)
	return result
}

// GetPolicy returns an object representing the IAM Policy
// with the given ARN.
func (p *IAMPolicyHandler) GetPolicy(ARN string) iam.Policy {
	service := p.createIAMService()
	result, err := service.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: &ARN,
	})

	if err != nil {
		log.Fatalf("Error getting IAM policy: %v", err)
	}

	return *result.Policy
}

func (p *IAMPolicyHandler) createIAMService() *iam.IAM {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(p.region)},
	)
	if err != nil {
		log.Fatalf("Error creating IAM service: %v", err)
	}

	return iam.New(sess)
}
