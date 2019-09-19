package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battlesget"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/votes"
)

// Processor represents a request processor
type Processor struct {
	indexer           *battles.Indexer
	battlesRepository battles.RepositoryInterface
	votesRepository   votes.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(
	indexer *battles.Indexer,
	battlesRepository battles.RepositoryInterface,
	votesRepository votes.RepositoryInterface) *Processor {
	return &Processor{
		indexer:           indexer,
		battlesRepository: battlesRepository,
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
	responses, err := battlesget.CreateResponses(userID, battleIDs, processor.battlesRepository, processor.getCanVote)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(responses)
	return output, nil
}

func (processor *Processor) getCanVote(userID string, battle *battles.Battle) (bool, error) {
	if userID == "" || userID == battle.UserID {
		return false, nil
	}
	vote, err := processor.votesRepository.GetByUserAndBattleIDs(userID, battle.ID)
	if err != nil {
		return false, err
	}
	return (vote == nil), nil
}
