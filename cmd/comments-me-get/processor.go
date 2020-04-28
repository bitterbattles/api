package main

import (
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/commentsget"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	idParam = "id"
)

// Processor represents a request processor
type Processor struct {
	commentsIndex      comments.IndexInterface
	commentsRepository comments.RepositoryInterface
	usersRepository    users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(commentsIndex comments.IndexInterface, commentsRepository comments.RepositoryInterface, usersRepository users.RepositoryInterface) *Processor {
	return &Processor{
		commentsIndex:      commentsIndex,
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
	page := commentsget.GetPage(input)
	pageSize := commentsget.GetPageSize(input)
	userID := input.AuthContext.UserID
	commentIDs, err := processor.commentsIndex.GetByAuthor(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses, err := commentsget.CreateResponses(commentIDs, false, processor.commentsRepository, processor.usersRepository)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(responses)
	return output, nil
}
