package stream

import "github.com/aws/aws-sdk-go/service/dynamodb"

// EventChange represents a DynamoDB event record change
type EventChange struct {
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
	OldImage map[string]*dynamodb.AttributeValue `json:"OldImage"`
}

// EventRecord represents a DynamoDB event record
type EventRecord struct {
	Change  EventChange `json:"dynamodb"`
	EventID string      `json:"eventID"`
}

// Event represents a DynamoDB event
type Event struct {
	Records []EventRecord `json:"records"`
}
