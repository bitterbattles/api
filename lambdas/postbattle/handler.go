package postbattle

import (
	"time"

	"github.com/bitterbattles/api/common/handlers"
	"github.com/bitterbattles/api/common/loggers"

	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/errors"
	"github.com/bitterbattles/api/common/input"
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
func NewHandler(repository battles.RepositoryInterface, logger loggers.LoggerInterface) *handlers.APIHandler {
	handler := Handler{
		repository: repository,
	}
	return handlers.NewAPIHandler(handler.Handle, logger)
}

// Handle handles a request
func (handler *Handler) Handle(request *Request) error {
	title, err := handler.sanitizeTitle(request.Title)
	if err != nil {
		return err
	}
	description, err := handler.sanitizeDescription(request.Description)
	if err != nil {
		return err
	}
	battle := battles.Battle{
		ID:          xid.New().String(),
		UserID:      "bgttr132fopt0uo06vlg",
		Title:       title,
		Description: description,
		CreatedOn:   time.Now().Unix(),
	}
	return handler.repository.Add(battle)
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
