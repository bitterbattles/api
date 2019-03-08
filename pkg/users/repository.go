package users

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName         = "users"
	idFieldName       = "id"
	usernameFieldName = "username"
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
		Item:      item,
		TableName: aws.String(tableName),
	})
	return err
}

// GetByUsername looks up a User with the specified username
func (repository *Repository) GetByUsername(username string) (*User, error) {
	result, err := repository.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			usernameFieldName: {
				S: aws.String(strings.ToLower(username)),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
