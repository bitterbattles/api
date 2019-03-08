package main

import (
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/votes"
)

const battleIDLength = 20

// Processor represents a request processor
type Processor struct {
	repository votes.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository votes.RepositoryInterface) *Processor {
	return &Processor{
		repository: repository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return &Request{}
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	request, _ := input.RequestBody.(*Request)
	battleID, err := processor.sanitizeBattleID(request.BattleID)
	if err != nil {
		return nil, err
	}
	vote := votes.Vote{
		UserID:    input.AuthContext.UserID,
		BattleID:  battleID,
		IsVoteFor: request.IsVoteFor,
		CreatedOn: time.NowUnix(),
	}
	err = processor.repository.Add(vote)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.Created, nil)
	return output, nil
}

func (processor *Processor) sanitizeBattleID(battleID string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		Length:    battleIDLength,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid Battle ID: " + message)
	}
	return input.SanitizeString(battleID, rules, errorCreator)
}
