package main

import (
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/commentsget"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	idParam         = "id"
	pageParam       = "page"
	pageSizeParam   = "pageSize"
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
)

// Processor represents a request processor
type Processor struct {
	indexer            *comments.Indexer
	commentsRepository comments.RepositoryInterface
	usersRepository    users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *comments.Indexer, commentsRepository comments.RepositoryInterface, usersRepository users.RepositoryInterface) *Processor {
	return &Processor{
		indexer:            indexer,
		commentsRepository: commentsRepository,
		usersRepository:    usersRepository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return nil
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	battleID := input.PathParams[idParam]
	page := commentsget.GetPage(input)
	pageSize := commentsget.GetPageSize(input)
	commentIDs, err := processor.indexer.GetByBattle(battleID, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses, err := commentsget.CreateResponses(commentIDs, true, processor.commentsRepository, processor.usersRepository)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(responses)
	return output, nil
}
