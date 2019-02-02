package main

import (
	"time"

	"github.com/bitterbattles/api/pkg/common/handlers"

	"github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/input"
	"github.com/bitterbattles/api/pkg/votes"
)

const battleIDLength = 20

// Handler represents a request handler
type Handler struct {
	repository votes.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository votes.RepositoryInterface) *handlers.APIHandler {
	handler := Handler{
		repository: repository,
	}
	return handlers.NewAPIHandler(handler.Handle)
}

// Handle handles a request
func (handler *Handler) Handle(request *Request) error {
	battleID, err := handler.sanitizeBattleID(request.BattleID)
	if err != nil {
		return err
	}
	vote := votes.Vote{
		UserID:    "bgttr132fopt0uo06vlg",
		BattleID:  battleID,
		IsVoteFor: request.IsVoteFor,
		CreatedOn: time.Now().Unix(),
	}
	return handler.repository.Add(vote)
}

func (handler *Handler) sanitizeBattleID(battleID string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		Length:    battleIDLength,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid Battle ID: " + message)
	}
	return input.SanitizeString(battleID, rules, errorCreator)
}
