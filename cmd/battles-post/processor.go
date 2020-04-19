package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/guid"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	minTitleLength       = 4
	maxTitleLength       = 50
	minDescriptionLength = 4
	maxDescriptionLength = 500
)

// Processor represents a request processor
type Processor struct {
	battlesRepository battles.RepositoryInterface
	usersRepository   users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(battlesRepository battles.RepositoryInterface, usersRepository users.RepositoryInterface) *Processor {
	return &Processor{
		battlesRepository: battlesRepository,
		usersRepository:   usersRepository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return &Request{}
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	request, _ := input.RequestBody.(*Request)
	title, err := processor.sanitizeTitle(request.Title)
	if err != nil {
		return nil, err
	}
	description, err := processor.sanitizeDescription(request.Description)
	if err != nil {
		return nil, err
	}
	userID := input.AuthContext.UserID
	user, err := processor.usersRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewForbiddenError("User not found.")
	}
	battle := battles.Battle{
		ID:          guid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		CreatedOn:   time.NowUnix(),
		State:       battles.Active,
	}
	err = processor.battlesRepository.Add(&battle)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.Created, nil)
	return output, nil
}

func (processor *Processor) sanitizeTitle(title string) (string, error) {
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

func (processor *Processor) sanitizeDescription(description string) (string, error) {
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
