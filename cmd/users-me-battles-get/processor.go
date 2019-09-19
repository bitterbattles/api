package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battlesget"
	"github.com/bitterbattles/api/pkg/lambda/api"
)

// Processor represents a request processor
type Processor struct {
	indexer    *battles.Indexer
	repository battles.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *battles.Indexer, repository battles.RepositoryInterface) *Processor {
	return &Processor{
		indexer:    indexer,
		repository: repository,
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
	userID := input.AuthContext.UserID
	battleIDs, err := processor.indexer.GetByAuthor(userID, sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses, err := battlesget.CreateResponses(userID, battleIDs, processor.repository, processor.getCanVote)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(responses)
	return output, nil
}

func (processor *Processor) getCanVote(userID string, battle *battles.Battle) (bool, error) {
	return false, nil
}
