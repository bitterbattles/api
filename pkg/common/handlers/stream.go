package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoEventChange represents a DynamoDB event record change
type DynamoEventChange struct {
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
	OldImage map[string]*dynamodb.AttributeValue `json:"OldImage"`
}

// DynamoEventRecord represents a DynamoDB event record
type DynamoEventRecord struct {
	Change    DynamoEventChange `json:"dynamodb"`
	EventName string            `json:"eventName"`
	EventID   string            `json:"eventID"`
}

// DynamoEvent represents a DynamoDB event
type DynamoEvent struct {
	Records []DynamoEventRecord `json:"records"`
}

// StreamHandlerInterface defines an interface for handling DynamoDB stream events
type StreamHandlerInterface interface {
	Handle(*DynamoEvent) error
}

// StreamHandler represents a lambda handler for DynamoDB stream events
type StreamHandler struct {
	StreamHandlerInterface
}

// NewStreamHandler creates a new StreamHandler instance
func NewStreamHandler(handler StreamHandlerInterface) *StreamHandler {
	return &StreamHandler{handler}
}

// Invoke handles a DynamoDB stream event
func (handler *StreamHandler) Invoke(context context.Context, inputBytes []byte) ([]byte, error) {
	var err error
	event := DynamoEvent{}
	err = json.Unmarshal(inputBytes, &event)
	if err != nil {
		return nil, err
	}
	err = handler.Handle(&event)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
