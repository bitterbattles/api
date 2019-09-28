package comments

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName       = "comments"
	idFieldName     = "id"
	userIDFieldName = "userId"
)

// RepositoryInterface defines an interface for a Comment repository
type RepositoryInterface interface {
	Add(*Comment) error
	DeleteByID(string) error
	GetByID(string) (*Comment, error)
	UpdateUsername(string, string) error
}

// Repository is an implementation of RepositoryInterface using DynamoDB
type Repository struct {
	client *dynamodb.DynamoDB
}

// NewRepository creates a new Comments repository instance
func NewRepository(client *dynamodb.DynamoDB) *Repository {
	return &Repository{client}
}

// Add is used to insert a new Comment document
func (repository *Repository) Add(comment *Comment) error {
	item, err := dynamodbattribute.MarshalMap(comment)
	if err != nil {
		return err
	}
	_, err = repository.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

// DeleteByID deletes a Comment by ID
func (repository *Repository) DeleteByID(id string) error {
	conditionExpression := "id = :id"
	updateExpression := "SET #state = :state"
	stateAttribute := "state"
	_, err := repository.client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			idFieldName: {
				S: aws.String(id),
			},
		},
		ConditionExpression: &conditionExpression,
		UpdateExpression:    &updateExpression,
		ExpressionAttributeNames: map[string]*string{
			"#state": &stateAttribute,
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
			":state": {
				N: aws.String(fmt.Sprintf("%d", Deleted)),
			},
		},
	})
	return err
}

// GetByID is used to get a Comment by ID
func (repository *Repository) GetByID(id string) (*Comment, error) {
	result, err := repository.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			idFieldName: {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, nil
	}
	comment := &Comment{}
	err = dynamodbattribute.UnmarshalMap(result.Item, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// UpdateUsername updates the specified user's username
func (repository *Repository) UpdateUsername(userID string, username string) error {
	conditionExpression := "userId = :userId"
	updateExpression := "SET username = :username"
	usernameAttribute := "username"
	_, err := repository.client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			userIDFieldName: {
				S: aws.String(userID),
			},
		},
		ConditionExpression: &conditionExpression,
		UpdateExpression:    &updateExpression,
		ExpressionAttributeNames: map[string]*string{
			"username": &usernameAttribute,
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {
				S: aws.String(userID),
			},
			":username": {
				S: aws.String(username),
			},
		},
	})
	return err
}
