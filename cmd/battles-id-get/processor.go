package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battlesget"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/votes"
)

const idParam = "id"

// Processor represents a request processor
type Processor struct {
	battlesRepository battles.RepositoryInterface
	usersRepository   users.RepositoryInterface
	votesRepository   votes.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(battlesRepository battles.RepositoryInterface, usersRepository users.RepositoryInterface, votesRepository votes.RepositoryInterface) *Processor {
	return &Processor{
		battlesRepository: battlesRepository,
		usersRepository:   usersRepository,
		votesRepository:   votesRepository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return nil
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	battleID := input.PathParams[idParam]
	var userID string
	if input.AuthContext != nil {
		userID = input.AuthContext.UserID
	}
	response, err := battlesget.CreateResponse(
		userID,
		battleID,
		processor.battlesRepository,
		processor.usersRepository,
		processor.votesRepository)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.NewNotFoundError("Battle with ID " + battleID + " is missing or deleted.")
	}
	output := api.NewOutput(response)
	return output, nil
}
