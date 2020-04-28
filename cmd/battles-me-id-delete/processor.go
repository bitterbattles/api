package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
)

const (
	idParam = "id"
)

// Processor represents a request processor
type Processor struct {
	repository battles.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository battles.RepositoryInterface) *Processor {
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
	battleID := input.PathParams[idParam]
	battle, err := processor.repository.GetByID(battleID)
	if err != nil {
		return nil, err
	}
	if battle == nil || battle.UserID != userID {
		output := api.NewOutputWithStatus(http.NotFound, nil)
		return output, nil
	}
	err = processor.repository.DeleteByID(battleID)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.NoContent, nil)
	return output, nil
}
