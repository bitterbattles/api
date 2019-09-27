package main

import (
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/commentsget"
	"github.com/bitterbattles/api/pkg/lambda/api"
)

const (
	idParam = "id"
)

// Processor represents a request processor
type Processor struct {
	indexer    *comments.Indexer
	repository comments.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *comments.Indexer, repository comments.RepositoryInterface) *Processor {
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
	page := commentsget.GetPage(input)
	pageSize := commentsget.GetPageSize(input)
	userID := input.AuthContext.UserID
	commentIDs, err := processor.indexer.GetByAuthor(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses, err := commentsget.CreateResponses(commentIDs, processor.repository)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(responses)
	return output, nil
}
