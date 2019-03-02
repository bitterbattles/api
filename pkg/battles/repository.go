package battles

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName   = "battles"
	idFieldName = "id"
)

// RepositoryInterface defines an interface for a Battle repository
type RepositoryInterface interface {
	Add(Battle) error
	DeleteByID(string) error
	GetByID(string) (*Battle, error)
	IncrementVotes(string, int, int) error
}

// Repository is an implementation of RepositoryInterface using DynamoDB
type Repository struct {
	client *dynamodb.DynamoDB
}

// NewRepository creates a new Battles repository instance
func NewRepository(client *dynamodb.DynamoDB) *Repository {
	return &Repository{client}
}

// Add is used to insert a new Battle document
func (repository *Repository) Add(battle Battle) error {
	item, err := dynamodbattribute.MarshalMap(battle)
	if err != nil {
		return err
	}
	_, err = repository.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	return err
}

// DeleteByID deletes a Battle by ID
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

// GetByID is used to get a Battle by ID
func (repository *Repository) GetByID(id string) (*Battle, error) {
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
	battle := Battle{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &battle)
	if err != nil {
		return nil, err
	}
	return &battle, nil
}

// IncrementVotes increments the votes for a given Battle ID
func (repository *Repository) IncrementVotes(id string, deltaVotesFor int, deltaVotesAgainst int) error {
	conditionExpression := "id = :id"
	updateExpression := "ADD votesFor :deltaVotesFor, votesAgainst :deltaVotesAgainst"
	_, err := repository.client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			idFieldName: {
				S: aws.String(id),
			},
		},
		ConditionExpression: &conditionExpression,
		UpdateExpression:    &updateExpression,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
			":deltaVotesFor": {
				N: aws.String(fmt.Sprintf("%d", deltaVotesFor)),
			},
			":deltaVotesAgainst": {
				N: aws.String(fmt.Sprintf("%d", deltaVotesAgainst)),
			},
		},
	})
	return err
}
