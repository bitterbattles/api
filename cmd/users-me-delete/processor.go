package main

import (
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	idParam = "id"
)

// Processor represents a request processor
type Processor struct {
	repository users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository users.RepositoryInterface) *Processor {
	return &Processor{
		repository: repository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return nil
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	userID := input.AuthContext.UserID
	err := processor.repository.DeleteByID(userID)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.NoContent, nil)
	return output, nil
}
