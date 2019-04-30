package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

// Processor represents a request processor
type Processor struct {
	indexer           *battles.Indexer
	battlesRepository battles.RepositoryInterface
	usersRepository   users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *battles.Indexer, battlesRepository battles.RepositoryInterface, usersRepository users.RepositoryInterface) *Processor {
	return &Processor{
		indexer:           indexer,
		battlesRepository: battlesRepository,
		usersRepository:   usersRepository,
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
	userID := input.AuthContext.UserID
	user, err := processor.usersRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}
	battleIDs, err := processor.indexer.GetByAuthor(userID, sort, page, pageSize)
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
			responses = append(responses, battles.ToGetResponse(battle, user, true))
		} else {
			log.Println("Failed to find battle ID", battleID, "referenced in author", userID, "+", sort, "index.")
		}
	}
	output := api.NewOutput(responses)
	return output, nil
}
