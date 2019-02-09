package main

import (
	"encoding/json"
	"time"

	"github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/http"
	"github.com/bitterbattles/api/pkg/common/input"
	"github.com/bitterbattles/api/pkg/common/lambda/api"
	"github.com/bitterbattles/api/pkg/votes"
)

const battleIDLength = 20

// Handler represents a request handler
type Handler struct {
	repository votes.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository votes.RepositoryInterface) *api.Handler {
	handler := Handler{
		repository: repository,
	}
	return api.NewHandler(&handler)
}

// Handle handles a request
func (handler *Handler) Handle(request *http.Request) (*http.Response, error) {
	requestBody := Request{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return nil, errors.NewBadRequestError("Failed to decode request JSON.")
	}
	battleID, err := handler.sanitizeBattleID(requestBody.BattleID)
	if err != nil {
		return nil, err
	}
	vote := votes.Vote{
		UserID:    "bgttr132fopt0uo06vlg",
		BattleID:  battleID,
		IsVoteFor: requestBody.IsVoteFor,
		CreatedOn: time.Now().Unix(),
	}
	err = handler.repository.Add(vote)
	if err != nil {
		return nil, err
	}
	return http.NewResponseWithStatus(nil, nil, http.Created)
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
