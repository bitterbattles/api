package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewDynamoClient creates a new DynamoDB client
func NewDynamoClient(session *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(session)
}
