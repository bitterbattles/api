package main

import (
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
)

const (
	idParam = "id"
)

// Processor represents a request processor
type Processor struct {
	repository comments.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository comments.RepositoryInterface) *Processor {
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
	commentID := input.PathParams[idParam]
	comment, err := processor.repository.GetByID(commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil || comment.UserID != userID {
		output := api.NewOutputWithStatus(http.NotFound, nil)
		return output, nil
	}
	err = processor.repository.DeleteByID(commentID)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.NoContent, nil)
	return output, nil
}
