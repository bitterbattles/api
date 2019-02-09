package stream

import (
	"context"
	"encoding/json"
)

// HandlerInterface defines an interface for handling DynamoDB stream events
type HandlerInterface interface {
	Handle(*Event) error
}

// Handler represents a lambda handler for DynamoDB stream events
type Handler struct {
	HandlerInterface
}

// NewHandler creates a new Handler instance
func NewHandler(handler HandlerInterface) *Handler {
	return &Handler{handler}
}

// Invoke handles a DynamoDB stream event
func (handler *Handler) Invoke(context context.Context, inputBytes []byte) ([]byte, error) {
	var err error
	event := Event{}
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
