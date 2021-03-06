package main

import (
	"github.com/bitterbattles/api/pkg/crypto"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/loginspost"
	"github.com/bitterbattles/api/pkg/users"
)

// Processor represents a request processor
type Processor struct {
	repository         users.RepositoryInterface
	accessTokenSecret  string
	refreshTokenSecret string
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository users.RepositoryInterface, accessTokenSecret string, refreshTokenSecret string) *Processor {
	return &Processor{
		repository:         repository,
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return &Request{}
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	request, _ := input.RequestBody.(*Request)
	username, err := processor.sanitizeUsername(request.Username)
	if err != nil {
		return nil, err
	}
	password, err := processor.sanitizePassword(request.Password)
	if err != nil {
		return nil, err
	}
	user, err := processor.repository.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewBadRequestError("Invalid credentials.")
	}
	passwordsMatch := crypto.VerifyPasswordHash(password, user.PasswordHash)
	if !passwordsMatch {
		return nil, errors.NewBadRequestError("Invalid credentials.")
	}
	response, err := loginspost.CreateResponse(user.ID, processor.accessTokenSecret, processor.refreshTokenSecret)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(response)
	return output, nil
}

func (processor *Processor) sanitizeUsername(value string) (string, error) {
	rules := input.StringRules{
		Required: true,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid email or username: " + message)
	}
	return input.SanitizeString(value, rules, errorCreator)
}

func (processor *Processor) sanitizePassword(value string) (string, error) {
	rules := input.StringRules{
		Required: true,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid password: " + message)
	}
	return input.SanitizeString(value, rules, errorCreator)
}
