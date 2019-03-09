package votes

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName         = "votes"
	userIDFieldName   = "userId"
	battleIDFieldName = "battleId"
)

// RepositoryInterface defines an interface for a Vote repository
type RepositoryInterface interface {
	Add(*Vote) error
	GetByUserAndBattleIDs(string, string) (*Vote, error)
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
func (repository *Repository) Add(vote *Vote) error {
	item, err := dynamodbattribute.MarshalMap(vote)
	if err != nil {
		return err
	}
	conditionExpression := "attribute_not_exists(" + battleIDFieldName + ")"
	_, err = repository.client.PutItem(&dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(tableName),
		ConditionExpression: &conditionExpression,
	})
	return err
}

// GetByUserAndBattleIDs is used to get a Vote by user ID and battle ID
func (repository *Repository) GetByUserAndBattleIDs(userID string, battleID string) (*Vote, error) {
	result, err := repository.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			userIDFieldName: {
				S: aws.String(userID),
			},
			battleIDFieldName: {
				S: aws.String(battleID),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	vote := &Vote{}
	err = dynamodbattribute.UnmarshalMap(result.Item, vote)
	if err != nil {
		return nil, err
	}
	return vote, nil
}
