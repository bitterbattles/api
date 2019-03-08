package users

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName     = "users"
	usernameIndex = "username"
)

// RepositoryInterface defines an interface for a User repository
type RepositoryInterface interface {
	Add(*User) error
	GetByUsername(string) (*User, error)
}

// Repository is an implementation of RepositoryInterface using DynamoDB
type Repository struct {
	client *dynamodb.DynamoDB
}

// NewRepository creates a new User repository instance
func NewRepository(client *dynamodb.DynamoDB) *Repository {
	return &Repository{client}
}

// Add is used to insert a new User document
func (repository *Repository) Add(user *User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}
	_, err = repository.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

// GetByUsername looks up a User with the specified username
func (repository *Repository) GetByUsername(username string) (*User, error) {
	indexName := usernameIndex
	conditionExpression := "username = :username"
	var limit int64 = 1
	result, err := repository.client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              &indexName,
		KeyConditionExpression: &conditionExpression,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(strings.ToLower(username)),
			},
		},
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
