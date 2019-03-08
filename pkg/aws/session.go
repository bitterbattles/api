package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates a new AWS session
func NewSession(region string) *session.Session {
	config := aws.Config{
		Region: aws.String(region),
	}
	session, err := session.NewSession(&config)
	if err != nil {
		panic(err)
	}
	return session
}
