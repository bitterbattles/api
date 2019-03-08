package dynamo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewClient creates a new DynamoDB client
func NewClient(session *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(session)
}
