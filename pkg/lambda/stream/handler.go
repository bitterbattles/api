package stream

import (
	"context"
	"encoding/json"
)

// ProcessorInterface defines an interface for processing DynamoDB stream events
type ProcessorInterface interface {
	Process(*Event) error
}

// Handler represents a lambda handler for DynamoDB stream events
type Handler struct {
	processor ProcessorInterface
}

// NewHandler creates a new Handler instance
func NewHandler(processor ProcessorInterface) *Handler {
	return &Handler{
		processor: processor,
	}
}

// Invoke handles a DynamoDB stream event
func (handler *Handler) Invoke(context context.Context, inputBytes []byte) ([]byte, error) {
	var err error
	event := Event{}
	err = json.Unmarshal(inputBytes, &event)
	if err != nil {
		return nil, err
	}
	err = handler.processor.Process(&event)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
