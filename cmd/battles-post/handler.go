package main

import (
	"encoding/json"
	"time"

	"github.com/bitterbattles/api/pkg/common/http"
	"github.com/bitterbattles/api/pkg/common/lambda/api"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/input"
	"github.com/rs/xid"
)

const (
	minTitleLength       = 4
	maxTitleLength       = 50
	minDescriptionLength = 4
	maxDescriptionLength = 500
)

// Handler represents a request handler
type Handler struct {
	repository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository battles.RepositoryInterface) *api.Handler {
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
	title, err := handler.sanitizeTitle(requestBody.Title)
	if err != nil {
		return nil, err
	}
	description, err := handler.sanitizeDescription(requestBody.Description)
	if err != nil {
		return nil, err
	}
	battle := battles.Battle{
		ID:          xid.New().String(),
		UserID:      "bgttr132fopt0uo06vlg",
		Title:       title,
		Description: description,
		CreatedOn:   time.Now().Unix(),
		State:       battles.Active,
	}
	err = handler.repository.Add(battle)
	if err != nil {
		return nil, err
	}
	return http.NewResponseWithStatus(nil, nil, http.Created)
}

func (handler *Handler) sanitizeTitle(title string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		MinLength: minTitleLength,
		MaxLength: maxTitleLength,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid title: " + message)
	}
	return input.SanitizeString(title, rules, errorCreator)
}

func (handler *Handler) sanitizeDescription(description string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		MinLength: minDescriptionLength,
		MaxLength: maxDescriptionLength,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid description: " + message)
	}
	return input.SanitizeString(description, rules, errorCreator)
}
