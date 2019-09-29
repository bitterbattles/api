package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battlesget"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/votes"
)

// Processor represents a request processor
type Processor struct {
	indexer           *battles.Indexer
	battlesRepository battles.RepositoryInterface
	usersRepository   users.RepositoryInterface
	votesRepository   votes.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(
	indexer *battles.Indexer,
	battlesRepository battles.RepositoryInterface,
	usersRepository users.RepositoryInterface,
	votesRepository votes.RepositoryInterface) *Processor {
	return &Processor{
		indexer:           indexer,
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
	sort := battlesget.GetSort(input)
	page := battlesget.GetPage(input)
	pageSize := battlesget.GetPageSize(input)
	battleIDs, err := processor.indexer.GetGlobal(sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	var userID string
	if input.AuthContext != nil {
		userID = input.AuthContext.UserID
	}
	responses, err := battlesget.CreateResponses(
		userID,
		battleIDs,
		processor.battlesRepository,
		processor.usersRepository,
		processor.votesRepository)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(responses)
	return output, nil
}
