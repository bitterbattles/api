package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
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
	sort := battles.GetSort(input)
	page := battles.GetPage(input)
	pageSize := battles.GetPageSize(input)
	var userID string
	if input.AuthContext != nil {
		userID = input.AuthContext.UserID
	}
	battleIDs, err := processor.indexer.GetGlobal(sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses := make([]*battles.Response, 0, len(battleIDs))
	for _, battleID := range battleIDs {
		battle, err := processor.battlesRepository.GetByID(battleID)
		if err != nil {
			return nil, err
		}
		if battle != nil {
			user, err := processor.usersRepository.GetByID(battle.UserID)
			if err != nil {
				return nil, err
			}
			var vote *votes.Vote
			if userID != "" {
				vote, err = processor.votesRepository.GetByUserAndBattleIDs(userID, battleID)
				if err != nil {
					return nil, err
				}
			}
			responses = append(responses, battles.ToGetResponse(battle, user, vote != nil))
		} else {
			log.Println("Failed to find battle ID", battleID, "referenced in", sort, "index.")
		}
	}
	output := api.NewOutput(responses)
	return output, nil
}
