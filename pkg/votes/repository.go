package votes

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "votes"

// RepositoryInterface defines an interface for a Vote repository
type RepositoryInterface interface {
	Add(Vote) error
}

// Repository is an implementation of RepositoryInterface using DynamoDB
type Repository struct {
	client *dynamodb.DynamoDB
}

// NewRepository creates a new Votes repository instance
func NewRepository(client *dynamodb.DynamoDB) *Repository {
	return &Repository{client}
}

// Add is used to insert a new Vote document
func (repository *Repository) Add(vote Vote) error {
	item, err := dynamodbattribute.MarshalMap(vote)
	if err != nil {
		return err
	}
	_, err = repository.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	return err
}
